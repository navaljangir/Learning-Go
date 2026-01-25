// Package handlers contains HTTP handler functions for the API.
// Handlers are responsible for processing requests and returning responses.
package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the health check response structure.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthCheck handles the /health endpoint.
// It returns a simple JSON response indicating the service is healthy.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Message: "Service is healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
