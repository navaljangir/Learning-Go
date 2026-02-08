package service

import (
	"context"
	"fmt"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	domainService "todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"

	"github.com/google/uuid"
)

// TodoListServiceImpl implements todo list-related business logic
type TodoListServiceImpl struct {
	listRepo repository.TodoListRepository
	todoRepo repository.TodoRepository
	userRepo repository.UserRepository
}

// Compile-time check to ensure TodoListServiceImpl implements TodoListService interface
var _ domainService.TodoListService = (*TodoListServiceImpl)(nil)

// NewTodoListService creates a new todo list service
func NewTodoListService(
	listRepo repository.TodoListRepository,
	todoRepo repository.TodoRepository,
	userRepo repository.UserRepository,
) domainService.TodoListService {
	return &TodoListServiceImpl{
		listRepo: listRepo,
		todoRepo: todoRepo,
		userRepo: userRepo,
	}
}

// Create creates a new todo list
func (s *TodoListServiceImpl) Create(ctx context.Context, userID uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error) {
	// Create list entity
	list := entity.NewTodoList(userID, req.Name)

	// Save to database
	if err := s.listRepo.Create(ctx, list); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to create list",
			StatusCode: 500,
		}
	}

	response := dto.ListToResponse(list)
	return &response, nil
}

// GetByID retrieves a specific list by ID with its todos
func (s *TodoListServiceImpl) GetByID(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error) {
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "List not found",
			StatusCode: 404,
		}
	}

	// Authorization check: ensure list belongs to the requesting user
	if !list.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this list",
			StatusCode: 403,
		}
	}

	// Get todos in this list
	listTodos, err := s.todoRepo.FindByListID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to fetch todos",
			StatusCode: 500,
		}
	}

	response := dto.ListWithTodosToResponse(list, listTodos)
	return &response, nil
}

// List retrieves all lists for a user
func (s *TodoListServiceImpl) List(ctx context.Context, userID uuid.UUID) (*dto.ListsResponse, error) {
	lists, err := s.listRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to fetch lists",
			StatusCode: 500,
		}
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
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "List not found",
			StatusCode: 404,
		}
	}

	// Authorization check
	if !list.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this list",
			StatusCode: 403,
		}
	}

	// Update name
	list.UpdateName(req.Name)

	// Save changes
	if err := s.listRepo.Update(ctx, list); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to update list",
			StatusCode: 500,
		}
	}

	response := dto.ListToResponse(list)
	return &response, nil
}

// Delete soft deletes a list
func (s *TodoListServiceImpl) Delete(ctx context.Context, listID, userID uuid.UUID) error {
	// Fetch existing list
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "List not found",
			StatusCode: 404,
		}
	}

	// Authorization check
	if !list.BelongsToUser(userID) {
		return &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this list",
			StatusCode: 403,
		}
	}

	// Soft delete the list
	// Note: All todos in this list will be permanently deleted via ON DELETE CASCADE
	if err := s.listRepo.Delete(ctx, listID); err != nil {
		return &utils.AppError{
			Err:        err,
			Message:    "Failed to delete list",
			StatusCode: 500,
		}
	}
	return nil
}

// Duplicate creates a copy of a list with all its todos
func (s *TodoListServiceImpl) Duplicate(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error) {
	// Fetch existing list
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "List not found",
			StatusCode: 404,
		}
	}

	// Authorization check
	if !list.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this list",
			StatusCode: 403,
		}
	}

	// Get todos in this list
	listTodos, err := s.todoRepo.FindByListID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to fetch todos",
			StatusCode: 500,
		}
	}

	// Create new list with "(Copy)" suffix
	newName := fmt.Sprintf("%s (Copy)", list.Name)
	newList := entity.NewTodoList(userID, newName)

	// Save new list
	if err := s.listRepo.Create(ctx, newList); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to create duplicate list",
			StatusCode: 500,
		}
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
			return nil, &utils.AppError{
				Err:        err,
				Message:    "Failed to duplicate todos",
				StatusCode: 500,
			}
		}

		newTodos = append(newTodos, newTodo)
	}

	response := dto.ListWithTodosToResponse(newList, newTodos)
	return &response, nil
}

// Share creates a copy of a list with all its todos for a different user
func (s *TodoListServiceImpl) Share(ctx context.Context, listID, ownerUserID, targetUserID uuid.UUID, req dto.ShareListRequest) (*dto.ListWithTodosResponse, error) {
	// Verify that owner and target are different users
	if ownerUserID == targetUserID {
		return nil, &utils.AppError{
			Err:        utils.ErrBadRequest,
			Message:    "Cannot share list with yourself",
			StatusCode: 400,
		}
	}

	// Verify target user exists
	targetUser, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "Target user not found",
			StatusCode: 404,
		}
	}

	// Fetch the list to be shared
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "List not found",
			StatusCode: 404,
		}
	}

	// Authorization check: ensure list belongs to the owner
	if !list.BelongsToUser(ownerUserID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this list",
			StatusCode: 403,
		}
	}

	// Get todos in this list
	listTodos, err := s.todoRepo.FindByListID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to fetch todos",
			StatusCode: 500,
		}
	}

	// Determine the name for the shared list
	newName := req.CustomName
	if newName == "" {
		newName = fmt.Sprintf("%s (from %s)", list.Name, targetUser.Username)
	}

	// Create new list for the target user
	newList := entity.NewTodoList(targetUserID, newName)

	// Save new list
	if err := s.listRepo.Create(ctx, newList); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to create shared list",
			StatusCode: 500,
		}
	}

	// Duplicate all todos for the target user
	var newTodos []*entity.Todo
	for _, todo := range listTodos {
		// Create new todo with same properties but for target user
		newTodo := entity.NewTodo(
			targetUserID,
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
			return nil, &utils.AppError{
				Err:        err,
				Message:    "Failed to share todos",
				StatusCode: 500,
			}
		}

		newTodos = append(newTodos, newTodo)
	}

	response := dto.ListWithTodosToResponse(newList, newTodos)
	return &response, nil
}
