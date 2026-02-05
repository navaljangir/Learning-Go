package service

import (
	"context"
	"errors"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	domainService "todo_app/domain/service"
	"todo_app/internal/dto"

	"github.com/google/uuid"
)

// TodoServiceImpl implements todo-related business logic
type TodoServiceImpl struct {
	todoRepo repository.TodoRepository
}

// Compile-time check to ensure TodoServiceImpl implements TodoService interface
var _ domainService.TodoService = (*TodoServiceImpl)(nil)

// NewTodoService creates a new todo service
func NewTodoService(todoRepo repository.TodoRepository) domainService.TodoService {
	return &TodoServiceImpl{todoRepo: todoRepo}
}

// Create creates a new todo
func (s *TodoServiceImpl) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
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
func (s *TodoServiceImpl) GetByID(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
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
func (s *TodoServiceImpl) List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error) {
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
func (s *TodoServiceImpl) Update(ctx context.Context, todoID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
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
func (s *TodoServiceImpl) ToggleComplete(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
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
func (s *TodoServiceImpl) Delete(ctx context.Context, todoID, userID uuid.UUID) error {
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

// MoveTodos moves multiple todos to a specific list or to global (nil list_id)
func (s *TodoServiceImpl) MoveTodos(ctx context.Context, userID uuid.UUID, req dto.MoveTodosRequest) error {
	// Validate input
	if len(req.TodoIDs) == 0 {
		return errors.New("no todos specified")
	}

	// Convert string IDs to UUIDs
	todoIDs := make([]uuid.UUID, len(req.TodoIDs))
	for i, idStr := range req.TodoIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return errors.New("invalid todo ID format")
		}
		todoIDs[i] = id
	}

	// Authorization check: verify all todos belong to the user
	// This is important for security - don't move todos the user doesn't own
	for _, todoID := range todoIDs {
		todo, err := s.todoRepo.FindByID(ctx, todoID)
		if err != nil {
			return errors.New("one or more todos not found")
		}
		if !todo.BelongsToUser(userID) {
			return errors.New("unauthorized access to one or more todos")
		}
	}

	// Convert list_id string to UUID pointer
	var listIDPtr *uuid.UUID
	if req.ListID != nil {
		listID, err := uuid.Parse(*req.ListID)
		if err != nil {
			return errors.New("invalid list ID format")
		}
		listIDPtr = &listID
		// TODO: In Phase 2, verify the list exists and belongs to the user
	}

	// Perform the bulk update
	return s.todoRepo.UpdateListID(ctx, todoIDs, listIDPtr, userID)
}
