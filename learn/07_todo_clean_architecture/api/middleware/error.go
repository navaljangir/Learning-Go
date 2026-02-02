package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorRecoveryMiddleware recovers from panics and returns a 500 error
func ErrorRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "internal server error",
					"message": "an unexpected error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
