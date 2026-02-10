package service

import (
	"context"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	domainService "todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"

	"github.com/google/uuid"
)

// TodoServiceImpl implements todo-related business logic
type TodoServiceImpl struct {
	todoRepo repository.TodoRepository
	listRepo repository.TodoListRepository
}

// Compile-time check to ensure TodoServiceImpl implements TodoService interface
var _ domainService.TodoService = (*TodoServiceImpl)(nil)

// NewTodoService creates a new todo service
func NewTodoService(todoRepo repository.TodoRepository, listRepo repository.TodoListRepository) domainService.TodoService {
	return &TodoServiceImpl{todoRepo: todoRepo, listRepo: listRepo}
}

// Create creates a new todo
func (s *TodoServiceImpl) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
	// Convert priority string to entity type
	priority := entity.Priority(req.Priority)

	// Create todo entity
	todo := entity.NewTodo(userID, req.Title, req.Description, priority, req.DueDate)

	// Handle optional completion status
	if req.Completed {
		// If creating as completed, set completed status
		if req.CompletedAt != nil {
			// Validate: completed_at must not be in the future
			if req.CompletedAt.After(time.Now()) {
				return nil, &utils.AppError{
					Err:        utils.ErrBadRequest,
					Message:    "completed_at cannot be in the future",
					StatusCode: 400,
				}
			}
			// Use provided completion date
			todo.Completed = true
			todo.CompletedAt = req.CompletedAt
		} else {
			// Use current time as completion date
			todo.MarkAsCompleted()
		}
	}

	// Handle optional list assignment
	if req.ListID != nil {
		listID, err := uuid.Parse(*req.ListID)
		if err != nil {
			return nil, &utils.AppError{
				Err:        utils.ErrBadRequest,
				Message:    "Invalid list ID format",
				StatusCode: 400,
			}
		}

		// Verify list exists and belongs to user
		// If list doesn't exist or belongs to another user, create as global todo
		list, err := s.listRepo.FindByID(ctx, listID)
		if err == nil && list.BelongsToUser(userID) {
			todo.ListID = &listID
		}
		// else: silently skip — todo is created without a list (global)
	}

	// Save to database
	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to create todo",
			StatusCode: 500,
		}
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// GetByID retrieves a specific todo by ID
func (s *TodoServiceImpl) GetByID(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "Todo not found",
			StatusCode: 404,
		}
	}

	// Authorization check: ensure todo belongs to the requesting user
	if !todo.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this todo",
			StatusCode: 403,
		}
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
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to fetch todos",
			StatusCode: 500,
		}
	}

	// Get total count for pagination
	total, err := s.todoRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to count todos",
			StatusCode: 500,
		}
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
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "Todo not found",
			StatusCode: 404,
		}
	}

	// Authorization check
	if !todo.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this todo",
			StatusCode: 403,
		}
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

	// Handle completion status update
	if req.Completed != nil {
		if *req.Completed {
			// Mark as completed
			if req.CompletedAt != nil {
				// Validate: completed_at must not be in the future
				if req.CompletedAt.After(time.Now()) {
					return nil, &utils.AppError{
						Err:        utils.ErrBadRequest,
						Message:    "completed_at cannot be in the future",
						StatusCode: 400,
					}
				}
				// Use provided completion date
				todo.Completed = true
				todo.CompletedAt = req.CompletedAt
			} else {
				// Use current time
				todo.MarkAsCompleted()
			}
		} else {
			// Mark as incomplete
			todo.MarkAsIncomplete()
		}
	}

	todo.UpdatedAt = time.Now()

	// Save changes
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to update todo",
			StatusCode: 500,
		}
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// ToggleComplete toggles the completion status of a todo
func (s *TodoServiceImpl) ToggleComplete(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
	// Fetch existing todo
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "Todo not found",
			StatusCode: 404,
		}
	}

	// Authorization check
	if !todo.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this todo",
			StatusCode: 403,
		}
	}

	// Toggle completion using domain logic
	if todo.Completed {
		todo.MarkAsIncomplete()
	} else {
		todo.MarkAsCompleted()
	}

	// Save changes
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to update todo",
			StatusCode: 500,
		}
	}

	response := dto.TodoToResponse(todo)
	return &response, nil
}

// Delete soft deletes a todo
func (s *TodoServiceImpl) Delete(ctx context.Context, todoID, userID uuid.UUID) error {
	// Fetch existing todo
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "Todo not found",
			StatusCode: 404,
		}
	}

	// Authorization check
	if !todo.BelongsToUser(userID) {
		return &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this todo",
			StatusCode: 403,
		}
	}

	// Soft delete the todo
	if err := s.todoRepo.Delete(ctx, todoID); err != nil {
		return &utils.AppError{
			Err:        err,
			Message:    "Failed to delete todo",
			StatusCode: 500,
		}
	}
	return nil
}

// MoveTodos moves multiple todos to a specific list or to global (nil list_id)
func (s *TodoServiceImpl) MoveTodos(ctx context.Context, userID uuid.UUID, req dto.MoveTodosRequest) error {
	// Validate input
	if len(req.TodoIDs) == 0 {
		return &utils.AppError{
			Err:        utils.ErrBadRequest,
			Message:    "No todos specified",
			StatusCode: 400,
		}
	}

	// Convert string IDs to UUIDs
	todoIDs := make([]uuid.UUID, len(req.TodoIDs))
	for i, idStr := range req.TodoIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return &utils.AppError{
				Err:        utils.ErrBadRequest,
				Message:    "Invalid todo ID format",
				StatusCode: 400,
			}
		}
		todoIDs[i] = id
	}

	// Authorization check: verify all todos belong to the user
	// This is important for security - don't move todos the user doesn't own
	for _, todoID := range todoIDs {
		todo, err := s.todoRepo.FindByID(ctx, todoID)
		if err != nil {
			return &utils.AppError{
				Err:        utils.ErrNotFound,
				Message:    "One or more todos not found",
				StatusCode: 404,
			}
		}
		if !todo.BelongsToUser(userID) {
			return &utils.AppError{
				Err:        utils.ErrForbidden,
				Message:    "Unauthorized access to one or more todos",
				StatusCode: 403,
			}
		}
	}

	// Convert list_id string to UUID pointer
	var listIDPtr *uuid.UUID
	if req.ListID != nil {
		listID, err := uuid.Parse(*req.ListID)
		if err != nil {
			return &utils.AppError{
				Err:        utils.ErrBadRequest,
				Message:    "Invalid list ID format",
				StatusCode: 400,
			}
		}
		listIDPtr = &listID
		// TODO: In Phase 2, verify the list exists and belongs to the user
	}

	// Perform the bulk update
	if err := s.todoRepo.UpdateListID(ctx, todoIDs, listIDPtr, userID); err != nil {
		return &utils.AppError{
			Err:        err,
			Message:    "Failed to move todos",
			StatusCode: 500,
		}
	}
	return nil
}
