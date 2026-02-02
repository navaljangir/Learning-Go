package service

import (
	"context"
	"errors"
	"todo_concurrency/domain/entity"
	"todo_concurrency/domain/repository"
	"todo_concurrency/internal/dto"
)

// TodoService handles business logic for todos
//
// KEY LEARNING - INTERFACES:
// The service depends on the INTERFACE, not a concrete implementation.
// This means we can swap storage backends without changing this code!
type TodoService struct {
	repo repository.TodoRepository // Interface, not concrete type!
}

// NewTodoService creates a new service
func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// Create creates a new todo
func (s *TodoService) Create(ctx context.Context, req dto.CreateTodoRequest) (*entity.Todo, error) {
	todo := &entity.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Completed:   false,
	}

	if !todo.IsValid() {
		return nil, errors.New("invalid todo data")
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// GetByID retrieves a todo by ID
func (s *TodoService) GetByID(ctx context.Context, id string) (*entity.Todo, error) {
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// GetAll retrieves all todos
func (s *TodoService) GetAll(ctx context.Context) ([]*entity.Todo, error) {
	return s.repo.FindAll(ctx)
}

// Update updates an existing todo
func (s *TodoService) Update(ctx context.Context, id string, req dto.UpdateTodoRequest) (*entity.Todo, error) {
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}

	// Update fields if provided
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Priority > 0 {
		todo.Priority = req.Priority
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// Delete deletes a todo
func (s *TodoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ToggleComplete toggles the completion status
func (s *TodoService) ToggleComplete(ctx context.Context, id string) (*entity.Todo, error) {
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}

	if todo.Completed {
		todo.MarkIncomplete()
	} else {
		todo.MarkComplete()
	}

	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// GetRepository returns the underlying repository
// This allows handlers to access storage-specific features like stats
func (s *TodoService) GetRepository() repository.TodoRepository {
	return s.repo
}
