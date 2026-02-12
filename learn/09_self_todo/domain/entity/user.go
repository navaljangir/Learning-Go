package entity

import "time"
type User struct {
	ID string
	Username string
	Email string
	PasswordHash string
	FullName string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}