package middlewares

import (
	"gin_server/constants"
	"gin_server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, constants.MsgTokenRequired)
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization header format. Use: Bearer <token>")
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, constants.MsgTokenInvalid)
			c.Abort()
			return
		}

		// Set user info in context for handlers to use
		c.Set(constants.ContextUserID, claims.UserID)
		c.Set(constants.ContextUsername, claims.Username)

		c.Next()
	}
}
