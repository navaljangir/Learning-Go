package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/constants"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock TodoService
// =============================================================================

type mockTodoService struct {
	createFunc         func(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error)
	getByIDFunc        func(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error)
	listFunc           func(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error)
	updateFunc         func(ctx context.Context, todoID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error)
	toggleCompleteFunc func(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error)
	deleteFunc         func(ctx context.Context, todoID, userID uuid.UUID) error
	moveTodosFunc      func(ctx context.Context, userID uuid.UUID, req dto.MoveTodosRequest) error
}

func (m *mockTodoService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, userID, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoService) GetByID(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, todoID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoService) List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, userID, page, pageSize)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoService) Update(ctx context.Context, todoID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, todoID, userID, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoService) ToggleComplete(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error) {
	if m.toggleCompleteFunc != nil {
		return m.toggleCompleteFunc(ctx, todoID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoService) Delete(ctx context.Context, todoID, userID uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, todoID, userID)
	}
	return errors.New("not implemented")
}

func (m *mockTodoService) MoveTodos(ctx context.Context, userID uuid.UUID, req dto.MoveTodosRequest) error {
	if m.moveTodosFunc != nil {
		return m.moveTodosFunc(ctx, userID, req)
	}
	return errors.New("not implemented")
}

var _ service.TodoService = (*mockTodoService)(nil)

// =============================================================================
// Test Helpers
// =============================================================================

func setupTodoTestRouter(todoService service.TodoService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Error handler middleware
	router.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *utils.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			} else {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.Abort()
		}
	})

	// Middleware to simulate auth
	authMiddleware := func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr != "" {
			userID, _ := uuid.Parse(userIDStr)
			c.Set(constants.ContextUserID, userID)
		}
		c.Next()
	}

	handler := NewTodoHandler(todoService)

	router.POST("/api/v1/todos", authMiddleware, handler.Create)
	router.GET("/api/v1/todos", authMiddleware, handler.List)
	router.GET("/api/v1/todos/:id", authMiddleware, handler.GetByID)
	router.PUT("/api/v1/todos/:id", authMiddleware, handler.Update)
	router.PATCH("/api/v1/todos/:id/toggle", authMiddleware, handler.ToggleComplete)
	router.DELETE("/api/v1/todos/:id", authMiddleware, handler.Delete)

	return router
}

// =============================================================================
// Create Tests
// =============================================================================

