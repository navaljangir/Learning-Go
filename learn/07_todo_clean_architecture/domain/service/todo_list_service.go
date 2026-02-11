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
	// If req.KeepCompleted is true, completed status is preserved; otherwise all todos start incomplete
	// Returns error if list doesn't exist or user is not authorized
	Duplicate(ctx context.Context, listID, userID uuid.UUID, req dto.DuplicateListRequest) (*dto.ListWithTodosResponse, error)

	// GenerateShareLink generates a shareable URL with a random token for a list
	// If the list already has a token, returns the existing one (idempotent)
	// Returns error if list doesn't exist or user is not the owner
	GenerateShareLink(ctx context.Context, listID, userID uuid.UUID) (*dto.ShareLinkResponse, error)

	// ImportSharedList imports a shared list (and its todos) into the caller's account
	// Creates a new independent copy of the list and all its todos
	// If req.KeepCompleted is true, completed status is preserved; otherwise all todos start incomplete
	// Returns error if token is invalid or caller is the list owner
	ImportSharedList(ctx context.Context, token string, userID uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error)
}
