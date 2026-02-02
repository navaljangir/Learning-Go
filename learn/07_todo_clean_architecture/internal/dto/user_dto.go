package dto

import (
	"time"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required,max=100"`
}

// UserToResponse converts a user entity to a response DTO
func UserToResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
