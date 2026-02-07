package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter creates a test router with auth middleware
// Returns: router and jwtUtil for creating tokens
func setupTestRouter(jwtUtil *utils.JWTUtil) *gin.Engine {
	// Set Gin to test mode (disables debug logs)
	gin.SetMode(gin.TestMode)

	// Create router with auth middleware
	router := gin.New()
	router.Use(AuthMiddleware(jwtUtil))

	// Add a protected endpoint
	router.GET("/protected", func(c *gin.Context) {
		// This endpoint only runs if middleware passes
		userID := c.MustGet("user_id").(uuid.UUID)
		username := c.MustGet("username").(string)

		c.JSON(http.StatusOK, gin.H{
			"message":  "success",
			"user_id":  userID.String(),
			"username": username,
		})
	})

	return router
}

// TestAuthMiddlewareNoToken tests request without Authorization header
func TestAuthMiddlewareNoToken(t *testing.T) {
	// ARRANGE: Setup router with middleware
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	router := setupTestRouter(jwtUtil)

	// ACT: Make request WITHOUT Authorization header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Should get 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// NOTE: Middleware calls c.Abort() when token is missing
	// NOTE: This prevents the handler from executing at all
}

// TestAuthMiddlewareEmptyToken tests Authorization header with empty value
// Expected: Should reject with 401
func TestAuthMiddlewareEmptyToken(t *testing.T) {
	// ARRANGE
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	router := setupTestRouter(jwtUtil)

	// ACT: Send empty Authorization header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "") // Empty!
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddlewareInvalidToken tests with malformed JWT tokens
// Expected: Should reject tokens that can't be parsed
func TestAuthMiddlewareInvalidToken(t *testing.T) {
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	router := setupTestRouter(jwtUtil)

	// Test various invalid tokens
	invalidTokens := []string{
		"not-a-jwt-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", // Only header part
		"invalid.token.format",
		"",
	}

	for _, token := range invalidTokens {
		// ACT: Send request with invalid token
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ASSERT: Should reject
		assert.Equal(t, http.StatusUnauthorized, w.Code, "token '%s' should be rejected", token)
	}

	// NOTE: jwtUtil.ValidateToken() will return error for these
	// NOTE: Middleware catches error and calls c.Abort()
}

// TestAuthMiddlewareExpiredToken tests with expired token
// Expected: Should reject old tokens
func TestAuthMiddlewareExpiredToken(t *testing.T) {
	// ARRANGE: Create token that expires immediately
	jwtUtil := utils.NewJWTUtil("test-secret", 0, "test-issuer")
	token, _, _ := jwtUtil.GenerateToken("user-123", "john")

	// Wait for expiration
	time.Sleep(time.Second * 2)

	router := setupTestRouter(jwtUtil)

	// ACT: Try to use expired token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Should reject
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// NOTE: ValidateToken checks ExpiresAt.Before(time.Now())
	// NOTE: This is why tokens have expiration - security best practice
}

// TestAuthMiddlewareWrongSecret tests token signed with different secret
// Expected: Should reject (security test)
func TestAuthMiddlewareWrongSecret(t *testing.T) {
	// ARRANGE: Create token with one secret
	jwtUtil1 := utils.NewJWTUtil("secret-key-1", 24, "test-issuer")
	token, _, _ := jwtUtil1.GenerateToken("user-123", "john")

	// Create router with DIFFERENT secret
	jwtUtil2 := utils.NewJWTUtil("secret-key-2", 24, "test-issuer")
	router := setupTestRouter(jwtUtil2)

	// ACT: Try to use token from different app
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Should reject
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// NOTE: This is critical security - prevents token reuse across apps
	// NOTE: JWT signature validation fails when secrets don't match
}

// TestAuthMiddlewareValidToken tests successful authentication
// Expected: Request should reach handler with user data in context
func TestAuthMiddlewareValidToken(t *testing.T) {
	// ARRANGE: Create valid token
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	userID := uuid.New().String()
	username := "john_doe"
	token, _, err := jwtUtil.GenerateToken(userID, username)
	assert.NoError(t, err, "token generation should succeed")

	router := setupTestRouter(jwtUtil)

	// ACT: Send request with valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token) // ✅ Correct format
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Should succeed
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body contains user data
	// WAY 1 (loose - uses map, no type safety):
	//   var response map[string]interface{}
	//   json.Unmarshal(w.Body.Bytes(), &response)
	//   assert.Equal(t, "success", response["message"])
	//
	// WAY 2 (better - uses struct, type safe!):
	type ProtectedResponse struct {
		Message  string `json:"message"`
		UserID   string `json:"user_id"`
		Username string `json:"username"`
	}
	var response ProtectedResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "response should be valid JSON")
	assert.Equal(t, "success", response.Message)    // Type-safe string!
	assert.NotEmpty(t, response.UserID)              // Autocomplete works!
	assert.Equal(t, "john_doe", response.Username)   // No typo risk!

	// NOTE: Middleware sets user_id (as uuid.UUID) and username (as string) in context
	// NOTE: Handlers access them with c.MustGet(constants.ContextUserID).(uuid.UUID)
}

