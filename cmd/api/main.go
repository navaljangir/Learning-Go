// Package main is the entry point for the API server.
// It initializes and starts the HTTP server on port 8080.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tejas/learningGo/internal/routes"
)

func main() {
	// Initialize the router with all routes
	router := routes.SetupRoutes()

	// Server configuration
	port := ":8080"

	fmt.Printf("Starting server on http://localhost%s\n", port)

	// Start the HTTP server
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
