package entity

import (
	"time"

	"github.com/google/uuid"
)

// TodoList represents a list/collection of todos
type TodoList struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewTodoList creates a new todo list
func NewTodoList(userID uuid.UUID, name string) *TodoList {
	now := time.Now()
	return &TodoList{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// BelongsToUser checks if the list belongs to the given user
func (l *TodoList) BelongsToUser(userID uuid.UUID) bool {
	return l.UserID == userID
}

// IsDeleted checks if the list is soft deleted
func (l *TodoList) IsDeleted() bool {
	return l.DeletedAt != nil
}

// MarkDeleted marks the list as deleted
func (l *TodoList) MarkDeleted() {
	now := time.Now()
	l.DeletedAt = &now
}

// UpdateName updates the list name
func (l *TodoList) UpdateName(name string) {
	l.Name = name
	l.UpdatedAt = time.Now()
}
