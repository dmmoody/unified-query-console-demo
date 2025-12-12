package console

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	commonhttp "ach-concourse/internal/common/http"
)

// Handler handles HTTP requests for console
type Handler struct {
	service *Service
}

// NewHandler creates a new console handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all console routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	// Unified ACH items (legacy endpoints for backward compatibility)
	r.Route("/api/v1/ach-items", func(r chi.Router) {
		r.Get("/", h.GetAchItems)
		r.Get("/{side}/{id}", h.GetAchItem)
		r.Post("/{side}/{id}/return", h.ReturnEntry)
	})

	// ODFI operations via gateway
	r.Route("/api/v1/odfi/entries", func(r chi.Router) {
		r.Post("/", h.CreateODFIEntry)
		r.Get("/", h.ListODFIEntries)
		r.Get("/{id}", h.GetODFIEntry)
		r.Patch("/{id}/status", h.UpdateODFIStatus)
	})

	// RDFI operations via gateway
	r.Route("/api/v1/rdfi/entries", func(r chi.Router) {
		r.Post("/", h.CreateRDFIEntry)
		r.Get("/", h.ListRDFIEntries)
		r.Get("/{id}", h.GetRDFIEntry)
		r.Post("/{id}/return", h.ReturnRDFIEntry)
	})

	// Ledger operations via gateway
	r.Route("/api/v1/ledger", func(r chi.Router) {
		r.Post("/postings", h.CreateLedgerPosting)
		r.Get("/postings", h.ListLedgerPostings)
		r.Get("/balances", h.GetBalances)
	})

	// EIP operations via gateway
	r.Route("/api/v1/eip/cases", func(r chi.Router) {
		r.Post("/", h.CreateEIPCase)
		r.Get("/", h.ListEIPCases)
		r.Get("/{id}", h.GetEIPCase)
		r.Patch("/{id}/status", h.UpdateEIPCaseStatus)
	})

	r.Get("/healthz", h.Health)
}

// ========== Legacy Unified ACH Items Handlers ==========

// GetAchItems handles GET /api/v1/ach-items
func (h *Handler) GetAchItems(w http.ResponseWriter, r *http.Request) {
	side := r.URL.Query().Get("side")
	status := r.URL.Query().Get("status")
	traceNumber := r.URL.Query().Get("trace_number")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	// Pagination parameters
	limit := 100 // Default limit
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 1000 {
				parsedLimit = 1000 // Max limit for safety
			}
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Validate side if provided
	if side != "" {
		side = strings.ToUpper(side)
		if side != "ODFI" && side != "RDFI" {
			commonhttp.Error(w, http.StatusBadRequest, "side must be ODFI or RDFI")
			return
		}
	}

	// Validate sort_order if provided
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		commonhttp.Error(w, http.StatusBadRequest, "sort_order must be 'asc' or 'desc'")
		return
	}

	// Validate sort_by if provided
	validSortFields := map[string]bool{
		"created_at":   true,
		"status":       true,
		"amount":       true,
		"amount_cents": true,
		"trace_number": true,
		"side":         true,
	}
	if sortBy != "" && !validSortFields[sortBy] {
		commonhttp.Error(w, http.StatusBadRequest, "sort_by must be one of: created_at, status, amount, trace_number, side")
		return
	}

	items, err := h.service.GetAchItems(r.Context(), side, status, traceNumber, sortBy, sortOrder, limit, offset)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to fetch ACH items")
		return
	}

	if items == nil {
		items = []*UnifiedAchItem{}
	}

	commonhttp.JSON(w, http.StatusOK, items)
}

// GetAchItem handles GET /api/v1/ach-items/{side}/{id}
func (h *Handler) GetAchItem(w http.ResponseWriter, r *http.Request) {
	side := chi.URLParam(r, "side")
	id := chi.URLParam(r, "id")

	item, err := h.service.GetAchItem(r.Context(), side, id)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if item == nil {
		commonhttp.Error(w, http.StatusNotFound, "entry not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, item)
}

// ReturnEntry handles POST /api/v1/ach-items/{side}/{id}/return
func (h *Handler) ReturnEntry(w http.ResponseWriter, r *http.Request) {
	side := strings.ToUpper(chi.URLParam(r, "side"))
	id := chi.URLParam(r, "id")

	// Only support RDFI returns in this POC
	if side != "RDFI" {
		commonhttp.Error(w, http.StatusBadRequest, "only RDFI entries can be returned")
		return
	}

	var req ReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.ReturnEntry(r.Context(), id, req.Reason)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusOK, entry)
}

// ========== ODFI Handlers ==========

// CreateODFIEntry handles POST /api/v1/odfi/entries
func (h *Handler) CreateODFIEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateODFIEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.CreateODFIEntry(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, entry)
}

