package middleware

import (
	"bytes"
	"encoding/json"
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

// TestErrorHandlerWithValidationError tests handling of ValidationError
func TestErrorHandlerWithValidationError(t *testing.T) {
	router := setupErrorHandlerTest()
	router.GET("/test", func(c *gin.Context) {
		c.Error(&utils.ValidationError{
			Message: "Validation failed",
			Fields: map[string]string{
				"email": "invalid email format",
				"name":  "name is required",
			},
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "Validation failed", resp["error"])
	fields := resp["fields"].(map[string]interface{})
	assert.Equal(t, "invalid email format", fields["email"])
	assert.Equal(t, "name is required", fields["name"])
}

// TestErrorHandlerWithValidatorErrors tests handling of validator.ValidationErrors from ShouldBindJSON
func TestErrorHandlerWithValidatorErrors(t *testing.T) {
	router := setupErrorHandlerTest()

	type testInput struct {
		Email string `json:"email" binding:"required,email"`
	}

	router.POST("/test", func(c *gin.Context) {
		var input testInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.Error(err)
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Send invalid data (missing required email)
	jsonData, _ := json.Marshal(map[string]string{"email": "not-an-email"})
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "Validation failed", resp["error"])
	assert.NotNil(t, resp["fields"])
}

// TestErrorHandlerWithMalformedJSON tests handling of malformed JSON body
func TestErrorHandlerWithMalformedJSON(t *testing.T) {
	router := setupErrorHandlerTest()

	type testInput struct {
		Name string `json:"name" binding:"required"`
	}

	router.POST("/test", func(c *gin.Context) {
		var input testInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.Error(err)
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
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
