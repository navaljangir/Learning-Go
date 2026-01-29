package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs request details
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get request details
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get response status
		statusCode := c.Writer.Status()

		// Log format: [HTTP] METHOD PATH STATUS LATENCY IP
		log.Printf("[HTTP] %s %s %d %v %s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}
