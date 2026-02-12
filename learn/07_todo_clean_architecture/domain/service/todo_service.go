package service

import (
	"context"
	"todo_app/dto"

	"github.com/google/uuid"
)

// TodoService defines the interface for todo-related business logic
type TodoService interface {
	// Create creates a new todo item for a user
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error)

	// GetByID retrieves a specific todo by ID
	// Returns error if todo doesn't exist or user is not authorized
	GetByID(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error)

	// List retrieves a paginated list of todos for a user
	List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error)

	// Update updates an existing todo
	// Returns error if todo doesn't exist or user is not authorized
	Update(ctx context.Context, todoID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error)

	// ToggleComplete toggles the completion status of a todo
	// Returns error if todo doesn't exist or user is not authorized
	ToggleComplete(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error)

	// Delete soft deletes a todo
	// Returns error if todo doesn't exist or user is not authorized
	Delete(ctx context.Context, todoID, userID uuid.UUID) error

	// MoveTodos moves multiple todos to a specific list or to global (nil list_id)
	// Returns error if any todo doesn't exist or user is not authorized
	MoveTodos(ctx context.Context, userID uuid.UUID, req dto.MoveTodosRequest) error
}
