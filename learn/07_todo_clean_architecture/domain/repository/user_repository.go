package repository

import (
	"context"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByUsername finds a user by username
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a list of users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)

	// ExistsByUsername checks if a username already exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail checks if an email already exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
