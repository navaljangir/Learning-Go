package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock TodoRepository
// =============================================================================
//
// Stores todos in a map (in-memory). Has injectable error fields to simulate
// repo failures.
//
// Example:
//
//	repo := newMockTodoRepo()
//	repo.createErr = errors.New("db down")  // next Create() call will fail
//
// =============================================================================

type mockTodoRepo struct {
	todos map[uuid.UUID]*entity.Todo

	createErr      error
	findByIDErr    error
	findByUserErr  error
	updateErr      error
	deleteErr      error
	countByUserErr error
	updateListErr  error

	countOverride *int64
}

func newMockTodoRepo() *mockTodoRepo {
	return &mockTodoRepo{
		todos: make(map[uuid.UUID]*entity.Todo),
	}
}

func (m *mockTodoRepo) Create(_ context.Context, todo *entity.Todo) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *mockTodoRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Todo, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	todo, ok := m.todos[id]
	if !ok {
		return nil, utils.ErrNotFound
	}
	return todo, nil
}

func (m *mockTodoRepo) FindByUserID(_ context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Todo, error) {
	if m.findByUserErr != nil {
		return nil, m.findByUserErr
	}
	var result []*entity.Todo
	for _, todo := range m.todos {
		if todo.UserID == userID {
			result = append(result, todo)
		}
	}
	if offset >= len(result) {
		return []*entity.Todo{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *mockTodoRepo) FindByListID(_ context.Context, _ uuid.UUID) ([]*entity.Todo, error) {
	return nil, nil
}

func (m *mockTodoRepo) FindWithFilters(_ context.Context, _ repository.TodoFilter, _, _ int) ([]*entity.Todo, error) {
	return nil, nil
}

func (m *mockTodoRepo) Update(_ context.Context, todo *entity.Todo) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *mockTodoRepo) Delete(_ context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.todos, id)
	return nil
}

func (m *mockTodoRepo) Count(_ context.Context, _ repository.TodoFilter) (int64, error) {
	return 0, nil
}

func (m *mockTodoRepo) CountByUser(_ context.Context, userID uuid.UUID) (int64, error) {
	if m.countByUserErr != nil {
		return 0, m.countByUserErr
	}
	if m.countOverride != nil {
		return *m.countOverride, nil
	}
	var count int64
	for _, todo := range m.todos {
		if todo.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockTodoRepo) UpdateListID(_ context.Context, _ []uuid.UUID, _ *uuid.UUID, _ uuid.UUID) error {
	if m.updateListErr != nil {
		return m.updateListErr
	}
	return nil
}

// =============================================================================
// Mock TodoListRepository
// =============================================================================
//
// The todo service now needs a list repo to verify list ownership when
// creating a todo with a list_id. This mock stores lists in a map.
//
// =============================================================================

type mockListRepo struct {
	lists map[uuid.UUID]*entity.TodoList

	findByIDErr error
}

func newMockListRepo() *mockListRepo {
	return &mockListRepo{
		lists: make(map[uuid.UUID]*entity.TodoList),
	}
}

func (m *mockListRepo) Create(_ context.Context, list *entity.TodoList) error {
	m.lists[list.ID] = list
	return nil
}

func (m *mockListRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.TodoList, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	list, ok := m.lists[id]
	if !ok {
		return nil, utils.ErrNotFound
	}
	return list, nil
}

func (m *mockListRepo) FindByUserID(_ context.Context, _ uuid.UUID) ([]*entity.TodoList, error) {
	return nil, nil
}

func (m *mockListRepo) Update(_ context.Context, list *entity.TodoList) error {
	m.lists[list.ID] = list
	return nil
}

func (m *mockListRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.lists, id)
	return nil
}

func (m *mockListRepo) CountByUser(_ context.Context, _ uuid.UUID) (int64, error) {
	return int64(len(m.lists)), nil
}

func (m *mockListRepo) Duplicate(_ context.Context, _ uuid.UUID, _ string) (*entity.TodoList, error) {
	return nil, nil
}

// =============================================================================
// Test helpers
// =============================================================================

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

