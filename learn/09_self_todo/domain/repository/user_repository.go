package repository

import (
	"context"
	"todo_app/domain/entity"
)

type UserRepository interface {
	// Creating new user
	CreateUser(ctx context.Context, user *entity.User) error

	// Retrieving user by email or ID
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// Retrieving user by ID
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// Updaing user details 
	UpdateUser(ctx context.Context, user *entity.User) error
}