package handler

import (
	"todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	BaseHandler
	userService service.UserService
}

// NewUserHandler creates a new user handler
// Returns UserHandlerInterface to enforce dependency on abstraction
func NewUserHandler(userService service.UserService) UserHandlerInterface {
	return &UserHandler{userService: userService}
}

// GetProfile handles getting the current user's profile
// @Summary Get user profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} handler.Response
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	response, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	h.Success(c, response)
}

// UpdateProfile handles updating the current user's profile
// @Summary Update user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "Profile update details"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} handler.Response
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	h.Success(c, response)
}
