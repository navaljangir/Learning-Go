package middleware

import (
	"net/http"
	"strings"
	"todo_app/pkg/constants"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(&utils.AppError{Err: utils.ErrInvalidCredentials, Message: "authorization token required", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(&utils.AppError{Err: utils.ErrInvalidCredentials, Message: "invalid authorization header format", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			c.Error(&utils.AppError{Err: utils.ErrInvalidCredentials, Message: "invalid or expired token", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.Error(&utils.AppError{Err: utils.ErrInvalidCredentials, Message: "invalid user ID in token", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(constants.ContextUserID, userID)
		c.Set(constants.ContextUsername, claims.Username)

		c.Next()
	}
}
