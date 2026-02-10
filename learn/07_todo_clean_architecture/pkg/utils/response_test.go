package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestContext creates a test Gin context with response recorder
// This is a helper function used by all tests
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	// Set Gin to test mode (disables debug output)
	gin.SetMode(gin.TestMode)

	// Create response recorder to capture HTTP response
	w := httptest.NewRecorder()

	// Create test context
	c, _ := gin.CreateTestContext(w)

	return c, w
}

// parseResponse is a helper to decode JSON response
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "should parse JSON response")
	return resp
}

// TestSuccess tests the Success response helper
func TestSuccess(t *testing.T) {
	// ARRANGE: Setup test context and data
	c, w := setupTestContext()
	testData := map[string]string{"message": "operation successful"}

	// ACT: Call Success
	Success(c, testData)

	// ASSERT: Check status code
	assert.Equal(t, http.StatusOK, w.Code, "should return 200 OK")

	// ASSERT: Check response body
	resp := parseResponse(t, w)
	assert.True(t, resp.Success, "success should be true")
	assert.Empty(t, resp.Error, "error should be empty")
	assert.NotNil(t, resp.Data, "data should not be nil")
}

// TestCreated tests the Created response helper
func TestCreated(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()
	testData := map[string]interface{}{
		"id":   "123",
		"name": "New Item",
	}

	// ACT
	Created(c, testData)

	// ASSERT: Check status code
	assert.Equal(t, http.StatusCreated, w.Code, "should return 201 Created")

	// ASSERT: Check response body
	resp := parseResponse(t, w)
	assert.True(t, resp.Success, "success should be true")
	assert.Empty(t, resp.Error, "error should be empty")
	assert.NotNil(t, resp.Data, "data should not be nil")
}

// TestBadRequest tests the BadRequest response helper
func TestBadRequest(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()
	errorMessage := "invalid input provided"

	// ACT
	BadRequest(c, errorMessage)

	// ASSERT: Check status code
	assert.Equal(t, http.StatusBadRequest, w.Code, "should return 400 Bad Request")

	// ASSERT: Check response body
	resp := parseResponse(t, w)
	assert.False(t, resp.Success, "success should be false")
	assert.Equal(t, errorMessage, resp.Error, "error message should match")
	assert.Nil(t, resp.Data, "data should be nil")
}

// TestUnauthorized tests the Unauthorized response helper
func TestUnauthorized(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()
	errorMessage := "authentication required"

	// ACT
	Unauthorized(c, errorMessage)

	// ASSERT: Check status code
	assert.Equal(t, http.StatusUnauthorized, w.Code, "should return 401 Unauthorized")

	// ASSERT: Check response body
	resp := parseResponse(t, w)
	assert.False(t, resp.Success, "success should be false")
	assert.Equal(t, errorMessage, resp.Error, "error message should match")
}

// TestSuccessWithNilData tests Success with nil data
func TestSuccessWithNilData(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()

	// ACT: Pass nil as data
	Success(c, nil)

	// ASSERT: Should still work
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp.Success)
}

// TestSuccessWithComplexData tests Success with nested data structures
func TestSuccessWithComplexData(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()
	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    123,
			"name":  "John Doe",
			"roles": []string{"admin", "user"},
		},
		"metadata": map[string]interface{}{
			"page":  1,
			"total": 100,
		},
	}

	// ACT
	Success(c, complexData)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

// TestErrorResponsesWithEmptyMessages tests error helpers with empty strings
func TestErrorResponsesWithEmptyMessages(t *testing.T) {
	tests := []struct {
		name         string
		errorFunc    func(*gin.Context, string)
		expectedCode int
	}{
		{"BadRequest empty", BadRequest, http.StatusBadRequest},
		{"Unauthorized empty", Unauthorized, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			c, w := setupTestContext()

			// ACT: Call with empty message
			tt.errorFunc(c, "")

			// ASSERT
			assert.Equal(t, tt.expectedCode, w.Code)
			resp := parseResponse(t, w)
			assert.False(t, resp.Success)
			assert.Empty(t, resp.Error, "error should be empty string")
		})
	}
}

// TestResponseStructure tests the Response struct JSON serialization
func TestResponseStructure(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		expected map[string]interface{}
	}{
		{
			name: "success with data",
			response: Response{
				Success: true,
				Data:    "test data",
			},
			expected: map[string]interface{}{
				"success": true,
				"data":    "test data",
			},
		},
		{
			name: "error response",
			response: Response{
				Success: false,
				Error:   "test error",
			},
			expected: map[string]interface{}{
				"success": false,
				"error":   "test error",
			},
		},
		{
			name: "complete response",
			response: Response{
				Success: true,
				Message: "operation completed",
				Data:    map[string]string{"key": "value"},
			},
			expected: map[string]interface{}{
				"success": true,
				"message": "operation completed",
				"data": map[string]interface{}{
					"key": "value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT: Marshal to JSON
			jsonData, err := json.Marshal(tt.response)
			assert.NoError(t, err)

			// Unmarshal to map for comparison
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			assert.NoError(t, err)

			// ASSERT: Check fields exist correctly
			assert.Equal(t, tt.expected["success"], result["success"])
		})
	}
}

// TestContentTypeHeader tests that responses set correct content-type
func TestContentTypeHeader(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()

	// ACT
	Success(c, map[string]string{"test": "data"})

	// ASSERT: Should have JSON content type
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "should set JSON content type")
}

// TestMultipleResponseCalls tests calling response helper multiple times
// (In real scenarios, only one should be called per request)
func TestMultipleResponseCalls(t *testing.T) {
	// ARRANGE
	c, w := setupTestContext()

	// ACT: Call Success twice (this creates invalid JSON - multiple objects)
	Success(c, "first")
	Success(c, "second")

	// ASSERT: Response code should be set
	// NOTE: In real code, calling response helpers multiple times is a bug
	// NOTE: Multiple JSON writes create invalid JSON, but status code is still set
	assert.Equal(t, 200, w.Code, "should have status code from last call")

	// Body will have invalid JSON (two separate JSON objects)
	// This is expected behavior when misusing the API
}

