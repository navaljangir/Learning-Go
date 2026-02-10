package handler

import "github.com/gin-gonic/gin"

// AuthHandlerInterface defines methods for authentication handlers
// This interface allows for easier testing and dependency injection
type AuthHandlerInterface interface {
	// Register handles user registration
	Register(c *gin.Context)

	// Login handles user authentication
	Login(c *gin.Context)
}

// UserHandlerInterface defines methods for user profile handlers
// This interface allows for easier testing and dependency injection
type UserHandlerInterface interface {
	// GetProfile retrieves the current user's profile
	GetProfile(c *gin.Context)

	// UpdateProfile updates the current user's profile
	UpdateProfile(c *gin.Context)
}

// TodoHandlerInterface defines methods for todo handlers
// This interface allows for easier testing and dependency injection
type TodoHandlerInterface interface {
	// Create creates a new todo
	Create(c *gin.Context)

	// List retrieves todos with pagination
	List(c *gin.Context)

	// GetByID retrieves a specific todo by ID
	GetByID(c *gin.Context)

	// Update updates an existing todo
	Update(c *gin.Context)

	// ToggleComplete toggles the completion status of a todo
	ToggleComplete(c *gin.Context)

	// Delete soft deletes a todo
	Delete(c *gin.Context)

	// MoveTodos moves multiple todos to a list or to global
	MoveTodos(c *gin.Context)
}

// TodoListHandlerInterface defines methods for todo list handlers
// This interface allows for easier testing and dependency injection
type TodoListHandlerInterface interface {
	// Create creates a new list
	Create(c *gin.Context)

	// List retrieves all lists for a user
	List(c *gin.Context)

	// GetByID retrieves a specific list by ID with its todos
	GetByID(c *gin.Context)

	// Update updates a list (rename)
	Update(c *gin.Context)

	// Delete soft deletes a list
	Delete(c *gin.Context)

	// Duplicate creates a copy of a list with all its todos
	Duplicate(c *gin.Context)

	// GenerateShareLink generates a shareable URL token for a list
	GenerateShareLink(c *gin.Context)

	// ImportSharedList imports a shared list via token into the caller's account
	ImportSharedList(c *gin.Context)
}
