package service

import (
	"context"
	"errors"
	"fmt"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	domainService "todo_app/domain/service"
	"todo_app/internal/dto"

	"github.com/google/uuid"
)

// TodoListServiceImpl implements todo list-related business logic
type TodoListServiceImpl struct {
	listRepo repository.TodoListRepository
	todoRepo repository.TodoRepository
}

// Compile-time check to ensure TodoListServiceImpl implements TodoListService interface
var _ domainService.TodoListService = (*TodoListServiceImpl)(nil)

// NewTodoListService creates a new todo list service
func NewTodoListService(listRepo repository.TodoListRepository, todoRepo repository.TodoRepository) domainService.TodoListService {
	return &TodoListServiceImpl{
		listRepo: listRepo,
		todoRepo: todoRepo,
	}
}

// Create creates a new todo list
func (s *TodoListServiceImpl) Create(ctx context.Context, userID uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error) {
	// Create list entity
	list := entity.NewTodoList(userID, req.Name)

	// Save to database
	if err := s.listRepo.Create(ctx, list); err != nil {
		return nil, err
	}

	response := dto.ListToResponse(list)
	return &response, nil
}

// GetByID retrieves a specific list by ID with its todos
func (s *TodoListServiceImpl) GetByID(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error) {
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, err
	}

	// Authorization check: ensure list belongs to the requesting user
	if !list.BelongsToUser(userID) {
		return nil, errors.New("unauthorized access to this list")
	}

	// Get todos in this list
	listTodos, err := s.todoRepo.FindByListID(ctx, listID)
	if err != nil {
		return nil, err
	}

	response := dto.ListWithTodosToResponse(list, listTodos)
	return &response, nil
}

// List retrieves all lists for a user
func (s *TodoListServiceImpl) List(ctx context.Context, userID uuid.UUID) (*dto.ListsResponse, error) {
	lists, err := s.listRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := dto.ListsToResponse(lists)

	return &dto.ListsResponse{
		Lists: responses,
		Total: len(responses),
	}, nil
}

// Update updates an existing list (rename)
func (s *TodoListServiceImpl) Update(ctx context.Context, listID, userID uuid.UUID, req dto.UpdateListRequest) (*dto.ListResponse, error) {
	// Fetch existing list
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if !list.BelongsToUser(userID) {
		return nil, errors.New("unauthorized access to this list")
	}

	// Update name
	list.UpdateName(req.Name)

	// Save changes
	if err := s.listRepo.Update(ctx, list); err != nil {
		return nil, err
	}

	response := dto.ListToResponse(list)
	return &response, nil
}

// Delete soft deletes a list
func (s *TodoListServiceImpl) Delete(ctx context.Context, listID, userID uuid.UUID) error {
	// Fetch existing list
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return err
	}

	// Authorization check
	if !list.BelongsToUser(userID) {
		return errors.New("unauthorized access to this list")
	}

	// Soft delete the list
	// Note: All todos in this list will be permanently deleted via ON DELETE CASCADE
	return s.listRepo.Delete(ctx, listID)
}

// Duplicate creates a copy of a list with all its todos
func (s *TodoListServiceImpl) Duplicate(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error) {
	// Fetch existing list
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if !list.BelongsToUser(userID) {
		return nil, errors.New("unauthorized access to this list")
	}

	// Get todos in this list
	listTodos, err := s.todoRepo.FindByListID(ctx, listID)
	if err != nil {
		return nil, err
	}

	// Create new list with "(Copy)" suffix
	newName := fmt.Sprintf("%s (Copy)", list.Name)
	newList := entity.NewTodoList(userID, newName)

	// Save new list
	if err := s.listRepo.Create(ctx, newList); err != nil {
		return nil, err
	}

	// Duplicate all todos
	var newTodos []*entity.Todo
	for _, todo := range listTodos {
		// Create new todo with same properties
		newTodo := entity.NewTodo(
			userID,
			todo.Title,
			todo.Description,
			todo.Priority,
			todo.DueDate,
		)
		// Set the new list ID
		newListID := newList.ID
		newTodo.ListID = &newListID

		if err := s.todoRepo.Create(ctx, newTodo); err != nil {
			// Consider transaction rollback here in production
			return nil, err
		}

		newTodos = append(newTodos, newTodo)
	}

	response := dto.ListWithTodosToResponse(newList, newTodos)
	return &response, nil
}
