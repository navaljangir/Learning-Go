package repository

import (
	"context"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// TodoListRepository defines the interface for todo list data access
type TodoListRepository interface {
	// Create creates a new todo list
	Create(ctx context.Context, list *entity.TodoList) error

	// FindByID finds a list by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.TodoList, error)

	// FindByUserID retrieves all lists for a specific user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.TodoList, error)

	// Update updates a list
	Update(ctx context.Context, list *entity.TodoList) error

	// Delete soft deletes a list and permanently deletes all its todos (CASCADE)
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByUser returns the total count of lists for a specific user
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	// Duplicate creates a copy of a list (without todos)
	// Use this with todo operations to duplicate list + todos
	Duplicate(ctx context.Context, sourceListID uuid.UUID, newName string) (*entity.TodoList, error)
}
