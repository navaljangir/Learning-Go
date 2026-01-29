package utils

import (
	"gin_server/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: constants.MsgSuccess,
		Data:    data,
	})
}

// Created sends a 201 created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: constants.MsgSuccess,
		Data:    data,
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: constants.MsgBadRequest,
		Error:   message,
	})
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: constants.MsgUnauthorized,
		Error:   message,
	})
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: constants.MsgNotFound,
		Error:   message,
	})
}

// InternalError sends a 500 internal server error response
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: constants.MsgInternalError,
		Error:   message,
	})
}
