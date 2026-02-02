package middleware

import (
	"log"
	"time"
	"todo_app/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates and injects a unique request ID for tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store request ID in context
		c.Set(constants.ContextRequestID, requestID)

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggerMiddleware logs HTTP requests with request ID
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		requestID := c.GetString(constants.ContextRequestID)

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details with request ID
		log.Printf(
			"[%s] RequestID=%s %s %s - %d - %v - %s",
			method,
			requestID,
			path,
			clientIP,
			statusCode,
			duration,
			c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}
