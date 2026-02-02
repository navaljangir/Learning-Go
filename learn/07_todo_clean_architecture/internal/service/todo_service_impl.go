package service

import (
	"context"
	"errors"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	"todo_app/internal/dto"

	"github.com/google/uuid"
)

// TodoService implements todo-related business logic
type TodoService struct {
	todoRepo repository.TodoRepository
}

// NewTodoService creates a new todo service
func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

// Create creates a new todo
func (s *TodoService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
	// Convert priority string to entity type
	priority := entity.Priority(req.Priority)

	// Create todo entity
	todo := entity.NewTodo(userID, req.Title, req.Description, priority, req.DueDate)

	// Save to database
	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// GetByID retrieves a specific todo by ID
func (s *TodoService) GetByID(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// Authorization check: ensure todo belongs to the requesting user
	if !todo.BelongsToUser(userID) {
		return nil, errors.New("unauthorized access to this todo")
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// List retrieves a paginated list of todos for a user
func (s *TodoService) List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Fetch todos from repository
	todos, err := s.todoRepo.FindByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	total, err := s.todoRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert entities to DTOs
	responses := dto.TodosToResponse(todos)

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &dto.TodoListResponse{
		Todos:      responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates an existing todo
func (s *TodoService) Update(ctx context.Context, todoID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
	// Fetch existing todo
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if !todo.BelongsToUser(userID) {
		return nil, errors.New("unauthorized access to this todo")
	}

	// Update fields if provided
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != nil {
		todo.Priority = entity.Priority(*req.Priority)
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}
	todo.UpdatedAt = time.Now()

	// Save changes
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// ToggleComplete toggles the completion status of a todo
func (s *TodoService) ToggleComplete(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
	// Fetch existing todo
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if !todo.BelongsToUser(userID) {
		return nil, errors.New("unauthorized access to this todo")
	}

	// Toggle completion using domain logic
	if todo.Completed {
		todo.MarkAsIncomplete()
	} else {
		todo.MarkAsCompleted()
	}

	// Save changes
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// Delete soft deletes a todo
func (s *TodoService) Delete(ctx context.Context, todoID, userID uuid.UUID) error {
	// Fetch existing todo
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return err
	}

	// Authorization check
	if !todo.BelongsToUser(userID) {
		return errors.New("unauthorized access to this todo")
	}

	// Soft delete the todo
	return s.todoRepo.Delete(ctx, todoID)
}
