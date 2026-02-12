package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo_app/api/handler"
	"todo_app/api/middleware"
	"todo_app/domain/service"
	"todo_app/dto"
	"todo_app/pkg/constants"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock TodoListService
// =============================================================================

type mockTodoListService struct {
	createFunc            func(ctx context.Context, userID uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error)
	getByIDFunc           func(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error)
	listFunc              func(ctx context.Context, userID uuid.UUID) (*dto.ListsResponse, error)
	updateFunc            func(ctx context.Context, listID, userID uuid.UUID, req dto.UpdateListRequest) (*dto.ListResponse, error)
	deleteFunc            func(ctx context.Context, listID, userID uuid.UUID) error
	duplicateFunc         func(ctx context.Context, listID, userID uuid.UUID, req dto.DuplicateListRequest) (*dto.ListWithTodosResponse, error)
	generateShareLinkFunc func(ctx context.Context, listID, userID uuid.UUID) (*dto.ShareLinkResponse, error)
	importSharedListFunc  func(ctx context.Context, token string, userID uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error)
}

func (m *mockTodoListService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, userID, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoListService) GetByID(ctx context.Context, listID, userID uuid.UUID) (*dto.ListWithTodosResponse, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, listID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoListService) List(ctx context.Context, userID uuid.UUID) (*dto.ListsResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoListService) Update(ctx context.Context, listID, userID uuid.UUID, req dto.UpdateListRequest) (*dto.ListResponse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, listID, userID, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoListService) Delete(ctx context.Context, listID, userID uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, listID, userID)
	}
	return errors.New("not implemented")
}

