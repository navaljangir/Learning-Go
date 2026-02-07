package middleware

import (
	"errors"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware catches errors from handlers and returns appropriate HTTP responses
// This centralizes error handling so handlers don't need to check error types
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error (most recent)
			err := c.Errors.Last().Err

			// Check if it's an AppError with status code
			var appErr *utils.AppError
			if errors.As(err, &appErr) {
				// Use the status code from AppError
				c.JSON(appErr.StatusCode, gin.H{
					"success": false,
					"error":   appErr.Message,
				})
				return
			}

			// Check for specific sentinel errors
			if errors.Is(err, utils.ErrNotFound) {
				c.JSON(404, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			if errors.Is(err, utils.ErrForbidden) {
				c.JSON(403, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			if errors.Is(err, utils.ErrBadRequest) {
				c.JSON(400, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			if errors.Is(err, utils.ErrInvalidCredentials) {
				c.JSON(401, gin.H{
					"success": false,
					"error":   "invalid credentials",
				})
				return
			}

			// Default to 500 for unknown errors
			c.JSON(500, gin.H{
				"success": false,
				"error":   "internal server error",
			})
		}
	}
}
