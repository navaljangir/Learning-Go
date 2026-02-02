package entity

import "time"

// Todo represents a task in our system
// This is a domain entity - pure business logic with no dependencies
type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    int       `json:"priority"` // 1=Low, 2=Medium, 3=High
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsValid checks if the todo has required fields
func (t *Todo) IsValid() bool {
	return t.Title != "" && t.Priority >= 1 && t.Priority <= 3
}

// MarkComplete marks the todo as completed
func (t *Todo) MarkComplete() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkIncomplete marks the todo as not completed
func (t *Todo) MarkIncomplete() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}
