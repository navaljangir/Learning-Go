package dto

import (
	"time"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// CreateTodoRequest represents a request to create a new todo
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=2000"`
	Priority    string     `json:"priority" binding:"required,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTodoRequest represents a request to update a todo
type UpdateTodoRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
}

// TodoResponse represents a todo in API responses
type TodoResponse struct {
	ID          uuid.UUID  `json:"id"`
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
