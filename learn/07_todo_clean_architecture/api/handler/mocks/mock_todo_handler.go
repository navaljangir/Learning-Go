package mocks

import "github.com/gin-gonic/gin"

// MockTodoHandler is a fake handler for testing
// ⭐ This implements handler.TodoHandlerInterface
// ⭐ But it doesn't use real database or service!
type MockTodoHandler struct {
	CreateCalled        bool
	ListCalled          bool
	GetByIDCalled       bool
	UpdateCalled        bool
	ToggleCompleteCalled bool
	DeleteCalled        bool

	// Track how many times methods were called
	CreateCount int
	ListCount   int
}

// NewMockTodoHandler creates a new mock todo handler
func NewMockTodoHandler() *MockTodoHandler {
	return &MockTodoHandler{}
}

// Create implements TodoHandlerInterface.Create
// ⭐ Returns fake todo data (no database insert)
func (m *MockTodoHandler) Create(c *gin.Context) {
	m.CreateCalled = true
	m.CreateCount++

	// Return FAKE created todo
	c.JSON(201, gin.H{
		"id":          "mock-todo-id-123",
		"title":       "Mock Todo",
		"description": "This is a mock todo",
		"completed":   false,
		"priority":    "medium",
	})
}

// List implements TodoHandlerInterface.List
// ⭐ Returns fake todo list (no database query)
func (m *MockTodoHandler) List(c *gin.Context) {
	m.ListCalled = true
	m.ListCount++

	// Return FAKE todo list
	c.JSON(200, gin.H{
		"todos": []gin.H{
			{
				"id":          "mock-todo-1",
				"title":       "Mock Todo 1",
				"description": "First mock todo",
				"completed":   false,
			},
			{
				"id":          "mock-todo-2",
				"title":       "Mock Todo 2",
				"description": "Second mock todo",
				"completed":   true,
			},
		},
		"total": 2,
		"page":  1,
		"page_size": 10,
	})
}

// GetByID implements TodoHandlerInterface.GetByID
// ⭐ Returns fake single todo (no database query)
func (m *MockTodoHandler) GetByID(c *gin.Context) {
	m.GetByIDCalled = true

	c.JSON(200, gin.H{
		"id":          "mock-todo-id",
		"title":       "Mock Todo",
		"description": "This is a mock todo",
		"completed":   false,
	})
}

// Update implements TodoHandlerInterface.Update
// ⭐ Returns fake updated todo (no database update)
func (m *MockTodoHandler) Update(c *gin.Context) {
	m.UpdateCalled = true

	c.JSON(200, gin.H{
		"id":          "mock-todo-id",
		"title":       "Updated Mock Todo",
		"description": "This todo was updated",
		"completed":   false,
	})
}

// ToggleComplete implements TodoHandlerInterface.ToggleComplete
// ⭐ Returns fake toggled todo (no database update)
func (m *MockTodoHandler) ToggleComplete(c *gin.Context) {
	m.ToggleCompleteCalled = true

	c.JSON(200, gin.H{
		"id":          "mock-todo-id",
		"title":       "Mock Todo",
		"completed":   true, // Toggled
	})
}

// Delete implements TodoHandlerInterface.Delete
// ⭐ Returns success (no database delete)
func (m *MockTodoHandler) Delete(c *gin.Context) {
	m.DeleteCalled = true

	c.JSON(200, gin.H{
		"message": "todo deleted successfully",
	})
}
