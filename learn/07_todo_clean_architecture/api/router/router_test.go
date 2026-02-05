package router

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"todo_app/api/handler/mocks"

	"github.com/stretchr/testify/assert"
)

// TestRegisterEndpoint tests the /auth/register endpoint
// NOTE: Uses MOCK handler - NO database needed!
// NOTE: This is called a "UNIT TEST" - tests only the router
func TestRegisterEndpoint(t *testing.T) {
	// Step 1: Create MOCK handlers
	// NOTE: These implement the interfaces but don't use database
	mockAuthHandler := mocks.NewMockAuthHandler()
	mockUserHandler := mocks.NewMockUserHandler()
	mockTodoHandler := mocks.NewMockTodoHandler()
	mockListHandler := mocks.NewMockTodoListHandler()

	// Step 2: Create router with MOCK handlers
	// NOTE: This works because SetupRouter accepts INTERFACES, not concrete types
	// NOTE: Router doesn't know (or care) if handlers are real or mock!
	router := SetupRouter(mockAuthHandler, mockUserHandler, mockTodoHandler, mockListHandler, nil)

	// Step 3: Create a fake HTTP request
	requestBody := map[string]string{
		"username":  "testuser",
		"email":     "test@example.com",
		"password":  "TestPass123",
		"full_name": "Test User",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	// Step 4: Record the response
	w := httptest.NewRecorder()

	// Step 5: Call the router
	// NOTE: Router will call mockAuthHandler.Register() instead of real handler
	router.ServeHTTP(w, req)

	// Step 6: Verify the response
	assert.Equal(t, 201, w.Code, "Expected 201 Created status")

	// NOTE: Verify mock was called
	assert.True(t, mockAuthHandler.RegisterCalled, "Expected Register to be called")
	assert.Equal(t, 1, mockAuthHandler.RegisterCount, "Expected Register called once")

	// NOTE: Verify response body contains expected data
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Check fake data returned by mock
	assert.Equal(t, "mock-token-12345", response["token"])
	assert.NotNil(t, response["user"])

	user := response["user"].(map[string]interface{})
	assert.Equal(t, "mock-user-id", user["id"])
	assert.Equal(t, "mock-username", user["username"])
}

// TestLoginEndpoint tests the /auth/login endpoint
func TestLoginEndpoint(t *testing.T) {
	mockAuthHandler := mocks.NewMockAuthHandler()
	mockUserHandler := mocks.NewMockUserHandler()
	mockTodoHandler := mocks.NewMockTodoHandler()
	mockListHandler := mocks.NewMockTodoListHandler()

	router := SetupRouter(mockAuthHandler, mockUserHandler, mockTodoHandler, mockListHandler, nil)

	// Create login request
	requestBody := map[string]string{
		"username": "testuser",
		"password": "TestPass123",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, 200, w.Code)
	assert.True(t, mockAuthHandler.LoginCalled, "Expected Login to be called")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "mock-token-67890", response["token"])
}

// TestRegisterError tests error handling
// NOTE: We can control what the mock returns!
func TestRegisterError(t *testing.T) {
	mockAuthHandler := mocks.NewMockAuthHandler()
	mockUserHandler := mocks.NewMockUserHandler()
	mockTodoHandler := mocks.NewMockTodoHandler()
	mockListHandler := mocks.NewMockTodoListHandler()

	// NOTE: Configure mock to return error
	mockAuthHandler.ShouldReturnError = true

	router := SetupRouter(mockAuthHandler, mockUserHandler, mockTodoHandler, mockListHandler, nil)

	requestBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "TestPass123",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// NOTE: Verify error response
	assert.Equal(t, 400, w.Code, "Expected 400 Bad Request")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "mock error")
}

// TestMultipleRequests tests that router can handle multiple calls
func TestMultipleRequests(t *testing.T) {
	mockAuthHandler := mocks.NewMockAuthHandler()
	mockUserHandler := mocks.NewMockUserHandler()
	mockTodoHandler := mocks.NewMockTodoHandler()
	mockListHandler := mocks.NewMockTodoListHandler()

	router := SetupRouter(mockAuthHandler, mockUserHandler, mockTodoHandler, mockListHandler, nil)

	// Make 5 register requests
	for i := 0; i < 5; i++ {
		requestBody := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "TestPass123",
		}
		bodyJSON, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	}

	// NOTE: Mock tracked all calls
	assert.Equal(t, 5, mockAuthHandler.RegisterCount, "Expected 5 calls to Register")
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	mockAuthHandler := mocks.NewMockAuthHandler()
	mockUserHandler := mocks.NewMockUserHandler()
	mockTodoHandler := mocks.NewMockTodoHandler()
	mockListHandler := mocks.NewMockTodoListHandler()

	router := SetupRouter(mockAuthHandler, mockUserHandler, mockTodoHandler, mockListHandler, nil)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "todo-api", response["service"])
}

// TestRouteNotFound tests 404 handling
func TestRouteNotFound(t *testing.T) {
	mockAuthHandler := mocks.NewMockAuthHandler()
	mockUserHandler := mocks.NewMockUserHandler()
	mockTodoHandler := mocks.NewMockTodoHandler()
	mockListHandler := mocks.NewMockTodoListHandler()

	router := SetupRouter(mockAuthHandler, mockUserHandler, mockTodoHandler, mockListHandler, nil)

	// Request non-existent route
	req := httptest.NewRequest("POST", "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 404
	assert.Equal(t, 404, w.Code)

	// NOTE: Mock handlers should NOT be called
	assert.False(t, mockAuthHandler.RegisterCalled)
	assert.False(t, mockAuthHandler.LoginCalled)
}
