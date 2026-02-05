package mocks

import "github.com/gin-gonic/gin"

// MockUserHandler is a fake handler for testing
// This implements handler.UserHandlerInterface
// But it doesn't use real database or service!
type MockUserHandler struct {
	GetProfileCalled    bool
	UpdateProfileCalled bool
}

// NewMockUserHandler creates a new mock user handler
func NewMockUserHandler() *MockUserHandler {
	return &MockUserHandler{}
}

// GetProfile implements UserHandlerInterface.GetProfile
// Returns fake data immediately (no database query)
func (m *MockUserHandler) GetProfile(c *gin.Context) {
	m.GetProfileCalled = true

	// Return FAKE user data
	c.JSON(200, gin.H{
		"id":        "mock-user-id",
		"username":  "mock-user",
		"email":     "mock@example.com",
		"full_name": "Mock User",
	})
}

// UpdateProfile implements UserHandlerInterface.UpdateProfile
// Returns fake data immediately (no database update)
func (m *MockUserHandler) UpdateProfile(c *gin.Context) {
	m.UpdateProfileCalled = true

	// Return FAKE updated user data
	c.JSON(200, gin.H{
		"id":        "mock-user-id",
		"username":  "updated-user",
		"email":     "updated@example.com",
		"full_name": "Updated Mock User",
	})
}
