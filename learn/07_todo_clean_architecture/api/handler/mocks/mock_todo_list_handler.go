package mocks

import "github.com/gin-gonic/gin"

// MockTodoListHandler is a fake handler for testing list operations
type MockTodoListHandler struct {
	CreateCalled          bool
	ListCalled            bool
	GetByIDCalled         bool
	UpdateCalled          bool
	DeleteCalled          bool
	DuplicateCalled       bool
	GenerateShareCalled   bool
	ImportSharedCalled    bool

	CreateCount          int
	ListCount            int
	DuplicateCount       int
	GenerateShareCount   int
	ImportSharedCount    int
}

// NewMockTodoListHandler creates a new mock todo list handler
func NewMockTodoListHandler() *MockTodoListHandler {
	return &MockTodoListHandler{}
}

// Create implements TodoListHandlerInterface.Create
func (m *MockTodoListHandler) Create(c *gin.Context) {
	m.CreateCalled = true
	m.CreateCount++

	c.JSON(201, gin.H{
		"id":      "mock-list-id-123",
		"name":    "Mock List",
		"user_id": "mock-user-id",
	})
}

// List implements TodoListHandlerInterface.List
func (m *MockTodoListHandler) List(c *gin.Context) {
	m.ListCalled = true
	m.ListCount++

	c.JSON(200, gin.H{
		"lists": []gin.H{
			{
				"id":   "mock-list-1",
				"name": "Work",
			},
			{
				"id":   "mock-list-2",
				"name": "Personal",
			},
		},
		"total": 2,
	})
}

// GetByID implements TodoListHandlerInterface.GetByID
func (m *MockTodoListHandler) GetByID(c *gin.Context) {
	m.GetByIDCalled = true

	c.JSON(200, gin.H{
		"id":    "mock-list-id",
		"name":  "Mock List",
		"todos": []gin.H{},
	})
}

// Update implements TodoListHandlerInterface.Update
func (m *MockTodoListHandler) Update(c *gin.Context) {
	m.UpdateCalled = true

	c.JSON(200, gin.H{
		"id":   "mock-list-id",
		"name": "Updated List Name",
	})
}

// Delete implements TodoListHandlerInterface.Delete
func (m *MockTodoListHandler) Delete(c *gin.Context) {
	m.DeleteCalled = true

	c.JSON(200, gin.H{
		"message": "list deleted successfully",
	})
}

// Duplicate implements TodoListHandlerInterface.Duplicate
func (m *MockTodoListHandler) Duplicate(c *gin.Context) {
	m.DuplicateCalled = true
	m.DuplicateCount++

	c.JSON(201, gin.H{
		"id":    "mock-list-id-copy",
		"name":  "Mock List (Copy)",
		"todos": []gin.H{},
	})
}

// GenerateShareLink implements TodoListHandlerInterface.GenerateShareLink
func (m *MockTodoListHandler) GenerateShareLink(c *gin.Context) {
	m.GenerateShareCalled = true
	m.GenerateShareCount++

	c.JSON(200, gin.H{
		"share_url":   "/api/v1/lists/import/abc123token",
		"share_token": "abc123token",
	})
}

// ImportSharedList implements TodoListHandlerInterface.ImportSharedList
func (m *MockTodoListHandler) ImportSharedList(c *gin.Context) {
	m.ImportSharedCalled = true
	m.ImportSharedCount++

	c.JSON(201, gin.H{
		"id":    "mock-imported-list-id",
		"name":  "Mock List (shared)",
		"todos": []gin.H{},
	})
}
