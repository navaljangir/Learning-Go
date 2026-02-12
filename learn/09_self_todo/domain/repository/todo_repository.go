package repository

import (
	"context"
	"todo_app/domain/entity"
)

type TodoRepository interface {

	// CreateTodo creates a new todo item in the database
	CreateTodo( c context.Context,  todo *entity.Todo) (*entity.Todo, error)
	// FindById retrieves a todo item by its ID
	FindById (c context.Context, id string) (*entity.Todo, error)
	// FindByUserId retrieves all todo items for a specific user
	FindByUserId(c context.Context, userId string) ([]*entity.Todo, error)

	// UpdateTodo updates an existing todo item in the database
	UpdateTodo(c context.Context, todo *entity.Todo) (*entity.Todo, error)
	DeleteTodo(c context.Context, id string) error
	GetAllTodos(c context.Context) ([]*entity.Todo, error)
}