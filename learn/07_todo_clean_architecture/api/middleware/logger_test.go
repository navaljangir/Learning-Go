package middleware

import (
	"bytes"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"todo_app/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// captureLogOutput captures log output for testing
func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

// TestRequestIDMiddleware tests request ID generation and injection
func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	
	var capturedRequestID string
	router.GET("/test", func(c *gin.Context) {
		capturedRequestID = c.GetString(constants.ContextRequestID)
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.NotEmpty(t, capturedRequestID, "request ID should be set in context")
	assert.Equal(t, capturedRequestID, w.Header().Get("X-Request-ID"))
}

// TestRequestIDMiddlewareWithExistingID tests using provided request ID
func TestRequestIDMiddlewareWithExistingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	
	providedRequestID := "test-request-id-123"
	var capturedRequestID string
	
	router.GET("/test", func(c *gin.Context) {
		capturedRequestID = c.GetString(constants.ContextRequestID)
		c.JSON(200, gin.H{})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", providedRequestID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, providedRequestID, capturedRequestID)
}

// TestLoggerMiddleware tests basic logging functionality
func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.Use(LoggerMiddleware())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	logOutput := captureLogOutput(func() {
		router.ServeHTTP(w, req)
	})

	assert.NotEmpty(t, logOutput)
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "RequestID=")
}
