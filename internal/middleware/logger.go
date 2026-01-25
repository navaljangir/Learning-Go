// Package middleware contains HTTP middleware functions.
// Middleware can intercept requests before they reach handlers
// and responses before they are sent to clients.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger is middleware that logs incoming HTTP requests.
// It logs the method, path, and duration of each request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request details
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
