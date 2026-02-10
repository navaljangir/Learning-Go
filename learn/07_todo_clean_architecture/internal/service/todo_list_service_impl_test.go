package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"todo_app/domain/entity"
	"todo_app/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Test TodoListServiceImpl
// =============================================================================

func TestTodoListService_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success: create new list", func(t *testing.T) {
		listRepo := newMockListRepo()
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Create(ctx, userID, dto.CreateListRequest{
			Name: "Work Tasks",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Work Tasks", resp.Name)
		assert.Equal(t, userID, resp.UserID)
		assert.Equal(t, 1, len(listRepo.lists))
	})

	t.Run("fail: repo error on create", func(t *testing.T) {
		listRepo := newMockListRepo()
		listRepo.createErr = errors.New("db error")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Create(ctx, userID, dto.CreateListRequest{
			Name: "Work Tasks",
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to create list")
	})
}

func TestTodoListService_GetByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: get list with todos", func(t *testing.T) {
		listRepo := newMockListRepo()
		todoRepo := newMockTodoRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		seedTodoWithList(todoRepo, userID, "Task 1", false, list.ID)
		seedTodoWithList(todoRepo, userID, "Task 2", true, list.ID)

		svc := NewTodoListService(listRepo, todoRepo, newMockUserRepo(), "secret")

		resp, err := svc.GetByID(ctx, list.ID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Work Tasks", resp.Name)
		assert.Equal(t, 2, len(resp.Todos))
	})

	t.Run("fail: list not found", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.GetByID(ctx, uuid.New(), userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "List not found")
	})

	t.Run("fail: unauthorized access", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.GetByID(ctx, list.ID, otherUserID)

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this list")
	})

	t.Run("fail: repo error on get todos", func(t *testing.T) {
		listRepo := newMockListRepo()
		todoRepo := newMockTodoRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		todoRepo.findByListErr = errors.New("db error")

		svc := NewTodoListService(listRepo, todoRepo, newMockUserRepo(), "secret")

		resp, err := svc.GetByID(ctx, list.ID, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to fetch todos")
	})
}

func TestTodoListService_List(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success: list all user lists", func(t *testing.T) {
		listRepo := newMockListRepo()
		seedList(listRepo, userID, "Work")
		seedList(listRepo, userID, "Personal")
		seedList(listRepo, uuid.New(), "Other User List") // Different user

		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.List(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, resp.Total)
		assert.Equal(t, 2, len(resp.Lists))
	})

	t.Run("success: empty list", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.List(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 0, resp.Total)
		assert.Equal(t, 0, len(resp.Lists))
	})

	t.Run("fail: repo error", func(t *testing.T) {
		listRepo := newMockListRepo()
		listRepo.findByUserErr = errors.New("db error")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.List(ctx, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to fetch lists")
	})
}

func TestTodoListService_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: update list name", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Old Name")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Update(ctx, list.ID, userID, dto.UpdateListRequest{
			Name: "New Name",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "New Name", resp.Name)
	})

	t.Run("fail: list not found", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Update(ctx, uuid.New(), userID, dto.UpdateListRequest{
			Name: "New Name",
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "List not found")
	})

	t.Run("fail: unauthorized access", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Update(ctx, list.ID, otherUserID, dto.UpdateListRequest{
			Name: "New Name",
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this list")
	})

	t.Run("fail: repo error on update", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		listRepo.updateErr = errors.New("db error")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Update(ctx, list.ID, userID, dto.UpdateListRequest{
			Name: "New Name",
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to update list")
	})
}

func TestTodoListService_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: delete list", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		err := svc.Delete(ctx, list.ID, userID)

		assert.NoError(t, err)
	})

	t.Run("fail: list not found", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		err := svc.Delete(ctx, uuid.New(), userID)

		assertAppError(t, err, 404, "List not found")
	})

	t.Run("fail: unauthorized access", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		err := svc.Delete(ctx, list.ID, otherUserID)

		assertAppError(t, err, 403, "Unauthorized access to this list")
	})

	t.Run("fail: repo error on delete", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		listRepo.deleteErr = errors.New("db error")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		err := svc.Delete(ctx, list.ID, userID)

		assertAppError(t, err, 500, "Failed to delete list")
	})
}

