package handlers

import (
	"gin_server/constants"
	"gin_server/models"
	"gin_server/utils"

	"github.com/gin-gonic/gin"
)

// GetProfile returns the authenticated user's profile
func GetProfile(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	username, exists := c.Get(constants.ContextUsername)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	// Find user
	user, exists := users[username.(string)]
	if !exists {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user.ToResponse())
}

// GetAllUsers returns all users (for demo purposes)
func GetAllUsers(c *gin.Context) {
	var userList []models.UserResponse

	for _, user := range users {
		userList = append(userList, user.ToResponse())
	}

	utils.Success(c, userList)
}