// ListODFIEntries handles GET /api/v1/odfi/entries
func (h *Handler) ListODFIEntries(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	traceNumber := r.URL.Query().Get("trace_number")

	entries, err := h.service.ListODFIEntries(r.Context(), status, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list ODFI entries")
		return
	}

	if entries == nil {
		entries = []*ODFIEntry{}
	}

	commonhttp.JSON(w, http.StatusOK, entries)
}

// GetODFIEntry handles GET /api/v1/odfi/entries/{id}
func (h *Handler) GetODFIEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entry, err := h.service.GetODFIEntry(r.Context(), id)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get ODFI entry")
		return
	}

	if entry == nil {
		commonhttp.Error(w, http.StatusNotFound, "entry not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, entry)
}

// UpdateODFIStatus handles PATCH /api/v1/odfi/entries/{id}/status
func (h *Handler) UpdateODFIStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateODFIStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.UpdateODFIEntryStatus(r.Context(), id, req.Status)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if entry == nil {
		commonhttp.Error(w, http.StatusNotFound, "entry not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, entry)
}

// ========== RDFI Handlers ==========

// CreateRDFIEntry handles POST /api/v1/rdfi/entries
func (h *Handler) CreateRDFIEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateRDFIEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.CreateRDFIEntry(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, entry)
}

// ListRDFIEntries handles GET /api/v1/rdfi/entries
func (h *Handler) ListRDFIEntries(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	traceNumber := r.URL.Query().Get("trace_number")

	entries, err := h.service.ListRDFIEntries(r.Context(), status, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list RDFI entries")
		return
	}

	if entries == nil {
		entries = []*RDFIEntry{}
	}

	commonhttp.JSON(w, http.StatusOK, entries)
}

// GetRDFIEntry handles GET /api/v1/rdfi/entries/{id}
func (h *Handler) GetRDFIEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entry, err := h.service.GetRDFIEntry(r.Context(), id)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get RDFI entry")
		return
	}

	if entry == nil {
		commonhttp.Error(w, http.StatusNotFound, "entry not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, entry)
}

// ReturnRDFIEntry handles POST /api/v1/rdfi/entries/{id}/return
func (h *Handler) ReturnRDFIEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req ReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.ReturnEntry(r.Context(), id, req.Reason)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if entry == nil {
		commonhttp.Error(w, http.StatusNotFound, "entry not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, entry)
}

// ========== Ledger Handlers ==========

// CreateLedgerPosting handles POST /api/v1/ledger/postings
func (h *Handler) CreateLedgerPosting(w http.ResponseWriter, r *http.Request) {
	var req CreateLedgerPostingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.CreateLedgerPosting(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, entry)
}

// ListLedgerPostings handles GET /api/v1/ledger/postings
func (h *Handler) ListLedgerPostings(w http.ResponseWriter, r *http.Request) {
	achSide := r.URL.Query().Get("ach_side")
	traceNumber := r.URL.Query().Get("trace_number")

	entries, err := h.service.ListLedgerPostings(r.Context(), achSide, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list ledger postings")
		return
	}

	if entries == nil {
		entries = []*LedgerEntry{}
	}

	commonhttp.JSON(w, http.StatusOK, entries)
}

// GetBalances handles GET /api/v1/ledger/balances
func (h *Handler) GetBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := h.service.GetBalances(r.Context())
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get balances")
		return
	}

	commonhttp.JSON(w, http.StatusOK, balances)
}

// ========== EIP Handlers ==========

// CreateEIPCase handles POST /api/v1/eip/cases
func (h *Handler) CreateEIPCase(w http.ResponseWriter, r *http.Request) {
	var req CreateEIPCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	eipCase, err := h.service.CreateEIPCase(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, eipCase)
}

// ListEIPCases handles GET /api/v1/eip/cases
func (h *Handler) ListEIPCases(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	side := r.URL.Query().Get("side")
	traceNumber := r.URL.Query().Get("trace_number")

	cases, err := h.service.ListEIPCases(r.Context(), status, side, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list EIP cases")
		return
	}

	if cases == nil {
		cases = []*EIPCase{}
	}

	commonhttp.JSON(w, http.StatusOK, cases)
}

// GetEIPCase handles GET /api/v1/eip/cases/{id}
func (h *Handler) GetEIPCase(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	eipCase, err := h.service.GetEIPCase(r.Context(), id)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get EIP case")
		return
	}

	if eipCase == nil {
		commonhttp.Error(w, http.StatusNotFound, "case not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, eipCase)
}

// UpdateEIPCaseStatus handles PATCH /api/v1/eip/cases/{id}/status
func (h *Handler) UpdateEIPCaseStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateEIPCaseStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	eipCase, err := h.service.UpdateEIPCaseStatus(r.Context(), id, req.Status)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if eipCase == nil {
		commonhttp.Error(w, http.StatusNotFound, "case not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, eipCase)
}

// Health handles GET /healthz
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	commonhttp.Health(w)
}