func TestTodoListService_Duplicate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: duplicate list with todos", func(t *testing.T) {
		listRepo := newMockListRepo()
		todoRepo := newMockTodoRepo()
		list := seedList(listRepo, userID, "Original List")
		seedTodoWithList(todoRepo, userID, "Task 1", false, list.ID)
		seedTodoWithList(todoRepo, userID, "Task 2", true, list.ID)

		svc := NewTodoListService(listRepo, todoRepo, newMockUserRepo(), "secret")

		resp, err := svc.Duplicate(ctx, list.ID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Original List (Copy)", resp.Name)
		assert.Equal(t, 2, len(resp.Todos))
		assert.Equal(t, 2, len(listRepo.lists)) // original + copy
		assert.Equal(t, 4, len(todoRepo.todos)) // 2 original + 2 copied
	})

	t.Run("success: duplicate empty list", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Empty List")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Duplicate(ctx, list.ID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Empty List (Copy)", resp.Name)
		assert.Equal(t, 0, len(resp.Todos))
	})

	t.Run("fail: list not found", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Duplicate(ctx, uuid.New(), userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "List not found")
	})

	t.Run("fail: unauthorized access", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Duplicate(ctx, list.ID, otherUserID)

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this list")
	})

	t.Run("fail: repo error on create list", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		listRepo.createErr = errors.New("db error")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.Duplicate(ctx, list.ID, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to create duplicate list")
	})

	t.Run("fail: repo error on create todo", func(t *testing.T) {
		listRepo := newMockListRepo()
		todoRepo := newMockTodoRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		seedTodoWithList(todoRepo, userID, "Task 1", false, list.ID)
		todoRepo.createErr = errors.New("db error")

		svc := NewTodoListService(listRepo, todoRepo, newMockUserRepo(), "secret")

		resp, err := svc.Duplicate(ctx, list.ID, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to duplicate todos")
	})
}

func TestTodoListService_GenerateShareLink(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: generate share link", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Shared List")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.GenerateShareLink(ctx, list.ID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.ShareToken)
		assert.Equal(t, 64, len(resp.ShareToken)) // 32 chars UUID + 32 chars HMAC
		assert.True(t, strings.HasPrefix(resp.ShareURL, "/api/v1/lists/import/"))
	})

	t.Run("fail: list not found", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.GenerateShareLink(ctx, uuid.New(), userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "List not found")
	})

	t.Run("fail: unauthorized access", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, userID, "Work Tasks")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.GenerateShareLink(ctx, list.ID, otherUserID)

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this list")
	})
}

func TestTodoListService_ImportSharedList(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	importerID := uuid.New()

	t.Run("success: import shared list", func(t *testing.T) {
		listRepo := newMockListRepo()
		todoRepo := newMockTodoRepo()
		list := seedList(listRepo, ownerID, "Shared List")
		seedTodoWithList(todoRepo, ownerID, "Task 1", false, list.ID)

		svc := NewTodoListService(listRepo, todoRepo, newMockUserRepo(), "secret")

		// Generate share token
		shareResp, _ := svc.GenerateShareLink(ctx, list.ID, ownerID)

		// Import using token
		resp, err := svc.ImportSharedList(ctx, shareResp.ShareToken, importerID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Shared List (shared)", resp.Name)
		assert.Equal(t, importerID, resp.UserID)
		assert.Equal(t, 1, len(resp.Todos))
		assert.Equal(t, 2, len(listRepo.lists)) // original + imported
	})

	t.Run("fail: invalid token format", func(t *testing.T) {
		svc := NewTodoListService(newMockListRepo(), newMockTodoRepo(), newMockUserRepo(), "secret")

		resp, err := svc.ImportSharedList(ctx, "invalid-token", importerID)

		assert.Nil(t, resp)
		assertAppError(t, err, 400, "Invalid or malformed share token")
	})

	t.Run("fail: tampered token", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, ownerID, "Shared List")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		shareResp, _ := svc.GenerateShareLink(ctx, list.ID, ownerID)
		tamperedToken := shareResp.ShareToken[:62] + "xx" // Change last 2 chars

		resp, err := svc.ImportSharedList(ctx, tamperedToken, importerID)

		assert.Nil(t, resp)
		assertAppError(t, err, 400, "Invalid or malformed share token")
	})

	t.Run("fail: list not found", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, ownerID, "Shared List")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		shareResp, _ := svc.GenerateShareLink(ctx, list.ID, ownerID)

		// Delete the list
		delete(listRepo.lists, list.ID)

		resp, err := svc.ImportSharedList(ctx, shareResp.ShareToken, importerID)

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "Shared list not found")
	})

	t.Run("fail: cannot import own list", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, ownerID, "My List")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		shareResp, _ := svc.GenerateShareLink(ctx, list.ID, ownerID)

		resp, err := svc.ImportSharedList(ctx, shareResp.ShareToken, ownerID) // Same user

		assert.Nil(t, resp)
		assertAppError(t, err, 400, "Cannot import your own list, use duplicate instead")
	})

	t.Run("fail: repo error on create list", func(t *testing.T) {
		listRepo := newMockListRepo()
		list := seedList(listRepo, ownerID, "Shared List")
		svc := NewTodoListService(listRepo, newMockTodoRepo(), newMockUserRepo(), "secret")

		shareResp, _ := svc.GenerateShareLink(ctx, list.ID, ownerID)
		listRepo.createErr = errors.New("db error")

		resp, err := svc.ImportSharedList(ctx, shareResp.ShareToken, importerID)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to create imported list")
	})
}

// =============================================================================
// Additional mock helpers for list service tests
// =============================================================================

func seedTodoWithList(repo *mockTodoRepo, userID uuid.UUID, title string, completed bool, listID uuid.UUID) *entity.Todo {
	todo := entity.NewTodo(userID, title, "", entity.PriorityMedium, nil)
	todo.ListID = &listID
	if completed {
		todo.MarkAsCompleted()
	}
	repo.todos[todo.ID] = todo
	return todo
}