// seedTodo creates a todo in the mock repo and returns it
func seedTodo(repo *mockTodoRepo, userID uuid.UUID, title string, completed bool) *entity.Todo {
	todo := entity.NewTodo(userID, title, "", entity.PriorityMedium, nil)
	if completed {
		todo.MarkAsCompleted()
	}
	repo.todos[todo.ID] = todo
	return todo
}

// seedList creates a list in the mock list repo and returns it
func seedList(repo *mockListRepo, userID uuid.UUID, name string) *entity.TodoList {
	list := entity.NewTodoList(userID, name)
	repo.lists[list.ID] = list
	return list
}

// assertAppError checks that err is an *AppError with the expected status and message
func assertAppError(t *testing.T, err error, wantStatus int, wantMsg string) {
	t.Helper()
	assert.Error(t, err)
	var appErr *utils.AppError
	assert.True(t, errors.As(err, &appErr), "error should be *utils.AppError, got %T", err)
	assert.Equal(t, wantStatus, appErr.StatusCode)
	assert.Equal(t, wantMsg, appErr.Message)
}

// =============================================================================
// Create Tests
// =============================================================================

func TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success: basic todo with required fields", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())

		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:    "Buy groceries",
			Priority: "medium",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Buy groceries", resp.Title)
		assert.Equal(t, "medium", resp.Priority)
		assert.False(t, resp.Completed)
		assert.Nil(t, resp.CompletedAt)
		assert.Equal(t, 1, len(todoRepo.todos), "todo should be saved in repo")
	})

	t.Run("success: with description and due date", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		dueDate := time.Now().Add(48 * time.Hour)
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:       "Write docs",
			Description: "API documentation for v2",
			Priority:    "high",
			DueDate:     &dueDate,
		})

		assert.NoError(t, err)
		assert.Equal(t, "Write docs", resp.Title)
		assert.Equal(t, "API documentation for v2", resp.Description)
		assert.Equal(t, "high", resp.Priority)
		assert.NotNil(t, resp.DueDate)
	})

	t.Run("success: completed without completed_at auto-sets current time", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		before := time.Now()
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:     "Done task",
			Priority:  "low",
			Completed: true,
		})
		after := time.Now()

		assert.NoError(t, err)
		assert.True(t, resp.Completed)
		assert.NotNil(t, resp.CompletedAt)
		assert.False(t, resp.CompletedAt.Before(before))
		assert.False(t, resp.CompletedAt.After(after))
	})

	t.Run("success: completed with past completed_at", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		pastTime := time.Now().Add(-24 * time.Hour)
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:       "Imported task",
			Priority:    "medium",
			Completed:   true,
			CompletedAt: &pastTime,
		})

		assert.NoError(t, err)
		assert.True(t, resp.Completed)
		assert.True(t, resp.CompletedAt.Equal(pastTime))
	})

	t.Run("fail: completed_at in the future", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())

		futureTime := time.Now().Add(24 * time.Hour)
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:       "Future task",
			Priority:    "medium",
			Completed:   true,
			CompletedAt: &futureTime,
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 400, "completed_at cannot be in the future")
		assert.Equal(t, 0, len(todoRepo.todos), "todo should NOT be saved")
	})

	t.Run("fail: invalid list_id format", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		badListID := "not-a-uuid"
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:    "Task",
			Priority: "low",
			ListID:   &badListID,
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 400, "Invalid list ID format")
	})

	t.Run("success: valid list_id that belongs to user", func(t *testing.T) {
		listRepo := newMockListRepo()
		svc := NewTodoService(newMockTodoRepo(), listRepo)

		// Create a list owned by this user
		list := seedList(listRepo, userID, "Work Tasks")

		listID := list.ID.String()
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:    "Listed task",
			Priority: "high",
			ListID:   &listID,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp.ListID, "todo should be assigned to the list")
		assert.Equal(t, list.ID, *resp.ListID)
	})

	t.Run("list_id belongs to different user: creates as global todo", func(t *testing.T) {
		listRepo := newMockListRepo()
		svc := NewTodoService(newMockTodoRepo(), listRepo)

		// Create a list owned by ANOTHER user
		otherUserID := uuid.New()
		otherList := seedList(listRepo, otherUserID, "Other's list")

		listID := otherList.ID.String()
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:    "My task",
			Priority: "medium",
			ListID:   &listID,
		})

		assert.NoError(t, err, "should not error — todo is created as global")
		assert.NotNil(t, resp)
		assert.Nil(t, resp.ListID, "list_id should be nil (global) since list belongs to other user")
		assert.Equal(t, "My task", resp.Title)
	})

	t.Run("list_id does not exist: creates as global todo", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		// Valid UUID but no list with this ID exists
		nonExistentListID := uuid.New().String()
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:    "My task",
			Priority: "low",
			ListID:   &nonExistentListID,
		})

		assert.NoError(t, err, "should not error — todo is created as global")
		assert.NotNil(t, resp)
		assert.Nil(t, resp.ListID, "list_id should be nil (global) since list doesn't exist")
	})

	t.Run("fail: repo Create returns error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.createErr = errors.New("database connection lost")
		svc := NewTodoService(todoRepo, newMockListRepo())

		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:    "Task",
			Priority: "medium",
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to create todo")
	})

	t.Run("not completed ignores completed_at", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		pastTime := time.Now().Add(-1 * time.Hour)
		resp, err := svc.Create(ctx, userID, dto.CreateTodoRequest{
			Title:       "Incomplete",
			Priority:    "medium",
			Completed:   false,
			CompletedAt: &pastTime,
		})

		assert.NoError(t, err)
		assert.False(t, resp.Completed)
		assert.Nil(t, resp.CompletedAt, "completed_at should be nil when completed=false")
	})
}

