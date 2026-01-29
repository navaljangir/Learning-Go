package handlers

import (
	"gin_server/constants"
	"gin_server/models"
	"gin_server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// In-memory user storage (replace with database in production)
var users = make(map[string]*models.User) // key: username

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Check if user already exists
	if _, exists := users[req.Username]; exists {
		utils.BadRequest(c, constants.MsgUserExists)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	// Create user
	user := &models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user
	users[req.Username] = user

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	// Return response
	utils.Created(c, models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req models.LoginRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Find user
	user, exists := users[req.Username]
	if !exists {
		utils.Unauthorized(c, constants.MsgInvalidCredentials)
		return
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.Unauthorized(c, constants.MsgInvalidCredentials)
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	// Return response
	utils.Success(c, models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}