func (m *mockTodoListService) Duplicate(ctx context.Context, listID, userID uuid.UUID, req dto.DuplicateListRequest) (*dto.ListWithTodosResponse, error) {
	if m.duplicateFunc != nil {
		return m.duplicateFunc(ctx, listID, userID, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoListService) GenerateShareLink(ctx context.Context, listID, userID uuid.UUID) (*dto.ShareLinkResponse, error) {
	if m.generateShareLinkFunc != nil {
		return m.generateShareLinkFunc(ctx, listID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTodoListService) ImportSharedList(ctx context.Context, token string, userID uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error) {
	if m.importSharedListFunc != nil {
		return m.importSharedListFunc(ctx, token, userID, req)
	}
	return nil, errors.New("not implemented")
}

var _ service.TodoListService = (*mockTodoListService)(nil)

// =============================================================================
// Test Helpers
// =============================================================================

func setupTodoListTestRouter(listService service.TodoListService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Middleware to simulate auth
	authMiddleware := func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr != "" {
			userID, _ := uuid.Parse(userIDStr)
			c.Set(constants.ContextUserID, userID)
		}
		c.Next()
	}

	h := handler.NewTodoListHandler(listService)

	router.POST("/api/v1/lists", authMiddleware, h.Create)
	router.GET("/api/v1/lists", authMiddleware, h.List)
	router.GET("/api/v1/lists/:id", authMiddleware, h.GetByID)
	router.PUT("/api/v1/lists/:id", authMiddleware, h.Update)
	router.DELETE("/api/v1/lists/:id", authMiddleware, h.Delete)
	router.POST("/api/v1/lists/:id/duplicate", authMiddleware, h.Duplicate)
	router.POST("/api/v1/lists/:id/share", authMiddleware, h.GenerateShareLink)
	router.POST("/api/v1/lists/import/:token", authMiddleware, h.ImportSharedList)

	return router
}

// =============================================================================
// Create Tests
// =============================================================================

func TestTodoListHandler_Create(t *testing.T) {
	t.Run("success: create new list", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			createFunc: func(ctx context.Context, uid uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error) {
				assert.Equal(t, userID, uid)
				assert.Equal(t, "Work Tasks", req.Name)
				return &dto.ListResponse{
					ID:        listID,
					UserID:    userID,
					Name:      req.Name,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		reqBody := dto.CreateListRequest{
			Name: "Work Tasks",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/lists", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "Work Tasks", data["name"])
	})

	t.Run("fail: service error", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoListService{
			createFunc: func(ctx context.Context, uid uuid.UUID, req dto.CreateListRequest) (*dto.ListResponse, error) {
				return nil, &utils.AppError{
					Err:        errors.New("db error"),
					Message:    "Failed to create list",
					StatusCode: 500,
				}
			},
		}

		router := setupTodoListTestRouter(mockService)

		reqBody := dto.CreateListRequest{
			Name: "Work Tasks",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/lists", bytes.NewBuffer(jsonData))
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

func TestTodoListHandler_List(t *testing.T) {
	t.Run("success: list all user lists", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoListService{
			listFunc: func(ctx context.Context, uid uuid.UUID) (*dto.ListsResponse, error) {
				assert.Equal(t, userID, uid)
				return &dto.ListsResponse{
					Lists: []dto.ListResponse{
						{
							ID:     uuid.New(),
							UserID: userID,
							Name:   "Work",
						},
						{
							ID:     uuid.New(),
							UserID: userID,
							Name:   "Personal",
						},
					},
					Total: 2,
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/lists", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, float64(2), data["total"])
	})
}

// =============================================================================
// GetByID Tests
// =============================================================================

func TestTodoListHandler_GetByID(t *testing.T) {
	t.Run("success: get list by ID", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			getByIDFunc: func(ctx context.Context, lid, uid uuid.UUID) (*dto.ListWithTodosResponse, error) {
				assert.Equal(t, listID, lid)
				assert.Equal(t, userID, uid)
				return &dto.ListWithTodosResponse{
					ID:     listID,
					UserID: userID,
					Name:   "Work Tasks",
					Todos:  []dto.TodoResponse{},
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/lists/"+listID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail: invalid list ID", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoListService{}
		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/lists/invalid-uuid", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail: list not found", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			getByIDFunc: func(ctx context.Context, lid, uid uuid.UUID) (*dto.ListWithTodosResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrNotFound,
					Message:    "List not found",
					StatusCode: 404,
				}
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/lists/"+listID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// =============================================================================
// Update Tests
// =============================================================================

func TestTodoListHandler_Update(t *testing.T) {
	t.Run("success: update list name", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			updateFunc: func(ctx context.Context, lid, uid uuid.UUID, req dto.UpdateListRequest) (*dto.ListResponse, error) {
				assert.Equal(t, listID, lid)
				assert.Equal(t, userID, uid)
				assert.Equal(t, "New Name", req.Name)
				return &dto.ListResponse{
					ID:     listID,
					UserID: userID,
					Name:   req.Name,
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		reqBody := dto.UpdateListRequest{
			Name: "New Name",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/lists/"+listID.String(), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// =============================================================================
// Delete Tests
// =============================================================================

func TestTodoListHandler_Delete(t *testing.T) {
	t.Run("success: delete list", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			deleteFunc: func(ctx context.Context, lid, uid uuid.UUID) error {
				assert.Equal(t, listID, lid)
				assert.Equal(t, userID, uid)
				return nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("DELETE", "/api/v1/lists/"+listID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("fail: list not found", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			deleteFunc: func(ctx context.Context, lid, uid uuid.UUID) error {
				return &utils.AppError{
					Err:        utils.ErrNotFound,
					Message:    "List not found",
					StatusCode: 404,
				}
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("DELETE", "/api/v1/lists/"+listID.String(), nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// =============================================================================
// Duplicate Tests
// =============================================================================

func TestTodoListHandler_Duplicate(t *testing.T) {
	t.Run("success: duplicate list without body defaults keep_completed=false", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			duplicateFunc: func(ctx context.Context, lid, uid uuid.UUID, req dto.DuplicateListRequest) (*dto.ListWithTodosResponse, error) {
				assert.Equal(t, listID, lid)
				assert.Equal(t, userID, uid)
				assert.False(t, req.KeepCompleted)
				return &dto.ListWithTodosResponse{
					ID:     uuid.New(),
					UserID: userID,
					Name:   "Work Tasks (Copy)",
					Todos:  []dto.TodoResponse{},
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("POST", "/api/v1/lists/"+listID.String()+"/duplicate", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("success: duplicate list with keep_completed=true", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			duplicateFunc: func(ctx context.Context, lid, uid uuid.UUID, req dto.DuplicateListRequest) (*dto.ListWithTodosResponse, error) {
				assert.Equal(t, listID, lid)
				assert.Equal(t, userID, uid)
				assert.True(t, req.KeepCompleted)
				return &dto.ListWithTodosResponse{
					ID:     uuid.New(),
					UserID: userID,
					Name:   "Work Tasks (Copy)",
					Todos:  []dto.TodoResponse{},
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		reqBody := dto.DuplicateListRequest{KeepCompleted: true}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/lists/"+listID.String()+"/duplicate", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

// =============================================================================
// GenerateShareLink Tests
// =============================================================================

func TestTodoListHandler_GenerateShareLink(t *testing.T) {
	t.Run("success: generate share link", func(t *testing.T) {
		userID := uuid.New()
		listID := uuid.New()
		mockService := &mockTodoListService{
			generateShareLinkFunc: func(ctx context.Context, lid, uid uuid.UUID) (*dto.ShareLinkResponse, error) {
				assert.Equal(t, listID, lid)
				assert.Equal(t, userID, uid)
				return &dto.ShareLinkResponse{
					ShareURL:   "/api/v1/lists/import/abc123",
					ShareToken: "abc123",
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("POST", "/api/v1/lists/"+listID.String()+"/share", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.NotEmpty(t, data["share_token"])
	})
}

// =============================================================================
// ImportSharedList Tests
// =============================================================================

func TestTodoListHandler_ImportSharedList(t *testing.T) {
	t.Run("success: import shared list without body defaults keep_completed=false", func(t *testing.T) {
		userID := uuid.New()
		token := "abc123"
		mockService := &mockTodoListService{
			importSharedListFunc: func(ctx context.Context, tok string, uid uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error) {
				assert.Equal(t, token, tok)
				assert.Equal(t, userID, uid)
				assert.False(t, req.KeepCompleted)
				return &dto.ListWithTodosResponse{
					ID:     uuid.New(),
					UserID: userID,
					Name:   "Imported List (shared)",
					Todos:  []dto.TodoResponse{},
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("POST", "/api/v1/lists/import/"+token, nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("success: import shared list with keep_completed=true", func(t *testing.T) {
		userID := uuid.New()
		token := "abc123"
		mockService := &mockTodoListService{
			importSharedListFunc: func(ctx context.Context, tok string, uid uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error) {
				assert.Equal(t, token, tok)
				assert.Equal(t, userID, uid)
				assert.True(t, req.KeepCompleted)
				return &dto.ListWithTodosResponse{
					ID:     uuid.New(),
					UserID: userID,
					Name:   "Imported List (shared)",
					Todos:  []dto.TodoResponse{},
				}, nil
			},
		}

		router := setupTodoListTestRouter(mockService)

		reqBody := dto.ImportListRequest{KeepCompleted: true}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/lists/import/"+token, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("fail: invalid share token", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockTodoListService{
			importSharedListFunc: func(ctx context.Context, tok string, uid uuid.UUID, req dto.ImportListRequest) (*dto.ListWithTodosResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrBadRequest,
					Message:    "Invalid or malformed share token",
					StatusCode: 400,
				}
			},
		}

		router := setupTodoListTestRouter(mockService)

		req := httptest.NewRequest("POST", "/api/v1/lists/import/invalid", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
