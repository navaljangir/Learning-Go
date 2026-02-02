package repository

import (
	"context"
	"time"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// TodoFilter defines filtering options for querying todos
type TodoFilter struct {
	UserID    *uuid.UUID
	Completed *bool
	Priority  *entity.Priority
	FromDate  *time.Time
	ToDate    *time.Time
}

// TodoRepository defines the interface for todo data access
type TodoRepository interface {
	// Create creates a new todo
	Create(ctx context.Context, todo *entity.Todo) error

	// FindByID finds a todo by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error)

	// FindByUserID retrieves all todos for a specific user with pagination
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Todo, error)

	// FindWithFilters retrieves todos based on filter criteria with pagination
	FindWithFilters(ctx context.Context, filter TodoFilter, limit, offset int) ([]*entity.Todo, error)

	// Update updates a todo
	Update(ctx context.Context, todo *entity.Todo) error

	// Delete soft deletes a todo
	Delete(ctx context.Context, id uuid.UUID) error

	// Count returns the total count of todos matching the filter
	Count(ctx context.Context, filter TodoFilter) (int64, error)

	// CountByUser returns the total count of todos for a specific user
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}
