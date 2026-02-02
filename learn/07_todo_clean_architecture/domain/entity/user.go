package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// NewUser creates a new user with the given details
func NewUser(username, email, passwordHash, fullName string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsDeleted checks if the user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// MarkDeleted marks the user as deleted
func (u *User) MarkDeleted() {
	now := time.Now()
	u.DeletedAt = &now
}

// Update updates user fields
func (u *User) Update(fullName string) {
	u.FullName = fullName
	u.UpdatedAt = time.Now()
}
