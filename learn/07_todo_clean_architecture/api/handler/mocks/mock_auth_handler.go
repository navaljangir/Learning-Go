package mocks

import (
	"github.com/gin-gonic/gin"
)

// MockAuthHandler is a fake handler for testing
// It implements handler.AuthHandlerInterface
type MockAuthHandler struct {
	// Track if methods were called
	RegisterCalled bool
	LoginCalled    bool

	// Control what the mock returns
	ShouldReturnError bool
	RegisterCount     int
	LoginCount        int
}

// NewMockAuthHandler creates a new mock handler
func NewMockAuthHandler() *MockAuthHandler {
	return &MockAuthHandler{}
}

// Register implements AuthHandlerInterface.Register
// This is a FAKE implementation for testing
func (m *MockAuthHandler) Register(c *gin.Context) {
	m.RegisterCalled = true
	m.RegisterCount++

	if m.ShouldReturnError {
		c.JSON(400, gin.H{
			"error": "mock error - registration failed",
		})
		return
	}

	// Return fake success response
	c.JSON(201, gin.H{
		"token": "mock-token-12345",
		"user": gin.H{
			"id":       "mock-user-id",
			"username": "mock-username",
			"email":    "mock@example.com",
		},
	})
}

// Login implements AuthHandlerInterface.Login
// This is a FAKE implementation for testing
func (m *MockAuthHandler) Login(c *gin.Context) {
	m.LoginCalled = true
	m.LoginCount++

	if m.ShouldReturnError {
		c.JSON(401, gin.H{
			"error": "mock error - invalid credentials",
		})
		return
	}

	// Return fake success response
	c.JSON(200, gin.H{
		"token": "mock-token-67890",
		"user": gin.H{
			"id":       "mock-user-id",
			"username": "mock-username",
		},
	})
}

// Reset clears all tracking data
func (m *MockAuthHandler) Reset() {
	m.RegisterCalled = false
	m.LoginCalled = false
	m.ShouldReturnError = false
	m.RegisterCount = 0
	m.LoginCount = 0
}
