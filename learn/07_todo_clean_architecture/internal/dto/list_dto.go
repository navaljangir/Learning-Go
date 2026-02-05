package dto

import (
	"time"
	"todo_app/domain/entity"

	"github.com/google/uuid"
)

// CreateListRequest represents a request to create a new list
type CreateListRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// UpdateListRequest represents a request to update a list
type UpdateListRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// ListResponse represents a list in API responses
type ListResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	TodoCount int        `json:"todo_count,omitempty"` // Optional: number of todos in list
}

// ListWithTodosResponse represents a list with its todos
type ListWithTodosResponse struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Todos     []TodoResponse `json:"todos"`
}

// ListsResponse represents multiple lists
type ListsResponse struct {
	Lists []ListResponse `json:"lists"`
	Total int            `json:"total"`
}

// ListToResponse converts a list entity to a response DTO
func ListToResponse(list *entity.TodoList) ListResponse {
	return ListResponse{
		ID:        list.ID,
		UserID:    list.UserID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

// ListsToResponse converts multiple list entities to response DTOs
func ListsToResponse(lists []*entity.TodoList) []ListResponse {
	responses := make([]ListResponse, len(lists))
	for i, list := range lists {
		responses[i] = ListToResponse(list)
	}
	return responses
}

// ListWithTodosToResponse converts a list with todos to response DTO
func ListWithTodosToResponse(list *entity.TodoList, todos []*entity.Todo) ListWithTodosResponse {
	return ListWithTodosResponse{
		ID:        list.ID,
		UserID:    list.UserID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
		Todos:     TodosToResponse(todos),
	}
}
