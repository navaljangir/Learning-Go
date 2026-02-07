package dto

import (
	"time"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// CreateTodoRequest represents a request to create a new todo
// All fields are optional except title and priority
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=2000"`
	Priority    string     `json:"priority" binding:"required,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`                         // Optional: set due date
	Completed   bool       `json:"completed"`                        // Optional: create as completed (default: false)
	CompletedAt *time.Time `json:"completed_at"`                     // Optional: set completion date (only if completed=true)
	ListID      *string    `json:"list_id" binding:"omitempty,uuid"` // Optional: assign to list
}

// UpdateTodoRequest represents a request to update a todo
// All fields are optional - only provided fields will be updated
type UpdateTodoRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`     // null = remove due date
	Completed   *bool      `json:"completed"`    // Update completion status
	CompletedAt *time.Time `json:"completed_at"` // Update completion date
}

// MoveTodosRequest represents a request to move todos between lists or to global
type MoveTodosRequest struct {
	TodoIDs []string `json:"todo_ids" binding:"required,min=1,dive,uuid"`
	ListID  *string  `json:"list_id" binding:"omitempty,uuid"` // null = move to global
}

// TodoResponse represents a todo in API responses
type TodoResponse struct {
	ID          uuid.UUID  `json:"id"`
	ListID      *uuid.UUID `json:"list_id,omitempty"`   // null = global todo
	ListName    string     `json:"list_name,omitempty"` // Name of the list (empty if global)
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	IsOverdue   bool       `json:"is_overdue"`
}

// TodoListResponse represents a paginated list of todos
type TodoListResponse struct {
	Todos      []TodoResponse `json:"todos"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// TodoStatsResponse represents statistics about user's todos
type TodoStatsResponse struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Pending   int64 `json:"pending"`
	Overdue   int64 `json:"overdue"`
}

// TodoToResponse converts a todo entity to a response DTO
func TodoToResponse(todo *entity.Todo) TodoResponse {
	return TodoResponse{
		ID:          todo.ID,
		ListID:      todo.ListID,
		ListName:    todo.ListName,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		Priority:    string(todo.Priority),
		DueDate:     todo.DueDate,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		CompletedAt: todo.CompletedAt,
		IsOverdue:   todo.IsOverdue(),
	}
}

// TodosToResponse converts multiple todo entities to response DTOs
func TodosToResponse(todos []*entity.Todo) []TodoResponse {
	responses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = TodoToResponse(todo)
	}
	return responses
}
