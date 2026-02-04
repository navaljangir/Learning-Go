package service

import (
	"context"
	"todo_app/internal/dto"

	"github.com/google/uuid"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	// Register creates a new user account and returns authentication token
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error)

	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)

	// GetProfile retrieves the user's profile information
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)

	// UpdateProfile updates the user's profile information
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
}
