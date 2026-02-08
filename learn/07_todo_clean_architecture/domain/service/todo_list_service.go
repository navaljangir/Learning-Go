package service

import (
	"context"
	"todo_app/internal/dto"

	"github.com/google/uuid"
)

// TodoListService defines the interface for todo list-related business logic
type TodoListService interface {
	// Create creates a new todo list for a user
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error)

	// GetByID retrieves a specific list by ID with its todos
	// Returns error if list doesn't exist or user is not authorized
	GetByID(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error)

	// List retrieves all lists for a user
	List(ctx context.Context, userID uuid.UUID) (*dto.ListsResponse, error)

	// Update updates an existing list (rename)
	// Returns error if list doesn't exist or user is not authorized
	Update(ctx context.Context, listID, userID uuid.UUID, req dto.UpdateListRequest) (*dto.ListResponse, error)

	// Delete soft deletes a list and permanently deletes all its todos (CASCADE)
	// Returns error if list doesn't exist or user is not authorized
	Delete(ctx context.Context, listID, userID uuid.UUID) error

	// Duplicate creates a copy of a list with all its todos
	// Returns error if list doesn't exist or user is not authorized
	Duplicate(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error)

	// Share creates a copy of a list with all its todos for a different user
	// This allows users to share their lists with others by creating independent copies
	// Returns error if list doesn't exist, user is not authorized, or target user doesn't exist
	Share(ctx context.Context, listID, ownerUserID, targetUserID uuid.UUID, req dto.ShareListRequest) (*dto.ListWithTodosResponse, error)
}
