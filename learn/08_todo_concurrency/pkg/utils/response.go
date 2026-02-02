package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// RespondError sends an error JSON response
func RespondError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   errorMsg,
	})
}

// RespondCreated sends a 201 Created response
func RespondCreated(c *gin.Context, data interface{}) {
	RespondSuccess(c, http.StatusCreated, data, "Resource created successfully")
}

// RespondOK sends a 200 OK response
func RespondOK(c *gin.Context, data interface{}) {
	RespondSuccess(c, http.StatusOK, data, "")
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(c *gin.Context, errorMsg string) {
	RespondError(c, http.StatusBadRequest, errorMsg)
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(c *gin.Context, errorMsg string) {
	RespondError(c, http.StatusNotFound, errorMsg)
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(c *gin.Context, errorMsg string) {
	RespondError(c, http.StatusInternalServerError, errorMsg)
}
