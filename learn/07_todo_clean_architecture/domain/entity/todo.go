package entity

import (
	"time"

	"github.com/google/uuid"
)

// Priority represents the priority level of a todo
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Todo represents a todo item in the system
type Todo struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ListID      *uuid.UUID // Nullable: NULL = global/uncategorized todo
	ListName    string     // Name of the list (populated from JOIN, not stored)
	Title       string
	Description string
	Completed   bool
	Priority    Priority
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	DeletedAt   *time.Time
}

// NewTodo creates a new todo with the given details
func NewTodo(userID uuid.UUID, title, description string, priority Priority, dueDate *time.Time) *Todo {
	now := time.Now()
	return &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   false,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsCompleted marks the todo as completed
func (t *Todo) MarkAsCompleted() {
	now := time.Now()
	t.Completed = true
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkAsIncomplete marks the todo as incomplete
func (t *Todo) MarkAsIncomplete() {
	t.Completed = false
	t.CompletedAt = nil
	t.UpdatedAt = time.Now()
}

// IsOverdue checks if the todo is overdue
func (t *Todo) IsOverdue() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// BelongsToUser checks if the todo belongs to the given user
func (t *Todo) BelongsToUser(userID uuid.UUID) bool {
	return t.UserID == userID
}

// IsDeleted checks if the todo is soft deleted
func (t *Todo) IsDeleted() bool {
	return t.DeletedAt != nil
}

// MarkDeleted marks the todo as deleted
func (t *Todo) MarkDeleted() {
	now := time.Now()
	t.DeletedAt = &now
}

// Update updates todo fields
func (t *Todo) Update(title, description string, priority Priority, dueDate *time.Time) {
	t.Title = title
	t.Description = description
	t.Priority = priority
	t.DueDate = dueDate
	t.UpdatedAt = time.Now()
}
