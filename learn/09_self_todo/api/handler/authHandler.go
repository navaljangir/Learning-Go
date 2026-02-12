package handler

import (
	"todo_app/domain/service"
	"todo_app/internal/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService service.AuthService
}

func NewAuthHandler(userService service.AuthService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Implementation for user login
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	} 
}