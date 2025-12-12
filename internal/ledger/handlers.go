package ledger

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	commonhttp "ach-concourse/internal/common/http"
)

// Handler handles HTTP requests for ledger
type Handler struct {
	service *Service
}

// NewHandler creates a new ledger handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all ledger routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/postings", func(r chi.Router) {
		r.Post("/", h.CreatePosting)
		r.Get("/", h.ListPostings)
	})
	r.Get("/api/v1/balances", h.GetBalances)
	r.Get("/healthz", h.Health)
}

// CreatePosting handles POST /api/v1/postings
func (h *Handler) CreatePosting(w http.ResponseWriter, r *http.Request) {
	var req CreatePostingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.CreatePosting(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, entry)
}

// ListPostings handles GET /api/v1/postings
func (h *Handler) ListPostings(w http.ResponseWriter, r *http.Request) {
	achSide := r.URL.Query().Get("ach_side")
	traceNumber := r.URL.Query().Get("trace_number")

	entries, err := h.service.ListPostings(r.Context(), achSide, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list postings")
		return
	}

	if entries == nil {
		entries = []*LedgerEntry{}
	}

	commonhttp.JSON(w, http.StatusOK, entries)
}

// GetBalances handles GET /api/v1/balances
func (h *Handler) GetBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := h.service.GetBalances(r.Context())
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get balances")
		return
	}

	commonhttp.JSON(w, http.StatusOK, balances)
}

// Health handles GET /healthz
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	commonhttp.Health(w)
}

