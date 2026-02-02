package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details
		log.Printf(
			"[%s] %s %s - %d - %v - %s",
			method,
			path,
			clientIP,
			statusCode,
			duration,
			c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}
