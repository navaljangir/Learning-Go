package main

import (
	"context"
	"gin_server/constants"
	"gin_server/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize the server
	s := server.NewServer()

	// Create HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         constants.ServerPort,
		Handler:      s.Router(),
		ReadTimeout:  constants.ServerReadTimeout,
		WriteTimeout: constants.ServerWriteTimeout,
	}

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)

	// Notify on SIGINT (Ctrl+C) and SIGTERM (kill command)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on http://localhost" + constants.ServerPort)
		log.Println("")
		log.Println("Available endpoints:")
		log.Println("  POST /api/auth/register  - Register new user")
		log.Println("  POST /api/auth/login     - Login and get token")
		log.Println("  GET  /api/users/profile  - Get profile (requires auth)")
		log.Println("  GET  /api/users          - Get all users (requires auth)")
		log.Println("  GET  /health             - Health check")
		log.Println("")
		log.Println("Press Ctrl+C to stop the server gracefully")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Block until we receive a signal
	sig := <-quit
	log.Printf("\nReceived signal: %v", sig)
	log.Println("Shutting down server gracefully...")

	// Create context with timeout for graceful shutdown
	// Server has 10 seconds to finish processing existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}