// =============================================================================
// GetByID Tests
// =============================================================================

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: returns owned todo", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "My task", false)

		resp, err := svc.GetByID(ctx, todo.ID, userID)

		assert.NoError(t, err)
		assert.Equal(t, todo.ID, resp.ID)
		assert.Equal(t, "My task", resp.Title)
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.GetByID(ctx, uuid.New(), userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "Todo not found")
	})

	t.Run("fail: todo belongs to different user", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, otherUserID, "Other user's task", false)

		resp, err := svc.GetByID(ctx, todo.ID, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this todo")
	})
}

// =============================================================================
// List Tests
// =============================================================================

func TestList(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success: returns paginated todos", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		seedTodo(todoRepo, userID, "Task 1", false)
		seedTodo(todoRepo, userID, "Task 2", true)

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, 2, len(resp.Todos))
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 10, resp.PageSize)
		assert.Equal(t, 1, resp.TotalPages)
	})

	t.Run("page < 1 defaults to 1", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		seedTodo(todoRepo, userID, "Task", false)

		resp, err := svc.List(ctx, userID, 0, 10)

		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Page)
	})

	t.Run("negative page defaults to 1", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.List(ctx, userID, -5, 10)

		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Page)
	})

	t.Run("pageSize < 1 defaults to 10", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.List(ctx, userID, 1, 0)

		assert.NoError(t, err)
		assert.Equal(t, 10, resp.PageSize)
	})

	t.Run("pageSize > 100 defaults to 10", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.List(ctx, userID, 1, 200)

		assert.NoError(t, err)
		assert.Equal(t, 10, resp.PageSize)
	})

	t.Run("totalPages calculation", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())

		count := int64(25)
		todoRepo.countOverride = &count

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(25), resp.Total)
		assert.Equal(t, 3, resp.TotalPages) // 25/10 = 2.5, rounds up to 3
	})

	t.Run("totalPages exact division", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())

		count := int64(20)
		todoRepo.countOverride = &count

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, resp.TotalPages) // 20/10 = exactly 2
	})

	t.Run("fail: FindByUserID repo error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.findByUserErr = errors.New("query timeout")
		svc := NewTodoService(todoRepo, newMockListRepo())

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to fetch todos")
	})

	t.Run("fail: CountByUser repo error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.countByUserErr = errors.New("count failed")
		svc := NewTodoService(todoRepo, newMockListRepo())

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to count todos")
	})

	t.Run("empty result for user with no todos", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), resp.Total)
		assert.Equal(t, 0, len(resp.Todos))
		assert.Equal(t, 0, resp.TotalPages)
	})

	t.Run("does not return other user's todos", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		seedTodo(todoRepo, uuid.New(), "Other user's task", false)
		seedTodo(todoRepo, userID, "My task", false)

		resp, err := svc.List(ctx, userID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
		assert.Equal(t, 1, len(resp.Todos))
		assert.Equal(t, "My task", resp.Todos[0].Title)
	})
}

