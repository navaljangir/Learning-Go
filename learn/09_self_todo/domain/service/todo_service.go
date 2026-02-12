package service

import (
	"context"
	"todo_app/internal/dto"
)

type TodoService interface {
	CreateTodo (ctx context.Context, req *dto.CreateTodoRequest) (*dto.TodoResponse, error)

	// Get Todo by id
	GetTodoByID(ctx context.Context, userID, todoID string) (dto.TodoResponse, error)

	// List Todos for a user
	ListTodos(ctx context.Context, userID string) ([]dto.TodoResponse, error)

	// Update Todo
	UpdateTodo(ctx context.Context, userID , todoID, title, description string) (dto.TodoResponse, error)

	// Delete Todo
	DeleteTodo(ctx context.Context, userID, todoID string) error

	// Get all todos
	GetAllTodos(ctx context.Context) ([]dto.TodoResponse, error)
}