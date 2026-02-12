package entity

import "time"

type Todo struct {
	ID          string
	UserID      string
	ListID      *string // Nullable: NULL = global/uncategorized todo
	Title       string
	Description string
	Completed   bool
	Priority    string
	DueDate     *time.Time // ISO8601 format, nullable
	CreatedAt   time.Time  // ISO8601 format
	UpdatedAt   time.Time  // ISO8601 format
	CompletedAt *time.Time // ISO8601 format, nullable
	DeletedAt   *time.Time // ISO8601 format, nullable
}


