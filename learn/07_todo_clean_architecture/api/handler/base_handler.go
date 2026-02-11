package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response envelope.
// All handlers use this format via BaseHandler.Success() and BaseHandler.Created().
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// BaseHandler provides common response methods for all handlers.
// Embed this in any handler struct to get Success() and Created().
//
// If the response format ever needs to change (e.g., add "timestamp"),
// change it here — every handler inherits the update.
type BaseHandler struct{}

// Success sends a 200 OK response with the standard envelope.
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response with the standard envelope.
func (h *BaseHandler) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}
