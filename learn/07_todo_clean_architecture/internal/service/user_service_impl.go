package service

import (
	"context"
	"errors"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	domainService "todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"

	"github.com/google/uuid"
)

// UserServiceImpl implements user-related business logic
type UserServiceImpl struct {
	userRepo repository.UserRepository
	jwtUtil  *utils.JWTUtil
}

// Compile-time check to ensure UserServiceImpl implements UserService interface
var _ domainService.UserService = (*UserServiceImpl)(nil)

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) domainService.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

// Register creates a new user account
func (s *UserServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user := entity.NewUser(req.Username, req.Email, hashedPassword, req.FullName)

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtUtil.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      dto.UserToResponse(user),
		ExpiresAt: expiresAt,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is deleted
	if user.IsDeleted() {
		return nil, errors.New("account not found")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtUtil.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      dto.UserToResponse(user),
		ExpiresAt: expiresAt,
	}, nil
}

// GetProfile retrieves the user's profile
func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.IsDeleted() {
		return nil, errors.New("user not found")
	}

	response := dto.UserToResponse(user)
	return &response, nil
}

// UpdateProfile updates the user's profile
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.IsDeleted() {
		return nil, errors.New("user not found")
	}

	// Update user fields
	user.Update(req.FullName)

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := dto.UserToResponse(user)
	return &response, nil
}