func TestTodoHandler_Create(t *testing.T) {
	t.Run("success: create new todo", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()
		mockService := &mockTodoService{
			createFunc: func(ctx context.Context, uid uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
				assert.Equal(t, userID, uid)
				assert.Equal(t, "Buy groceries", req.Title)
				return &dto.TodoResponse{
					ID:        todoID,
					Title:     req.Title,
					Priority:  req.Priority,
					Completed: false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		router := setupTodoTestRouter(mockService)

		reqBody := dto.CreateTodoRequest{
			Title:    "Buy groceries",
			Priority: "medium",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "Buy groceries", data["title"])
	})

	t.Run("fail: service error", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoService{
			createFunc: func(ctx context.Context, uid uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
				return nil, &utils.AppError{
					Err:        errors.New("db error"),
					Message:    "Failed to create todo",
					StatusCode: 500,
				}
			},
		}

		router := setupTodoTestRouter(mockService)

		reqBody := dto.CreateTodoRequest{
			Title:    "Buy groceries",
			Priority: "medium",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// =============================================================================
// List Tests
// =============================================================================

func TestTodoHandler_List(t *testing.T) {
	t.Run("success: list todos with default pagination", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoService{
			listFunc: func(ctx context.Context, uid uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error) {
				assert.Equal(t, userID, uid)
				assert.Equal(t, 1, page)
				assert.Equal(t, 10, pageSize)
				return &dto.TodoListResponse{
					Todos:      []dto.TodoResponse{},
					Total:      0,
					Page:       page,
					PageSize:   pageSize,
					TotalPages: 0,
				}, nil
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/todos", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success: list todos with custom pagination", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoService{
			listFunc: func(ctx context.Context, uid uuid.UUID, page, pageSize int) (*dto.TodoListResponse, error) {
				assert.Equal(t, 2, page)
				assert.Equal(t, 20, pageSize)
				return &dto.TodoListResponse{
					Todos:      []dto.TodoResponse{},
					Total:      0,
					Page:       page,
					PageSize:   pageSize,
					TotalPages: 0,
				}, nil
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/todos?page=2&page_size=20", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// =============================================================================
// GetByID Tests
// =============================================================================

func TestTodoHandler_GetByID(t *testing.T) {
	t.Run("success: get todo by ID", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()
		mockService := &mockTodoService{
			getByIDFunc: func(ctx context.Context, tid, uid uuid.UUID) (*dto.TodoResponse, error) {
				assert.Equal(t, todoID, tid)
				assert.Equal(t, userID, uid)
				return &dto.TodoResponse{
					ID:        todoID,
					Title:     "Test Todo",
					Priority:  "medium",
					Completed: false,
				}, nil
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/todos/"+todoID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail: invalid todo ID", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoService{}
		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/todos/invalid-uuid", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()
		mockService := &mockTodoService{
			getByIDFunc: func(ctx context.Context, tid, uid uuid.UUID) (*dto.TodoResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrNotFound,
					Message:    "Todo not found",
					StatusCode: 404,
				}
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/todos/"+todoID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// =============================================================================
// Update Tests
// =============================================================================

func TestTodoHandler_Update(t *testing.T) {
	t.Run("success: update todo", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()
		newTitle := "Updated Title"

		mockService := &mockTodoService{
			updateFunc: func(ctx context.Context, tid, uid uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
				assert.Equal(t, todoID, tid)
				assert.Equal(t, userID, uid)
				assert.Equal(t, newTitle, *req.Title)
				return &dto.TodoResponse{
					ID:    todoID,
					Title: newTitle,
				}, nil
			},
		}

		router := setupTodoTestRouter(mockService)

		reqBody := dto.UpdateTodoRequest{
			Title: &newTitle,
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/todos/"+todoID.String(), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// =============================================================================
// ToggleComplete Tests
// =============================================================================

func TestTodoHandler_ToggleComplete(t *testing.T) {
	t.Run("success: toggle todo completion", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()

		mockService := &mockTodoService{
			toggleCompleteFunc: func(ctx context.Context, tid, uid uuid.UUID) (*dto.TodoResponse, error) {
				assert.Equal(t, todoID, tid)
				assert.Equal(t, userID, uid)
				now := time.Now()
				return &dto.TodoResponse{
					ID:          todoID,
					Title:       "Test Todo",
					Completed:   true,
					CompletedAt: &now,
				}, nil
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("PATCH", "/api/v1/todos/"+todoID.String()+"/toggle", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// =============================================================================
// Delete Tests
// =============================================================================

func TestTodoHandler_Delete(t *testing.T) {
	t.Run("success: delete todo", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()

		mockService := &mockTodoService{
			deleteFunc: func(ctx context.Context, tid, uid uuid.UUID) error {
				assert.Equal(t, todoID, tid)
				assert.Equal(t, userID, uid)
				return nil
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("DELETE", "/api/v1/todos/"+todoID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail: todo not found", func(t *testing.T) {
		userID := uuid.New()
		todoID := uuid.New()

		mockService := &mockTodoService{
			deleteFunc: func(ctx context.Context, tid, uid uuid.UUID) error {
				return &utils.AppError{
					Err:        utils.ErrNotFound,
					Message:    "Todo not found",
					StatusCode: 404,
				}
			},
		}

		router := setupTodoTestRouter(mockService)

		req := httptest.NewRequest("DELETE", "/api/v1/todos/"+todoID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