// TestAuthMiddlewareContextValues tests that middleware sets correct context values
// Expected: user_id and username should be available in context
func TestAuthMiddlewareContextValues(t *testing.T) {
	// ARRANGE
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	// NOTE: user_id must be valid UUID format (middleware validates it!)
	userID := uuid.New().String()
	username := "jane_doe"
	token, _, _ := jwtUtil.GenerateToken(userID, username)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware(jwtUtil))

	// Variables to capture context values
	var capturedUserID uuid.UUID
	var capturedUsername string

	router.GET("/test", func(c *gin.Context) {
		// Capture values from context (middleware stores UUID type, not string)
		capturedUserID = c.MustGet("user_id").(uuid.UUID)
		capturedUsername = c.MustGet("username").(string)
		c.Status(http.StatusOK)
	})

	// ACT: Make authenticated request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Request should succeed
	assert.Equal(t, http.StatusOK, w.Code)

	// Context values should match
	assert.Equal(t, userID, capturedUserID.String(), "user_id in context should match token")
	assert.Equal(t, username, capturedUsername, "username in context should match token")

	// NOTE: Your middleware uses constants.ContextUserID and constants.ContextUsername
	// NOTE: This is good practice - avoids typos in key names
}

// TestAuthMiddlewareInvalidUserID tests when token has malformed user ID
// Expected: Should reject if user_id in token is not a valid UUID
func TestAuthMiddlewareInvalidUserID(t *testing.T) {
	// ARRANGE: Create JWT util and manually create token with invalid UUID
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")

	// Create token with invalid UUID (not using GenerateToken)
	claims := &utils.JWTClaims{
		UserID:   "not-a-valid-uuid", // ❌ Invalid!
		Username: "john",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "test-issuer",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	router := setupTestRouter(jwtUtil)

	// ACT: Try to use token with invalid UUID
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Should reject
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// NOTE: Middleware calls uuid.Parse(claims.UserID)
	// NOTE: This will error if UserID is not a valid UUID format
}

// TestAuthMiddlewareMultipleRequests tests that middleware works for multiple requests
// Expected: Each request should be handled independently
func TestAuthMiddlewareMultipleRequests(t *testing.T) {
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	router := setupTestRouter(jwtUtil)

	// Create different tokens for different users (use valid UUIDs!)
	users := []struct {
		id       string
		username string
	}{
		{uuid.New().String(), "alice"},
		{uuid.New().String(), "bob"},
		{uuid.New().String(), "charlie"},
	}

	for _, user := range users {
		// ACT: Make request with each user's token
		token, _, _ := jwtUtil.GenerateToken(user.id, user.username)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ASSERT: All should succeed
		assert.Equal(t, http.StatusOK, w.Code, "request for %s should succeed", user.username)
	}

	// NOTE: Middleware is stateless - doesn't remember previous requests
	// NOTE: Each request is authenticated independently
}

// TestAuthMiddlewareDifferentHTTPMethods tests middleware on various HTTP methods
// Expected: Middleware should work for GET, POST, PUT, DELETE, etc.
func TestAuthMiddlewareDifferentHTTPMethods(t *testing.T) {
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	// NOTE: Must use valid UUID format!
	token, _, _ := jwtUtil.GenerateToken(uuid.New().String(), "john")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware(jwtUtil))

	// Register handlers for different methods
	handler := func(c *gin.Context) { c.Status(http.StatusOK) }
	router.GET("/resource", handler)
	router.POST("/resource", handler)
	router.PUT("/resource", handler)
	router.DELETE("/resource", handler)
	router.PATCH("/resource", handler)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		// ACT: Make request with each method
		req := httptest.NewRequest(method, "/resource", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ASSERT: All should succeed
		assert.Equal(t, http.StatusOK, w.Code, "%s request should succeed", method)
	}

	// NOTE: Middleware doesn't care about HTTP method
	// NOTE: It only checks Authorization header
}

// TestAuthMiddlewareCallsNext tests that middleware calls c.Next() on success
// Expected: Handler chain should continue after successful auth
func TestAuthMiddlewareCallsNext(t *testing.T) {
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")
	// NOTE: Must use valid UUID!
	token, _, _ := jwtUtil.GenerateToken(uuid.New().String(), "john")

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Track if next middleware was called
	nextCalled := false

	router.Use(AuthMiddleware(jwtUtil))
	router.Use(func(c *gin.Context) {
		nextCalled = true
		c.Next()
	})
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// ACT: Make authenticated request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Next middleware should have been called
	assert.True(t, nextCalled, "c.Next() should have been called")

	// NOTE: When auth succeeds, middleware calls c.Next()
	// NOTE: This continues the middleware chain
}

// TestAuthMiddlewareAbortsOnFailure tests that middleware stops chain on failure
// Expected: c.Abort() should prevent subsequent middleware from running
func TestAuthMiddlewareAbortsOnFailure(t *testing.T) {
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Track if next middleware was called
	nextCalled := false

	router.Use(AuthMiddleware(jwtUtil))
	router.Use(func(c *gin.Context) {
		nextCalled = true // Should NOT be called
		c.Next()
	})
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// ACT: Make request WITHOUT token
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT: Next middleware should NOT have been called
	assert.False(t, nextCalled, "middleware chain should be aborted")
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// NOTE: When auth fails, middleware calls c.Abort()
	// NOTE: This stops execution - subsequent handlers don't run
}