// =============================================================================
// Update Tests
// =============================================================================

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: update title only", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Old title", false)

		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Title: strPtr("New title"),
		})

		assert.NoError(t, err)
		assert.Equal(t, "New title", resp.Title)
		assert.Equal(t, "", resp.Description, "unchanged fields stay the same")
	})

	t.Run("success: update multiple fields", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		dueDate := time.Now().Add(72 * time.Hour)
		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Title:       strPtr("Updated task"),
			Description: strPtr("With details now"),
			Priority:    strPtr("urgent"),
			DueDate:     &dueDate,
		})

		assert.NoError(t, err)
		assert.Equal(t, "Updated task", resp.Title)
		assert.Equal(t, "With details now", resp.Description)
		assert.Equal(t, "urgent", resp.Priority)
		assert.NotNil(t, resp.DueDate)
	})

	t.Run("success: mark as completed auto-sets completed_at", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Pending task", false)

		before := time.Now()
		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Completed: boolPtr(true),
		})
		after := time.Now()

		assert.NoError(t, err)
		assert.True(t, resp.Completed)
		assert.NotNil(t, resp.CompletedAt)
		assert.False(t, resp.CompletedAt.Before(before))
		assert.False(t, resp.CompletedAt.After(after))
	})

	t.Run("success: mark as completed with past completed_at", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		pastTime := time.Now().Add(-48 * time.Hour)
		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Completed:   boolPtr(true),
			CompletedAt: &pastTime,
		})

		assert.NoError(t, err)
		assert.True(t, resp.Completed)
		assert.True(t, resp.CompletedAt.Equal(pastTime))
	})

	t.Run("success: mark as incomplete clears completed_at", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Done task", true)

		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Completed: boolPtr(false),
		})

		assert.NoError(t, err)
		assert.False(t, resp.Completed)
		assert.Nil(t, resp.CompletedAt)
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.Update(ctx, uuid.New(), userID, dto.UpdateTodoRequest{
			Title: strPtr("Doesn't matter"),
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "Todo not found")
	})

	t.Run("fail: todo belongs to different user", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, otherUserID, "Not mine", false)

		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Title: strPtr("Trying to update"),
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this todo")
	})

	t.Run("fail: completed_at in the future", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		futureTime := time.Now().Add(24 * time.Hour)
		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Completed:   boolPtr(true),
			CompletedAt: &futureTime,
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 400, "completed_at cannot be in the future")
	})

	t.Run("fail: repo Update returns error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.updateErr = errors.New("disk full")
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		resp, err := svc.Update(ctx, todo.ID, userID, dto.UpdateTodoRequest{
			Title: strPtr("Updated"),
		})

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to update todo")
	})
}

// =============================================================================
// ToggleComplete Tests
// =============================================================================

func TestToggleComplete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: toggle incomplete to completed", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Pending", false)

		resp, err := svc.ToggleComplete(ctx, todo.ID, userID)

		assert.NoError(t, err)
		assert.True(t, resp.Completed)
		assert.NotNil(t, resp.CompletedAt)
	})

	t.Run("success: toggle completed to incomplete", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Done", true)

		resp, err := svc.ToggleComplete(ctx, todo.ID, userID)

		assert.NoError(t, err)
		assert.False(t, resp.Completed)
		assert.Nil(t, resp.CompletedAt)
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		resp, err := svc.ToggleComplete(ctx, uuid.New(), userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 404, "Todo not found")
	})

	t.Run("fail: todo belongs to different user", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, otherUserID, "Not mine", false)

		resp, err := svc.ToggleComplete(ctx, todo.ID, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 403, "Unauthorized access to this todo")
	})

	t.Run("fail: repo Update returns error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.updateErr = errors.New("connection reset")
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		resp, err := svc.ToggleComplete(ctx, todo.ID, userID)

		assert.Nil(t, resp)
		assertAppError(t, err, 500, "Failed to update todo")
	})
}

