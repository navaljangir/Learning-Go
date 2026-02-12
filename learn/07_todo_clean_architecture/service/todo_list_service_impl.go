package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	domainService "todo_app/domain/service"
	"todo_app/dto"
	"todo_app/pkg/utils"

	"github.com/google/uuid"
)

// TodoListServiceImpl implements todo list-related business logic
type TodoListServiceImpl struct {
	listRepo    repository.TodoListRepository
	todoRepo    repository.TodoRepository
	userRepo    repository.UserRepository
	shareSecret string
}

// Compile-time check to ensure TodoListServiceImpl implements TodoListService interface
var _ domainService.TodoListService = (*TodoListServiceImpl)(nil)

// NewTodoListService creates a new todo list service
func NewTodoListService(
	listRepo repository.TodoListRepository,
	todoRepo repository.TodoRepository,
	userRepo repository.UserRepository,
	shareSecret string,
) domainService.TodoListService {
	return &TodoListServiceImpl{
		listRepo:    listRepo,
		todoRepo:    todoRepo,
		userRepo:    userRepo,
		shareSecret: shareSecret,
	}
}

// generateShareToken creates an HMAC-signed token from a list UUID
// Token format: {uuid_hex_no_dashes}{hmac_hex_first_32_chars} = 64 chars total
//
// How it works:
//  1. Take the list UUID, remove dashes → 32 hex chars
//  2. Compute HMAC-SHA256(uuid_string, secret) → 64 hex chars
//  3. Take first 32 hex chars of the HMAC (128 bits — plenty for integrity)
//  4. Concatenate: uuid_hex + hmac_hex_prefix = 64 chars
//
func (s *TodoListServiceImpl) generateShareToken(listID uuid.UUID) string {
	// Remove dashes from UUID: "550e8400-e29b-..." → "550e8400e29b..."
	uuidHex := strings.ReplaceAll(listID.String(), "-", "")

	// HMAC-SHA256 the full UUID string (with dashes) using the server secret
	mac := hmac.New(sha256.New, []byte(s.shareSecret))
	mac.Write([]byte(listID.String()))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Token = 32 chars (uuid hex) + 32 chars (hmac prefix) = 64 chars
	return uuidHex + signature[:32]
}

// verifyShareToken extracts and validates the list UUID from a share token
// Returns the list UUID if valid, error if the token is malformed or forged
func (s *TodoListServiceImpl) verifyShareToken(token string) (uuid.UUID, error) {
	if len(token) != 64 {
		return uuid.Nil, fmt.Errorf("invalid share token")
	}

	// Split: first 32 chars = uuid hex, last 32 chars = hmac prefix
	uuidHex := token[:32]
	providedSig := token[32:]

	// Reconstruct UUID with dashes: 8-4-4-4-12
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s",
		uuidHex[0:8], uuidHex[8:12], uuidHex[12:16], uuidHex[16:20], uuidHex[20:32])

	listID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid share token")
	}

	// Recompute HMAC and compare
	mac := hmac.New(sha256.New, []byte(s.shareSecret))
	mac.Write([]byte(listID.String()))
	expectedSig := hex.EncodeToString(mac.Sum(nil))[:32]

	if !hmac.Equal([]byte(providedSig), []byte(expectedSig)) {
		return uuid.Nil, fmt.Errorf("invalid share token")
	}

	return listID, nil
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
func (s *TodoListServiceImpl) Duplicate(ctx context.Context, listID, userID uuid.UUID, req dto.DuplicateListRequest) (*dto.ListWithTodosResponse, error) {
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

		// Preserve completed status if requested
		if req.KeepCompleted && todo.Completed {
			newTodo.MarkAsCompleted()
		}

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

// GenerateShareLink creates an HMAC-signed share token for a list
// The token encodes the list ID + a signature — no database storage needed
func (s *TodoListServiceImpl) GenerateShareLink(ctx context.Context, listID, userID uuid.UUID) (*dto.ShareLinkResponse, error) {
	// Fetch the list
	list, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "List not found",
			StatusCode: 404,
		}
	}

	// Authorization check: only the owner can generate a share link
	if !list.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrForbidden,
			Message:    "Unauthorized access to this list",
			StatusCode: 403,
		}
	}

	token := s.generateShareToken(listID)

	return &dto.ShareLinkResponse{
		ShareURL:   fmt.Sprintf("/api/v1/lists/import/%s", token),
		ShareToken: token,
	}, nil
}

// ImportSharedList verifies the share token, then copies the list + todos into the caller's account
func (s *TodoListServiceImpl) ImportSharedList(ctx context.Context, token string, userID uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error) {
	// Verify the HMAC token and extract list ID
	listID, err := s.verifyShareToken(token)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrBadRequest,
			Message:    "Invalid or malformed share token",
			StatusCode: 400,
		}
	}

	// Fetch the source list
	sourceList, err := s.listRepo.FindByID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        utils.ErrNotFound,
			Message:    "Shared list not found",
			StatusCode: 404,
		}
	}

	// Can't import your own list — use duplicate instead
	if sourceList.BelongsToUser(userID) {
		return nil, &utils.AppError{
			Err:        utils.ErrBadRequest,
			Message:    "Cannot import your own list, use duplicate instead",
			StatusCode: 400,
		}
	}

	// Get all todos from the source list
	sourceTodos, err := s.todoRepo.FindByListID(ctx, listID)
	if err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to fetch todos from shared list",
			StatusCode: 500,
		}
	}

	// Create new list for the caller
	newName := fmt.Sprintf("%s (shared)", sourceList.Name)
	newList := entity.NewTodoList(userID, newName)

	if err := s.listRepo.Create(ctx, newList); err != nil {
		return nil, &utils.AppError{
			Err:        err,
			Message:    "Failed to create imported list",
			StatusCode: 500,
		}
	}

	// Copy all todos to the new list
	var newTodos []*entity.Todo
	for _, todo := range sourceTodos {
		newTodo := entity.NewTodo(
			userID,
			todo.Title,
			todo.Description,
			todo.Priority,
			todo.DueDate,
		)
		newListID := newList.ID
		newTodo.ListID = &newListID

		// Preserve completed status if requested
		if req.KeepCompleted && todo.Completed {
			newTodo.MarkAsCompleted()
		}

		if err := s.todoRepo.Create(ctx, newTodo); err != nil {
			return nil, &utils.AppError{
				Err:        err,
				Message:    "Failed to copy todos",
				StatusCode: 500,
			}
		}

		newTodos = append(newTodos, newTodo)
	}

	response := dto.ListWithTodosToResponse(newList, newTodos)
	return &response, nil
}
