package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// JSON marshals v and writes it as JSON to the response writer
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// If encoding fails, log it but we've already written the status
		// In production, you'd want proper logging here
		http.Error(w, `{"error":"internal encoding error"}`, http.StatusInternalServerError)
	}
}

// Error writes a JSON error response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}

// HealthResponse represents a standard health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// Health writes a standard health check response
func Health(w http.ResponseWriter) {
	JSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

