package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupCORSTest creates a test router with CORS middleware
func setupCORSTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	return router
}

// TestCORSMiddlewareWithAllowedOrigin tests CORS with allowed origin
func TestCORSMiddlewareWithAllowedOrigin(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Test allowed origins
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:3000",
	}

	for _, origin := range allowedOrigins {
		t.Run("allowed origin: "+origin, func(t *testing.T) {
			// ARRANGE
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", origin)
			w := httptest.NewRecorder()

			// ACT
			router.ServeHTTP(w, req)

			// ASSERT: Should set CORS headers
			assert.Equal(t, 200, w.Code)
			assert.Equal(t, origin, w.Header().Get("Access-Control-Allow-Origin"),
				"should echo back the allowed origin")
			assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"),
				"should allow credentials")
		})
	}
}

// TestCORSMiddlewareWithDisallowedOrigin tests CORS with disallowed origin
func TestCORSMiddlewareWithDisallowedOrigin(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	disallowedOrigins := []string{
		"http://evil.com",
		"https://localhost:5173", // https not allowed (only http in dev)
		"http://localhost:8080",  // different port
		"http://example.com",
	}

	for _, origin := range disallowedOrigins {
		t.Run("disallowed origin: "+origin, func(t *testing.T) {
			// ARRANGE
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", origin)
			w := httptest.NewRecorder()

			// ACT
			router.ServeHTTP(w, req)

			// ASSERT: Should NOT set Access-Control-Allow-Origin
			assert.Equal(t, 200, w.Code, "request should still succeed")
			assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"),
				"should not set CORS origin header for disallowed origin")
			assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"),
				"should not set credentials header for disallowed origin")
		})
	}
}

// TestCORSMiddlewareWithNoOrigin tests CORS when no Origin header is sent
func TestCORSMiddlewareWithNoOrigin(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// Don't set Origin header
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Should still set common headers but not origin-specific ones
	assert.Equal(t, 200, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"),
		"should not set origin when none provided")

	// Common headers should still be set
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

// TestCORSPreflightRequest tests OPTIONS preflight requests
func TestCORSPreflightRequest(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/api/users", func(c *gin.Context) {
		c.JSON(200, gin.H{"users": []string{}})
	})

	req := httptest.NewRequest("OPTIONS", "/api/users", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Authorization")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Should return 204 No Content for preflight
	assert.Equal(t, 204, w.Code, "preflight should return 204")
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
}

// TestCORSAllowedMethods tests that correct methods are allowed
func TestCORSAllowedMethods(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Check allowed methods header
	allowedMethods := w.Header().Get("Access-Control-Allow-Methods")
	assert.Contains(t, allowedMethods, "GET")
	assert.Contains(t, allowedMethods, "POST")
	assert.Contains(t, allowedMethods, "PUT")
	assert.Contains(t, allowedMethods, "PATCH")
	assert.Contains(t, allowedMethods, "DELETE")
	assert.Contains(t, allowedMethods, "OPTIONS")
}

// TestCORSAllowedHeaders tests that correct headers are allowed
func TestCORSAllowedHeaders(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Check allowed headers
	allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "Authorization")
	assert.Contains(t, allowedHeaders, "Accept-Encoding")
}

// TestCORSSecurityHeaders tests security headers
func TestCORSSecurityHeaders(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Check security headers
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"),
		"should prevent clickjacking")
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"),
		"should prevent MIME type sniffing")
}

// TestCORSPreflightWithDisallowedOrigin tests preflight with disallowed origin
func TestCORSPreflightWithDisallowedOrigin(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()

	req := httptest.NewRequest("OPTIONS", "/api/users", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Should still return 204 but without origin header
	assert.Equal(t, 204, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"),
		"should not set origin for disallowed origin")
}

// TestCORSWithPOSTRequest tests CORS with POST request
func TestCORSWithPOSTRequest(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.POST("/api/users", func(c *gin.Context) {
		c.JSON(201, gin.H{"id": "123"})
	})

	req := httptest.NewRequest("POST", "/api/users", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

// TestCORSMiddlewareDoesNotBlockRequest tests that middleware passes request through
func TestCORSMiddlewareDoesNotBlockRequest(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	handlerCalled := false
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "handler called"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Handler should be called
	assert.True(t, handlerCalled, "handler should be called after middleware")
	assert.Equal(t, 200, w.Code)
}

// TestCORSPreflightDoesNotCallHandler tests that OPTIONS stops at middleware
func TestCORSPreflightDoesNotCallHandler(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	handlerCalled := false
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "handler called"})
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Handler should NOT be called for preflight
	assert.False(t, handlerCalled, "handler should not be called for OPTIONS")
	assert.Equal(t, 204, w.Code)
}

// TestCORSWithMultipleRequests tests CORS with multiple sequential requests
func TestCORSWithMultipleRequests(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	origins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
	}

	// ACT & ASSERT: Multiple requests should each get proper CORS headers
	for _, origin := range origins {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", origin)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, origin, w.Header().Get("Access-Control-Allow-Origin"))
	}
}

// TestCORSHeadersAlwaysPresent tests that common headers are always set
func TestCORSHeadersAlwaysPresent(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// No Origin header
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Common headers should always be present
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}

// TestCORSWithCredentials tests credentials flag is properly set
func TestCORSWithCredentials(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Credentials should be "true" (string, not boolean)
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

// TestCORSWithCaseSensitiveOrigin tests origin matching is case-sensitive
func TestCORSWithCaseSensitiveOrigin(t *testing.T) {
	// ARRANGE
	router := setupCORSTest()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Test with different case (should not match)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://LocalHost:5173") // Different case
	w := httptest.NewRecorder()

	// ACT
	router.ServeHTTP(w, req)

	// ASSERT: Should NOT match (origin matching is case-sensitive)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
