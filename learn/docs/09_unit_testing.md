# Unit Testing in Go - Complete Guide

**Target Audience:** Developers coming from Node.js/JavaScript background
**Focus:** Testing authentication layer by layer

---

## Table of Contents
1. [What is Unit Testing?](#what-is-unit-testing)
2. [Go Testing Basics](#go-testing-basics)
3. [Testing Pyramid](#testing-pyramid)
4. [Tools & Libraries](#tools--libraries)
5. [Testing Auth: Layer by Layer](#testing-auth-layer-by-layer)
6. [Mocking & Dependency Injection](#mocking--dependency-injection)
7. [Table-Driven Tests](#table-driven-tests)
8. [Test Coverage](#test-coverage)

---

## What is Unit Testing?

**Unit testing** = Testing a single "unit" (function, method, component) in isolation.

### In JavaScript/Node.js:
```javascript
// Using Jest
describe('add function', () => {
  test('adds 1 + 2 to equal 3', () => {
    expect(add(1, 2)).toBe(3);
  });
});
```

### In Go:
```go
// Using built-in testing package
func TestAdd(t *testing.T) {
    result := add(1, 2)
    if result != 3 {
        t.Errorf("add(1, 2) = %d; want 3", result)
    }
}
```

---

## Go Testing Basics

### 1. Test File Naming Convention

**RULE:** Test files MUST end with `_test.go`

```
project/
├── auth.go           # Your code
├── auth_test.go      # Tests for auth.go
├── jwt.go
└── jwt_test.go       # Tests for jwt.go
```

**Why?** Go compiler IGNORES `_test.go` files when building your app. They only run with `go test`.

---

### 2. Test Function Naming

**RULE:** Test functions MUST start with `Test`

```go
// ✅ CORRECT
func TestLogin(t *testing.T) { }
func TestRegister(t *testing.T) { }
func TestJWTGeneration(t *testing.T) { }

// ❌ WRONG - Go won't recognize these as tests
func login_test(t *testing.T) { }
func testLogin(t *testing.T) { }   // lowercase 't'
```

**Internal Check:** Go's test runner scans for functions matching this pattern:
```
func Test<Name>(t *testing.T)
```

---

### 3. The `testing.T` Type

`t *testing.T` is Go's test context object (like `assert` in Jest)

**Common Methods:**

| Method | Purpose | Example |
|--------|---------|---------|
| `t.Error()` | Mark test as FAILED (continue running) | `t.Error("login failed")` |
| `t.Errorf()` | Formatted error message | `t.Errorf("want %d, got %d", 5, result)` |
| `t.Fatal()` | Mark FAILED and STOP immediately | `t.Fatal("cannot continue")` |
| `t.Fatalf()` | Formatted fatal error | `t.Fatalf("setup failed: %v", err)` |
| `t.Skip()` | Skip this test | `t.Skip("not implemented yet")` |
| `t.Log()` | Print debug message (only if test fails) | `t.Log("user created")` |

**Example:**
```go
func TestDivide(t *testing.T) {
    result, err := divide(10, 0)

    // Fatal stops execution - use for critical failures
    if err == nil {
        t.Fatal("expected error for division by zero")
        // Code after t.Fatal() won't execute
    }

    // Error continues execution - use for assertions
    if result != 0 {
        t.Errorf("want 0, got %d", result)
        // Test continues even after failure
    }
}
```

---

### 4. Running Tests

```bash
# Run ALL tests in current directory
go test

# Run tests in ALL subdirectories
go test ./...

# Run with verbose output (see all test names)
go test -v

# Run specific test by name
go test -run TestLogin

# Run tests matching pattern
go test -run "^TestAuth"    # All tests starting with TestAuth

# Show test coverage
go test -cover

# Generate detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Output Example:**
```bash
$ go test -v
=== RUN   TestLogin
--- PASS: TestLogin (0.00s)
=== RUN   TestRegister
--- PASS: TestRegister (0.01s)
PASS
ok      todo_app/auth   0.123s
```

---

## Testing Pyramid

Different types of tests serve different purposes:

```
        /\
       /  \
      / E2E \         E2E Tests (slow, test entire system)
     /______\
    /        \
   / Integration \   Integration Tests (test multiple components)
  /______________\
 /                \
/   UNIT TESTS     \  Unit Tests (fast, test single function)
____________________
```

**For Auth:**
- **Unit Tests:** Test JWT generation, password hashing in isolation
- **Integration Tests:** Test auth handler + service + database together
- **E2E Tests:** Test full login flow via HTTP requests

**Focus Today:** Unit tests (fast, reliable, no database needed)

---

## Tools & Libraries

### 1. Built-in `testing` Package

**Pros:** No dependencies, always available
**Cons:** Verbose, manual assertions

```go
func TestAdd(t *testing.T) {
    result := add(2, 3)
    if result != 5 {
        t.Errorf("got %d, want 5", result)
    }
}
```

### 2. Testify (Most Popular)

**Install:**
```bash
go get github.com/stretchr/testify
```

**Pros:** Clean assertions, easy mocking
**Cons:** External dependency

```go
import "github.com/stretchr/testify/assert"

func TestAdd(t *testing.T) {
    result := add(2, 3)
    assert.Equal(t, 5, result)  // Much cleaner!
}
```

**Common Testify Functions:**

| Function | Purpose |
|----------|---------|
| `assert.Equal(t, expected, actual)` | Check equality |
| `assert.NotEqual(t, expected, actual)` | Check inequality |
| `assert.Nil(t, value)` | Check if nil |
| `assert.NotNil(t, value)` | Check if not nil |
| `assert.True(t, condition)` | Check if true |
| `assert.False(t, condition)` | Check if false |
| `assert.Contains(t, string, substring)` | Check substring |
| `assert.NoError(t, err)` | Check err == nil |
| `assert.Error(t, err)` | Check err != nil |

---

## Testing Auth: Layer by Layer

Your auth system has these layers:

```
┌─────────────────────────────────────┐
│  1. HTTP Handler (auth_handler.go)  │  ← Receives HTTP requests
├─────────────────────────────────────┤
│  2. Service (user_service.go)       │  ← Business logic
├─────────────────────────────────────┤
│  3. Repository (user_repository.go) │  ← Database operations
├─────────────────────────────────────┤
│  4. Utilities (jwt.go, bcrypt)      │  ← Helper functions
└─────────────────────────────────────┘
```

**Testing Strategy:**
- Test each layer INDEPENDENTLY
- Use MOCKS to replace dependencies
- Start from bottom (utilities) → top (handlers)

---

## Layer 1: Testing Utilities (No Dependencies)

Utilities like JWT and password hashing have NO dependencies → easiest to test!

### Example: Testing JWT Generation

```go
// pkg/utils/jwt_test.go
package utils

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

// TestJWTGeneration tests creating a valid JWT token
func TestJWTGeneration(t *testing.T) {
    // Setup: Create JWTUtil with test config
    jwtUtil := NewJWTUtil("test-secret-key", 24, "test-issuer")

    // Action: Generate token
    token, expiresAt, err := jwtUtil.GenerateToken("user-123", "john")

    // Assert: Check results
    assert.NoError(t, err, "token generation should not error")
    assert.NotEmpty(t, token, "token should not be empty")
    assert.Greater(t, expiresAt, time.Now().Unix(), "expires should be in future")
}

// TestJWTValidation tests validating a token
func TestJWTValidation(t *testing.T) {
    jwtUtil := NewJWTUtil("test-secret", 24, "test-issuer")

    // Generate a token first
    token, _, err := jwtUtil.GenerateToken("user-123", "john")
    assert.NoError(t, err)

    // Validate it
    claims, err := jwtUtil.ValidateToken(token)

    // Check validation succeeded
    assert.NoError(t, err)
    assert.Equal(t, "user-123", claims.UserID)
    assert.Equal(t, "john", claims.Username)
    assert.Equal(t, "test-issuer", claims.Issuer)
}

// TestJWTValidationWithWrongSecret tests security
func TestJWTValidationWithWrongSecret(t *testing.T) {
    // Create token with one secret
    jwtUtil1 := NewJWTUtil("secret-1", 24, "issuer")
    token, _, _ := jwtUtil1.GenerateToken("user-123", "john")

    // Try to validate with DIFFERENT secret
    jwtUtil2 := NewJWTUtil("secret-2", 24, "issuer")
    claims, err := jwtUtil2.ValidateToken(token)

    // Should FAIL
    assert.Error(t, err, "validation should fail with wrong secret")
    assert.Nil(t, claims, "claims should be nil for invalid token")
}

// TestExpiredToken tests token expiration
func TestExpiredToken(t *testing.T) {
    // Create token that expires in 0 hours (immediately)
    jwtUtil := NewJWTUtil("secret", 0, "issuer")
    token, _, _ := jwtUtil.GenerateToken("user-123", "john")

    // Wait a moment
    time.Sleep(time.Second)

    // Validation should fail
    claims, err := jwtUtil.ValidateToken(token)
    assert.Error(t, err)
    assert.Nil(t, claims)
}
```

**Run it:**
```bash
cd pkg/utils
go test -v
```

---

## Layer 2: Testing Middleware (With Mocks)

Middleware depends on:
1. HTTP request/response (use `httptest`)
2. JWTUtil (can use real one - it's just a utility)

### Example: Testing Auth Middleware

```go
// api/middleware/auth_test.go
package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "todo_app/pkg/utils"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// TestAuthMiddlewareNoToken tests missing token
func TestAuthMiddlewareNoToken(t *testing.T) {
    // Setup Gin in test mode
    gin.SetMode(gin.TestMode)

    // Create test router with middleware
    router := gin.New()
    jwtUtil := utils.NewJWTUtil("test-secret", 24, "test")
    router.Use(AuthMiddleware(jwtUtil))

    // Add protected endpoint
    router.GET("/protected", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Create request WITHOUT Authorization header
    req := httptest.NewRequest("GET", "/protected", nil)
    w := httptest.NewRecorder()

    // Execute
    router.ServeHTTP(w, req)

    // Assert: Should get 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddlewareInvalidFormat tests wrong header format
func TestAuthMiddlewareInvalidFormat(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    jwtUtil := utils.NewJWTUtil("test-secret", 24, "test")
    router.Use(AuthMiddleware(jwtUtil))

    router.GET("/protected", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Create request with WRONG format (missing "Bearer")
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "just-a-token")  // ❌ Should be "Bearer <token>"
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    // Should fail
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddlewareValidToken tests successful authentication
func TestAuthMiddlewareValidToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    jwtUtil := utils.NewJWTUtil("test-secret", 24, "test")
    router.Use(AuthMiddleware(jwtUtil))

    router.GET("/protected", func(c *gin.Context) {
        // Check if user data was set in context
        userID := c.GetString("user_id")
        username := c.GetString("username")
        c.JSON(200, gin.H{
            "message": "success",
            "user_id": userID,
            "username": username,
        })
    })

    // Generate VALID token
    token, _, _ := jwtUtil.GenerateToken("user-123", "john")

    // Create request with valid token
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)  // ✅ Correct format
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    // Should succeed
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**What's Happening Here:**

1. **httptest.NewRequest()** - Creates fake HTTP request (no real server needed!)
2. **httptest.NewRecorder()** - Records response (like Postman, but in code)
3. **router.ServeHTTP()** - Simulates handling the request
4. **w.Code** - Check HTTP status code
5. **w.Body** - Check response body

---

## Layer 3: Testing Handlers (With Service Mocks)

Handlers depend on:
1. HTTP context (use `httptest`)
2. Service layer (use MOCK)

### Why Mock the Service?

**WITHOUT Mock:**
```
Handler → Service → Repository → Database
                                   ↑
                            Need real database!
```

**WITH Mock:**
```
Handler → MockService (returns fake data)
          ↑
     No database needed!
```

### Example: Testing Auth Handler

```go
// api/handler/auth_handler_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    "todo_app/internal/dto"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// MockUserService is a FAKE service that implements UserService interface
type MockUserService struct {
    // Control what the mock returns
    ShouldReturnError bool
    ReturnedResponse  *dto.LoginResponse
}

// Register implements service.UserService interface
func (m *MockUserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
    if m.ShouldReturnError {
        return nil, errors.New("mock error")
    }

    // Return fake data
    return &dto.LoginResponse{
        Token: "mock-token-12345",
        User: dto.UserResponse{
            ID:       "mock-user-id",
            Username: "mock-username",
            Email:    "mock@example.com",
        },
        ExpiresAt: time.Now().Unix() + 3600,
    }, nil
}

// Login implements service.UserService interface
func (m *MockUserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
    if m.ShouldReturnError {
        return nil, errors.New("invalid credentials")
    }

    return m.ReturnedResponse, nil
}

// TestRegisterSuccess tests successful registration
func TestRegisterSuccess(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)

    // Create mock service (NO database needed!)
    mockService := &MockUserService{
        ShouldReturnError: false,
    }

    // Create handler with mock
    handler := NewAuthHandler(mockService)

    // Create test request
    reqBody := dto.RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
        FullName: "Test User",
    }
    bodyJSON, _ := json.Marshal(reqBody)

    // Create Gin context
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(bodyJSON))
    c.Request.Header.Set("Content-Type", "application/json")

    // Execute
    handler.Register(c)

    // Assert
    assert.Equal(t, 201, w.Code)

    var response dto.LoginResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "mock-token-12345", response.Token)
}

// TestRegisterInvalidJSON tests bad request body
func TestRegisterInvalidJSON(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockService := &MockUserService{}
    handler := NewAuthHandler(mockService)

    // Send INVALID JSON
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/register", bytes.NewBufferString("{invalid"))
    c.Request.Header.Set("Content-Type", "application/json")

    handler.Register(c)

    // Should return 400 Bad Request
    assert.Equal(t, 400, w.Code)
}

// TestRegisterServiceError tests service layer error
func TestRegisterServiceError(t *testing.T) {
    gin.SetMode(gin.TestMode)

    // Configure mock to return error
    mockService := &MockUserService{
        ShouldReturnError: true,
    }
    handler := NewAuthHandler(mockService)

    reqBody := dto.RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    bodyJSON, _ := json.Marshal(reqBody)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(bodyJSON))
    c.Request.Header.Set("Content-Type", "application/json")

    handler.Register(c)

    // Should return 400 (as per your handler code)
    assert.Equal(t, 400, w.Code)
}
```

---

## Table-Driven Tests

Testing multiple scenarios gets repetitive. **Table-driven tests** solve this:

### Without Table:
```go
func TestAdd(t *testing.T) {
    if add(1, 2) != 3 { t.Fail() }
}
func TestAddNegative(t *testing.T) {
    if add(-1, 1) != 0 { t.Fail() }
}
func TestAddZero(t *testing.T) {
    if add(0, 0) != 0 { t.Fail() }
}
```

### With Table (Better!):
```go
func TestAdd(t *testing.T) {
    // Define test cases
    tests := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"positive numbers", 1, 2, 3},
        {"negative numbers", -1, 1, 0},
        {"both zero", 0, 0, 0},
        {"large numbers", 1000, 2000, 3000},
    }

    // Run each test case
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := add(tt.a, tt.b)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**Output:**
```bash
=== RUN   TestAdd
=== RUN   TestAdd/positive_numbers
=== RUN   TestAdd/negative_numbers
=== RUN   TestAdd/both_zero
=== RUN   TestAdd/large_numbers
--- PASS: TestAdd (0.00s)
```

### Real Example: Testing JWT Validation

```go
func TestJWTValidation(t *testing.T) {
    tests := []struct {
        name        string
        secret      string
        createToken bool
        useWrongSecret bool
        wantError   bool
    }{
        {
            name:        "valid token",
            secret:      "test-secret",
            createToken: true,
            useWrongSecret: false,
            wantError:   false,
        },
        {
            name:        "wrong secret",
            secret:      "test-secret",
            createToken: true,
            useWrongSecret: true,
            wantError:   true,
        },
        {
            name:        "invalid token string",
            secret:      "test-secret",
            createToken: false,
            wantError:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jwtUtil := NewJWTUtil(tt.secret, 24, "issuer")

            var token string
            if tt.createToken {
                token, _, _ = jwtUtil.GenerateToken("user-123", "john")
            } else {
                token = "invalid-token-string"
            }

            if tt.useWrongSecret {
                jwtUtil = NewJWTUtil("different-secret", 24, "issuer")
            }

            _, err := jwtUtil.ValidateToken(token)

            if tt.wantError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## Test Coverage

Check how much of your code is tested:

```bash
# Show coverage percentage
go test -cover

# Output:
# PASS
# coverage: 75.0% of statements

# Generate detailed report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

This opens an HTML file showing:
- **Green lines:** Covered by tests ✅
- **Red lines:** NOT covered by tests ❌
- **Gray lines:** Not executable (comments, declarations)

**Good Coverage Targets:**
- **Utilities:** 90%+ (easy to test)
- **Handlers:** 70%+ (harder due to HTTP context)
- **Services:** 80%+ (business logic)

---

## Quick Reference: Testing Patterns

### Pattern 1: AAA (Arrange-Act-Assert)

```go
func TestLogin(t *testing.T) {
    // ARRANGE: Setup test data
    mockService := &MockUserService{}
    handler := NewAuthHandler(mockService)
    reqBody := dto.LoginRequest{Username: "test", Password: "pass"}

    // ACT: Execute the code
    response, err := handler.Login(reqBody)

    // ASSERT: Verify results
    assert.NoError(t, err)
    assert.NotNil(t, response)
}
```

### Pattern 2: Setup/Teardown

```go
func setupTest(t *testing.T) (*AuthHandler, func()) {
    // Setup
    mockService := &MockUserService{}
    handler := NewAuthHandler(mockService)

    // Return cleanup function
    cleanup := func() {
        // Clean up resources
    }

    return handler, cleanup
}

func TestSomething(t *testing.T) {
    handler, cleanup := setupTest(t)
    defer cleanup()  // Always run cleanup

    // Your test...
}
```

### Pattern 3: Parallel Tests

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name string
        // ...
    }{ /* test cases */ }

    for _, tt := range tests {
        tt := tt  // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // Run tests concurrently
            // Your test...
        })
    }
}
```

---

## Summary

### Testing Checklist for Auth:

- [ ] **JWT Utilities**
  - [ ] Token generation
  - [ ] Token validation
  - [ ] Expiration handling
  - [ ] Wrong secret rejection

- [ ] **Auth Middleware**
  - [ ] Missing token
  - [ ] Invalid format
  - [ ] Valid token
  - [ ] Expired token

- [ ] **Auth Handler**
  - [ ] Successful registration
  - [ ] Invalid JSON
  - [ ] Service errors
  - [ ] Successful login
  - [ ] Invalid credentials

- [ ] **Service Layer** (Next step!)
  - [ ] Password hashing
  - [ ] User creation
  - [ ] Duplicate username check
  - [ ] Login validation

---

## Next Steps

1. Write tests for `jwt.go` utilities
2. Write tests for `auth.go` middleware
3. Write tests for `auth_handler.go`
4. Create mocks for service layer
5. Run coverage report
6. Aim for 70%+ coverage

**Commands to remember:**
```bash
go test -v              # Run tests verbosely
go test -run TestLogin  # Run specific test
go test -cover          # Show coverage
go test ./...           # Test all packages
```

---

**Your Turn!** Let's start writing actual tests for your auth code. Which layer do you want to test first?
