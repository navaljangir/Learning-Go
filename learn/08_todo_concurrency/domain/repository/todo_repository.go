package repository

import (
	"context"
	"todo_concurrency/domain/entity"
)

// TodoRepository defines the contract for todo storage operations
//
// KEY LEARNING - INTERFACES:
// This interface allows us to have multiple implementations (in-memory, cache, database)
// without changing the service layer code. This is called "dependency inversion"
// and is a core principle of clean architecture.
//
// Any type that implements these methods satisfies this interface.
type TodoRepository interface {
	// Create adds a new todo to storage
	Create(ctx context.Context, todo *entity.Todo) error

	// FindByID retrieves a todo by its ID
	// Returns nil, nil if not found (not an error)
	FindByID(ctx context.Context, id string) (*entity.Todo, error)

	// FindAll retrieves all todos
	FindAll(ctx context.Context) ([]*entity.Todo, error)

	// Update modifies an existing todo
	Update(ctx context.Context, todo *entity.Todo) error

	// Delete removes a todo by ID
	Delete(ctx context.Context, id string) error

	// Count returns the total number of todos
	Count(ctx context.Context) (int, error)

	// CountCompleted returns the number of completed todos
	CountCompleted(ctx context.Context) (int, error)
}

// StorageInfo provides metadata about the storage implementation
// This interface demonstrates how we can add optional capabilities
type StorageInfo interface {
	// GetStorageType returns the name of the storage backend
	GetStorageType() string

	// GetStats returns storage-specific statistics
	GetStats() map[string]interface{}
}
