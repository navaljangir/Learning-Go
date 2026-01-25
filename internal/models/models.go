// Package models contains data structures and types used throughout the application.
// This includes database models, DTOs (Data Transfer Objects), and other shared types.
package models

// User represents a user in the system.
// Add fields as needed for your application.
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Add more model definitions here as your application grows.
