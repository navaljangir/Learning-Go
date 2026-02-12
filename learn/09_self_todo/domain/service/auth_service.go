package service

import (
	"context"
	"todo_app/internal/dto"
)

type AuthService interface {
	// Register a new user
	Register(ctx context.Context, req *dto.RegisterRequest) (dto.RegisterResponse, error)

	// Login an existing user and return a JWT token
	Login(ctx context.Context, req *dto.LoginRequest) (dto.LoginResponse, error)
}