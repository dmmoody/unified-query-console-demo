package eip

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	commonhttp "ach-concourse/internal/common/http"
)

// Handler handles HTTP requests for EIP
type Handler struct {
	service *Service
}

// NewHandler creates a new EIP handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all EIP routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/cases", func(r chi.Router) {
		r.Post("/", h.CreateCase)
		r.Get("/", h.ListCases)
		r.Get("/{id}", h.GetCase)
		r.Patch("/{id}/status", h.UpdateStatus)
	})
	r.Get("/healthz", h.Health)
}

// CreateCase handles POST /api/v1/cases
func (h *Handler) CreateCase(w http.ResponseWriter, r *http.Request) {
	var req CreateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	eipCase, err := h.service.CreateCase(r.Context(), &req)
	if err != nil {
		commonhttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	commonhttp.JSON(w, http.StatusCreated, eipCase)
}

// ListCases handles GET /api/v1/cases
func (h *Handler) ListCases(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	side := r.URL.Query().Get("side")
	traceNumber := r.URL.Query().Get("trace_number")

	cases, err := h.service.ListCases(r.Context(), status, side, traceNumber)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to list cases")
		return
	}

	if cases == nil {
		cases = []*EIPCase{}
	}

	commonhttp.JSON(w, http.StatusOK, cases)
}

// GetCase handles GET /api/v1/cases/{id}
func (h *Handler) GetCase(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	eipCase, err := h.service.GetCase(r.Context(), id)
	if err != nil {
		commonhttp.Error(w, http.StatusInternalServerError, "failed to get case")
		return
	}

	if eipCase == nil {
		commonhttp.Error(w, http.StatusNotFound, "case not found")
		return
	}

	commonhttp.JSON(w, http.StatusOK, eipCase)
}

// UpdateStatus handles PATCH /api/v1/cases/{id}/status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		commonhttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	eipCase, err := h.service.UpdateCaseStatus(r.Context(), id, req.Status)
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

