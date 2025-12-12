package odfi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	commonhttp "ach-concourse/internal/common/http"
)

// Handler handles HTTP requests for ODFI
type Handler struct {
	service *Service
}

// NewHandler creates a new ODFI handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all ODFI routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/entries", func(r chi.Router) {
		r.Post("/", h.CreateEntry)
		r.Get("/", h.ListEntries)
		r.Get("/{id}", h.GetEntry)
		r.Patch("/{id}/status", h.UpdateStatus)
	})
	r.Get("/healthz", h.Health)
}

// CreateEntry handles POST /api/v1/entries
func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.CreateEntry(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, entry)
}

// ListEntries handles GET /api/v1/entries
func (h *Handler) ListEntries(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	traceNumber := r.URL.Query().Get("trace_number")

	entries, err := h.service.ListEntries(r.Context(), status, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list entries")
		return
	}

	if entries == nil {
		entries = []*ODFIEntry{}
	}

	commonhttp.JSON(w, http.StatusOK, entries)
}

// GetEntry handles GET /api/v1/entries/{id}
func (h *Handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entry, err := h.service.GetEntry(r.Context(), id)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get entry")
		return
	}

	if entry == nil {
		commonhttp.Error(w, http.StatusNotFound, "entry not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, entry)
}

// UpdateStatus handles PATCH /api/v1/entries/{id}/status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.service.UpdateEntryStatus(r.Context(), id, req.Status)
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

// Health handles GET /healthz
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	commonhttp.Health(w)
}

