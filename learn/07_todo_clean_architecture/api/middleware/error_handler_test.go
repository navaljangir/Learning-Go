package middleware

import (
	"errors"
	"net/http/httptest"
	"testing"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupErrorHandlerTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlerMiddleware())
	return router
}

// TestErrorHandlerWithAppError tests handling of AppError
func TestErrorHandlerWithAppError(t *testing.T) {
	router := setupErrorHandlerTest()
	router.GET("/test", func(c *gin.Context) {
		c.Error(&utils.AppError{
			Err:        errors.New("something went wrong"),
			Message:    "Resource not available",
			StatusCode: 503,
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 503, w.Code)
	assert.Contains(t, w.Body.String(), "Resource not available")
}

// TestErrorHandlerWithSentinelErrors tests sentinel error handling
func TestErrorHandlerWithSentinelErrors(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"ErrNotFound", utils.ErrNotFound, 404},
		{"ErrForbidden", utils.ErrForbidden, 403},
		{"ErrBadRequest", utils.ErrBadRequest, 400},
		{"ErrInvalidCredentials", utils.ErrInvalidCredentials, 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupErrorHandlerTest()
			router.GET("/test", func(c *gin.Context) {
				c.Error(tt.err)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestErrorHandlerWithGenericError tests unknown error defaults to 500
func TestErrorHandlerWithGenericError(t *testing.T) {
	router := setupErrorHandlerTest()
	router.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("unexpected error"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
}

// TestErrorHandlerWithNoErrors tests middleware passes through without errors
func TestErrorHandlerWithNoErrors(t *testing.T) {
	router := setupErrorHandlerTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