// =============================================================================
// Delete Tests
// =============================================================================

func TestDelete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: deletes owned todo", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "To delete", false)

		err := svc.Delete(ctx, todo.ID, userID)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(todoRepo.todos), "todo should be removed from repo")
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		err := svc.Delete(ctx, uuid.New(), userID)

		assertAppError(t, err, 404, "Todo not found")
	})

	t.Run("fail: todo belongs to different user", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, otherUserID, "Not mine", false)

		err := svc.Delete(ctx, todo.ID, userID)

		assertAppError(t, err, 403, "Unauthorized access to this todo")
		assert.Equal(t, 1, len(todoRepo.todos), "todo should NOT be deleted")
	})

	t.Run("fail: repo Delete returns error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.deleteErr = errors.New("foreign key constraint")
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		err := svc.Delete(ctx, todo.ID, userID)

		assertAppError(t, err, 500, "Failed to delete todo")
	})
}

// =============================================================================
// MoveTodos Tests
// =============================================================================

func TestMoveTodos(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success: move todos to a list", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo1 := seedTodo(todoRepo, userID, "Task 1", false)
		todo2 := seedTodo(todoRepo, userID, "Task 2", false)

		listID := uuid.New().String()
		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{todo1.ID.String(), todo2.ID.String()},
			ListID:  &listID,
		})

		assert.NoError(t, err)
	})

	t.Run("success: move to global (nil list_id)", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Listed task", false)

		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{todo.ID.String()},
			ListID:  nil,
		})

		assert.NoError(t, err)
	})

	t.Run("fail: empty todo_ids", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{},
		})

		assertAppError(t, err, 400, "No todos specified")
	})

	t.Run("fail: invalid todo_id format", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{"not-a-valid-uuid"},
		})

		assertAppError(t, err, 400, "Invalid todo ID format")
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		svc := NewTodoService(newMockTodoRepo(), newMockListRepo())

		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{uuid.New().String()},
		})

		assertAppError(t, err, 404, "One or more todos not found")
	})

	t.Run("fail: todo belongs to different user", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, otherUserID, "Not mine", false)

		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{todo.ID.String()},
		})

		assertAppError(t, err, 403, "Unauthorized access to one or more todos")
	})

	t.Run("fail: invalid list_id format", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		badListID := "not-a-uuid"
		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{todo.ID.String()},
			ListID:  &badListID,
		})

		assertAppError(t, err, 400, "Invalid list ID format")
	})

	t.Run("fail: repo UpdateListID returns error", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		todoRepo.updateListErr = errors.New("bulk update failed")
		svc := NewTodoService(todoRepo, newMockListRepo())
		todo := seedTodo(todoRepo, userID, "Task", false)

		listID := uuid.New().String()
		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{todo.ID.String()},
			ListID:  &listID,
		})

		assertAppError(t, err, 500, "Failed to move todos")
	})

	t.Run("fail: mix of owned and unowned todos", func(t *testing.T) {
		todoRepo := newMockTodoRepo()
		svc := NewTodoService(todoRepo, newMockListRepo())
		myTodo := seedTodo(todoRepo, userID, "Mine", false)
		otherTodo := seedTodo(todoRepo, otherUserID, "Not mine", false)

		err := svc.MoveTodos(ctx, userID, dto.MoveTodosRequest{
			TodoIDs: []string{myTodo.ID.String(), otherTodo.ID.String()},
		})

		assertAppError(t, err, 403, "Unauthorized access to one or more todos")
	})
}
