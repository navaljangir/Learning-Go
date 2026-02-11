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
	"todo_app/internal/dto"
	"todo_app/pkg/constants"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Test Helpers
// =============================================================================

func setupUserTestRouter(userService *mockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandlerMiddleware())

	handler := NewUserHandler(userService)

	// Protected routes - simulate auth middleware by setting context
	router.GET("/api/v1/users/profile", func(c *gin.Context) {
		// Simulate extracting userID from JWT (normally done by auth middleware)
		userIDStr := c.GetHeader("X-User-ID") // Test helper header
		if userIDStr != "" {
			userID, _ := uuid.Parse(userIDStr)
			c.Set(constants.ContextUserID, userID)
		}
		c.Next()
	}, handler.GetProfile)

	router.PUT("/api/v1/users/profile", func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr != "" {
			userID, _ := uuid.Parse(userIDStr)
			c.Set(constants.ContextUserID, userID)
		}
		c.Next()
	}, handler.UpdateProfile)

	return router
}

// =============================================================================
// GetProfile Tests
// =============================================================================

func TestUserHandler_GetProfile(t *testing.T) {
	t.Run("success: get user profile", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			getProfileFunc: func(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
				assert.Equal(t, userID, id)
				return &dto.UserResponse{
					ID:        userID,
					Username:  "john_doe",
					Email:     "john@example.com",
					FullName:  "John Doe",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		router := setupUserTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "john_doe", data["username"])
		assert.Equal(t, "john@example.com", data["email"])
		assert.Equal(t, "John Doe", data["full_name"])
	})

	t.Run("fail: user not found", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			getProfileFunc: func(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrNotFound,
					Message:    "User not found",
					StatusCode: 404,
				}
			},
		}

		router := setupUserTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "User not found", resp["error"])
	})

	t.Run("fail: service error", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			getProfileFunc: func(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
				return nil, &utils.AppError{
					Err:        errors.New("db error"),
					Message:    "Internal server error",
					StatusCode: 500,
				}
			},
		}

		router := setupUserTestRouter(mockService)

		req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// =============================================================================
// UpdateProfile Tests
// =============================================================================

func TestUserHandler_UpdateProfile(t *testing.T) {
	t.Run("success: update user profile", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			updateProfileFunc: func(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
				assert.Equal(t, userID, id)
				assert.Equal(t, "John Smith", req.FullName)
				return &dto.UserResponse{
					ID:        userID,
					Username:  "john_doe",
					Email:     "john@example.com",
					FullName:  req.FullName,
					CreatedAt: time.Now().Add(-24 * time.Hour),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		router := setupUserTestRouter(mockService)

		reqBody := dto.UpdateUserRequest{
			FullName: "John Smith",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "John Smith", data["full_name"])
	})

	t.Run("fail: invalid JSON", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{}
		router := setupUserTestRouter(mockService)

		req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail: user not found", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			updateProfileFunc: func(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
				return nil, &utils.AppError{
					Err:        utils.ErrNotFound,
					Message:    "User not found",
					StatusCode: 404,
				}
			},
		}

		router := setupUserTestRouter(mockService)

		reqBody := dto.UpdateUserRequest{
			FullName: "John Smith",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("fail: service error", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mockUserService{
			updateProfileFunc: func(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
				return nil, &utils.AppError{
					Err:        errors.New("db error"),
					Message:    "Failed to update user",
					StatusCode: 500,
				}
			},
		}

		router := setupUserTestRouter(mockService)

		reqBody := dto.UpdateUserRequest{
			FullName: "John Smith",
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
