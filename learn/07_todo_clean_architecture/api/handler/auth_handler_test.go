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
	"todo_app/api/middleware"
	"todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock UserService
// =============================================================================

type mockUserService struct {
	registerFunc      func(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error)
	loginFunc         func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	getProfileFunc    func(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	updateProfileFunc func(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
}

func (m *mockUserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if m.updateProfileFunc != nil {
		return m.updateProfileFunc(ctx, userID, req)
	}
	return nil, errors.New("not implemented")
}

// Compile-time check
var _ service.UserService = (*mockUserService)(nil)

// =============================================================================
// Test Helpers
// =============================================================================

func setupAuthTestRouter(userService service.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandlerMiddleware())

	handler := NewAuthHandler(userService)
	router.POST("/api/v1/auth/register", handler.Register)
	router.POST("/api/v1/auth/login", handler.Login)

	return router
}

func makeJSONRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	jsonData, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// =============================================================================
// Register Tests
// =============================================================================

func TestAuthHandler_Register(t *testing.T) {
	t.Run("success: register new user", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			registerFunc: func(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
				return &dto.LoginResponse{
					Token: "jwt-token-here",
					User: dto.UserResponse{
						ID:        userID,
						Username:  req.Username,
						Email:     req.Email,
						FullName:  req.FullName,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				}, nil
			},
		}

		router := setupAuthTestRouter(mockService)

		reqBody := dto.RegisterRequest{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "Password123!",
			FullName: "John Doe",
		}

		w := makeJSONRequest(router, "POST", "/api/v1/auth/register", reqBody)

		assert.Equal(t, http.StatusCreated, w.Code)

		var wrapper struct {
			Data dto.LoginResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &wrapper)
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token-here", wrapper.Data.Token)
		assert.Equal(t, "john_doe", wrapper.Data.User.Username)
	})

	t.Run("fail: username already exists", func(t *testing.T) {
		mockService := &mockUserService{
			registerFunc: func(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrDuplicateKey,
					Message:    "Username already exists",
					StatusCode: 400,
				}
			},
		}

		router := setupAuthTestRouter(mockService)

		reqBody := dto.RegisterRequest{
			Username: "existing_user",
			Email:    "new@example.com",
			Password: "Password123!",
			FullName: "New User",
		}

		w := makeJSONRequest(router, "POST", "/api/v1/auth/register", reqBody)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Username already exists", resp["error"])
	})

	t.Run("fail: invalid JSON", func(t *testing.T) {
		mockService := &mockUserService{}
		router := setupAuthTestRouter(mockService)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail: service error", func(t *testing.T) {
		mockService := &mockUserService{
			registerFunc: func(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
				return nil, &utils.AppError{
					Err:        errors.New("db error"),
					Message:    "Failed to create user",
					StatusCode: 500,
				}
			},
		}

		router := setupAuthTestRouter(mockService)

		reqBody := dto.RegisterRequest{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "Password123!",
			FullName: "John Doe",
		}

		w := makeJSONRequest(router, "POST", "/api/v1/auth/register", reqBody)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// =============================================================================
// Login Tests
// =============================================================================

func TestAuthHandler_Login(t *testing.T) {
	t.Run("success: login with valid credentials", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			loginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
				return &dto.LoginResponse{
					Token: "jwt-token-here",
					User: dto.UserResponse{
						ID:        userID,
						Username:  req.Username,
						Email:     "john@example.com",
						FullName:  "John Doe",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				}, nil
			},
		}

		router := setupAuthTestRouter(mockService)

		reqBody := dto.LoginRequest{
			Username: "john_doe",
			Password: "Password123!",
		}

		w := makeJSONRequest(router, "POST", "/api/v1/auth/login", reqBody)

		assert.Equal(t, http.StatusOK, w.Code)

		var wrapper struct {
			Data dto.LoginResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &wrapper)
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token-here", wrapper.Data.Token)
		assert.Equal(t, "john_doe", wrapper.Data.User.Username)
	})

	t.Run("fail: invalid credentials", func(t *testing.T) {
		mockService := &mockUserService{
			loginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrInvalidCredentials,
					Message:    "Invalid credentials",
					StatusCode: 401,
				}
			},
		}

		router := setupAuthTestRouter(mockService)

		reqBody := dto.LoginRequest{
			Username: "john_doe",
			Password: "wrongpassword",
		}

		w := makeJSONRequest(router, "POST", "/api/v1/auth/login", reqBody)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Invalid credentials", resp["error"])
	})

	t.Run("fail: invalid JSON", func(t *testing.T) {
		mockService := &mockUserService{}
		router := setupAuthTestRouter(mockService)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail: missing required fields", func(t *testing.T) {
		mockService := &mockUserService{}
		router := setupAuthTestRouter(mockService)

		reqBody := map[string]string{
			"username": "john_doe",
			// Missing password
		}

		w := makeJSONRequest(router, "POST", "/api/v1/auth/login", reqBody)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
