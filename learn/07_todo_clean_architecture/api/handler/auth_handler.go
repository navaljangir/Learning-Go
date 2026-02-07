package handler

import (
	"todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"
	"todo_app/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
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
// @Failure 400 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Get user-friendly validation errors
		validationErrors := validator.GetValidationErrors(err)
		if len(validationErrors) > 0 {
			c.JSON(400, gin.H{
				"error":  "Validation failed",
				"fields": validationErrors,
			})
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, response)
}

// Login handles user login
// @Summary Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Get user-friendly validation errors
		validationErrors := validator.GetValidationErrors(err)
		if len(validationErrors) > 0 {
			c.JSON(400, gin.H{
				"error":  "Validation failed",
				"fields": validationErrors,
			})
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, response)
}
