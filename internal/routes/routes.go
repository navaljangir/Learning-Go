// Package routes defines the HTTP routing configuration for the API.
// It maps URL paths to their corresponding handlers.
package routes

import (
	"net/http"

	"github.com/tejas/learningGo/internal/handlers"
	"github.com/tejas/learningGo/internal/middleware"
)

// SetupRoutes initializes and returns the HTTP router with all routes configured.
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handlers.HealthCheck)

	// Root endpoint
	mux.HandleFunc("/", handlers.HealthCheck)

	// Wrap the router with middleware
	return middleware.Logger(mux)
}
