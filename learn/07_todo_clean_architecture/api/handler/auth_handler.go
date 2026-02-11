package handler

import (
	"todo_app/domain/service"
	"todo_app/internal/dto"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	BaseHandler
	userService service.UserService
}

// NewAuthHandler creates a new auth handler
// Returns AuthHandlerInterface to enforce dependency on abstraction
func NewAuthHandler(userService service.UserService) AuthHandlerInterface {

	return &AuthHandler{userService: userService}
}

// Register handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} handler.Response	
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		c.Error(err) // Middleware will handle AppError
		return
	}

	h.Created(c, response)
}

// Login handles user login
// @Summary Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} handler.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err) // Middleware will handle AppError
		return
	}

	h.Success(c, response)
}
