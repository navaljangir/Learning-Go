package middleware

import (
	"errors"
	"todo_app/pkg/utils"
	"todo_app/pkg/validator"

	"github.com/gin-gonic/gin"
	govalidator "github.com/go-playground/validator/v10"
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

			// Check if it's a ValidationError (custom per-field validation)
			var valErr *utils.ValidationError
			if errors.As(err, &valErr) {
				c.JSON(400, gin.H{
					"success": false,
					"error":   valErr.Message,
					"fields":  valErr.Fields,
				})
				return
			}

			// Check if it's Gin's validator.ValidationErrors (from ShouldBindJSON)
			var validationErrors govalidator.ValidationErrors
			if errors.As(err, &validationErrors) {
				fields := validator.GetValidationErrors(err)
				c.JSON(400, gin.H{
					"success": false,
					"error":   "Validation failed",
					"fields":  fields,
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

			// Any other binding/request errors → 400
			// This catches malformed JSON, type mismatches, etc.
			if isBadRequestError(err) {
				c.JSON(400, gin.H{
					"success": false,
					"error":   "Invalid request body",
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

// isBadRequestError checks if the error looks like a client-side binding error.
// Gin's ShouldBindJSON returns errors with specific types for JSON syntax issues.
func isBadRequestError(err error) bool {
	// Check for common encoding/json error types that indicate bad input
	errMsg := err.Error()
	for _, prefix := range []string{
		"invalid character",  // json.SyntaxError
		"unexpected end of",  // json.SyntaxError
		"cannot unmarshal",   // json.UnmarshalTypeError
		"json: cannot",       // json.UnmarshalTypeError
		"EOF",                // empty body
		"invalid request",    // generic bad input
	} {
		if len(errMsg) >= len(prefix) && errMsg[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
