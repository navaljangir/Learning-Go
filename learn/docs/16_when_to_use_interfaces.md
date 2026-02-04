# When to Use Interfaces in Go

## TL;DR

**Go philosophy:** "Accept interfaces, return structs" + Define interfaces at the point of **use**, not definition.

**What needs interfaces:**
- ❌ **Entities** (User, Todo) - NO interfaces needed
- ✅ **Repositories** - YES (database boundary)
- ⚠️  **Services** - Usually NO (unless multiple implementations)
- ✅ **External services** - YES (email, payment, storage APIs)

---

## Your Current Architecture

```
Handler → Service (concrete) → Repository (interface) → Database
                                      ↑
                              THIS IS CORRECT!
```

**This is the right pattern for most Go applications.**

---

## Part 1: What Needs Interfaces?

### 1. Entities: NO ❌

Entities are data structures with domain logic - they should be **concrete types**.

```go
// ✅ GOOD: Concrete struct with methods
type User struct {
    ID           uuid.UUID
    Username     string
    Email        string
    PasswordHash string
    FullName     string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    *time.Time
}

func (u *User) IsDeleted() bool {
    return u.DeletedAt != nil
}

func (u *User) MarkDeleted() {
    now := time.Now()
    u.DeletedAt = &now
}

// ❌ BAD: Don't do this
type User interface {
    IsDeleted() bool
    MarkDeleted()
}
```

**Why no interface?**
- No external dependencies
- Easy to test directly (just create instances)
- No need for multiple implementations

**Testing:**
```go
func TestUser_IsDeleted(t *testing.T) {
    user := entity.NewUser("john", "john@example.com", "hash", "John")
    assert.False(t, user.IsDeleted())

    user.MarkDeleted()
    assert.True(t, user.IsDeleted())
}
```

---

### 2. Repositories: YES ✅

Repositories are the **boundary between your business logic and external systems** (database).

```go
// ✅ GOOD: Interface in domain layer
package repository

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    FindByUsername(ctx context.Context, username string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// Implementation in infrastructure layer
package sqlc_impl

type userRepository struct {
    db      *sql.DB
    queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
    return &userRepository{
        db:      db,
        queries: sqlc.New(db),
    }
}
```

**Why interface?**
- Database is external/slow
- Easy to swap implementations (PostgreSQL, MySQL, In-Memory)
- Easy to mock for testing services

**You already have this - correct! ✅**

---

### 3. Services: Usually NO ⚠️

Services coordinate business logic. **Most Go projects do NOT create service interfaces.**

```go
// ✅ GOOD: Your current approach
type UserService struct {
    userRepo repository.UserRepository  // ← Interface (for mocking)
    jwtUtil  *utils.JWTUtil
}

func NewUserService(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) *UserService {
    return &UserService{
        userRepo: userRepo,
        jwtUtil:  jwtUtil,
    }
}
```

**Why no interface?**
- Single implementation (YAGNI principle)
- You can mock the repositories instead
- Simpler, less boilerplate

**When to add service interfaces:**
- Multiple implementations (e.g., different payment providers)
- Want to test handlers in complete isolation from service logic
- Wrapping third-party APIs

---

### 4. External Services: YES ✅

Any external API/service should be behind an interface.

```go
// ✅ GOOD: Interface for external service
type EmailService interface {
    SendWelcomeEmail(to, username string) error
    SendPasswordResetEmail(to, resetToken string) error
}

// Production implementation
type SendGridEmailService struct {
    apiKey string
}

func (s *SendGridEmailService) SendWelcomeEmail(to, username string) error {
    // Actually call SendGrid API
    return sendgrid.Send(/* ... */)
}

// Test/Dev implementation
type MockEmailService struct {
    SentEmails []string
}

func (m *MockEmailService) SendWelcomeEmail(to, username string) error {
    m.SentEmails = append(m.SentEmails, to)
    return nil  // Don't actually send
}
```

**Why interface?**
- External dependency (email, payment, SMS APIs)
- Expensive/slow to call in tests
- May want to swap providers (SendGrid → AWS SES)

---

## Part 2: Two Testing Strategies

### The Key Question: What Layer Are You Testing?

There are **two different approaches** depending on whether you're testing services or handlers.

---

### Strategy 1: Test Services by Mocking Repositories (Recommended ✅)

**When:** Testing service logic
**Mock:** Repositories (you already have repository interfaces!)
**Service Interface Needed:** ❌ NO

```go
// Your current architecture (NO service interface)
type UserHandler struct {
    userService *service.UserService  // Concrete service
}

type UserService struct {
    userRepo repository.UserRepository  // Interface
}

// Test 1: Test the SERVICE
func TestUserService_GetProfile(t *testing.T) {
    // Mock the repository
    mockRepo := &MockUserRepository{
        users: map[uuid.UUID]*entity.User{
            testUserID: {Username: "john", Email: "john@test.com"},
        },
    }
    mockJWT := &MockJWTUtil{}

    // Test REAL service with mocked repository
    service := service.NewUserService(mockRepo, mockJWT)
    result, err := service.GetProfile(ctx, testUserID)

    // You're testing actual service logic!
    assert.NoError(t, err)
    assert.Equal(t, "john", result.Username)
}

// Test 2: Test the HANDLER with real service
func TestUserHandler_GetProfile(t *testing.T) {
    // Use real service with mocked repository
    mockRepo := &MockUserRepository{
        users: map[uuid.UUID]*entity.User{
            testUserID: {Username: "john"},
        },
    }
    mockJWT := &MockJWTUtil{}

    realService := service.NewUserService(mockRepo, mockJWT)  // Real service!
    handler := handler.NewUserHandler(realService)

    // Test handler with real service logic
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Set(constants.ContextUserID, testUserID)

    handler.GetProfile(c)

    // This tests both handler AND service logic together
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**What you're testing:**
- ✅ Real service logic (business rules, validation, etc.)
- ✅ Real handler logic (HTTP handling, status codes, etc.)
- ❌ Not testing database (mocked via repository interface)

**Pros:**
- Tests actual service code
- Only one layer of mocking (repositories)
- Simpler, less code
- Most common in Go

**Cons:**
- Handler tests run service logic (slightly slower, but negligible)

---

### Strategy 2: Test Handlers by Mocking Services (Optional)

**When:** Testing handlers in complete isolation
**Mock:** Services
**Service Interface Needed:** ✅ YES

```go
// Need to add service interface
type UserService interface {
    GetProfile(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
    UpdateProfile(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type UserHandler struct {
    userService UserService  // Interface instead of concrete
}

// Test ONLY the handler (service logic is mocked)
func TestUserHandler_GetProfile(t *testing.T) {
    // Mock the entire service
    mockService := &MockUserService{
        response: &dto.UserResponse{Username: "john"},
        error:    nil,
    }

    handler := handler.NewUserHandler(mockService)

    // Test ONLY HTTP handling
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Set(constants.ContextUserID, testUserID)

    handler.GetProfile(c)

    // Only testing: "Does handler call service? Does it return correct status code?"
    // NOT testing: "Does GetProfile actually fetch from database?"
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**What you're testing:**
- ✅ Handler logic only (HTTP handling)
- ❌ NOT testing service logic (it's mocked with fake data)

**Pros:**
- Ultra-fast handler tests
- Handlers tested in complete isolation
- Good for large teams (handlers and services developed separately)

**Cons:**
- Service logic not tested in handler tests (need separate service tests)
- Requires service interfaces (more code)
- Two layers of mocking

---

### Comparison of Strategies

| Aspect | Strategy 1 (Mock Repos) | Strategy 2 (Mock Services) |
|--------|------------------------|---------------------------|
| **Service Interface** | ❌ Not needed | ✅ Required |
| **What's tested in handler tests** | Handler + Service logic | Handler logic only |
| **Mocking layers** | 1 (repositories) | 2 (repositories + services) |
| **Code complexity** | Lower | Higher |
| **Test speed** | Fast | Slightly faster |
| **Common in Go** | ✅ Very common | Less common |
| **Best for** | Small-medium teams, learning | Large teams, microservices |

---

### Which Strategy for Your Project?

**Recommendation: Strategy 1 (Your Current Approach) ✅**

**Why?**
- You test the actual service logic
- Simpler codebase (no service interfaces)
- Sufficient for most applications
- More idiomatic Go
- Repository interfaces provide the testability boundary you need

**When to switch to Strategy 2:**
- Large team where handlers and services are developed by different people
- You have very complex, slow service logic
- You want to test handlers in complete isolation
- You have multiple service implementations

---

## Part 3: The Go Philosophy

### Define Interfaces at the Point of Use

**Best practice:** Define interfaces where they're **consumed** (used), not where they're **implemented** (defined).

#### Bad (Java-style)
```go
// internal/service/user_service.go (provider defines interface)
type UserService interface {
    GetProfile(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
    UpdateProfile(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) error
    Register(ctx context.Context, req dto.RegisterRequest) error
    Login(ctx context.Context, req dto.LoginRequest) error
    // ... 10 more methods
}

type userServiceImpl struct { /* ... */ }
```

**Problem:** Interface is too large and defined where implemented.

#### Good (Go-style)
```go
// api/handler/user_handler.go (consumer defines minimal interface)
type userProfileService interface {
    GetProfile(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
}

type UserHandler struct {
    userService userProfileService  // Only needs GetProfile
}

// The concrete service has more methods, but handler only uses one
```

**Benefits:**
- Interface is minimal (Interface Segregation Principle)
- Clear dependencies
- Easier to test
- More flexible

---

## Part 4: When to Use Interfaces - Decision Matrix

| Component | Use Interface? | Why / Why Not |
|-----------|----------------|---------------|
| **Entities** | ❌ NO | No dependencies, test directly |
| **DTOs** | ❌ NO | Just data structures |
| **Repositories** | ✅ YES | Database boundary, external dependency |
| **Services (single impl)** | ❌ NO | Mock dependencies instead |
| **Services (multiple impl)** | ✅ YES | Need to swap implementations |
| **External APIs** | ✅ YES | Expensive, need to mock |
| **Utilities** | ⚠️ MAYBE | Only if multiple implementations |

### Specific Scenarios

| Scenario | Use Interface? | Define Where? |
|----------|----------------|---------------|
| Single service implementation | ❌ No | N/A |
| Multiple implementations (payment providers) | ✅ Yes | Domain layer |
| Need to swap at runtime (cache: Redis/Memory) | ✅ Yes | Domain layer |
| Third-party service wrapper (email, SMS) | ✅ Yes | Domain layer |
| Want to test handlers in isolation | ✅ Yes (optional) | Handler package or domain |
| Can mock dependencies instead | ❌ No | Mock the dependencies |

---

## Part 5: Real-World Examples

### Example 1: Email Service (Multiple Providers)
```go
// domain/service/email_service.go
type EmailService interface {
    SendWelcomeEmail(to, username string) error
    SendPasswordResetEmail(to, resetToken string) error
}

// internal/service/sendgrid_email.go
type SendGridEmailService struct {
    apiKey string
}

func (s *SendGridEmailService) SendWelcomeEmail(to, username string) error {
    // Use SendGrid
}

// internal/service/ses_email.go
type SESEmailService struct {
    client *ses.Client
}

func (s *SESEmailService) SendWelcomeEmail(to, username string) error {
    // Use AWS SES
}

// Factory chooses implementation
func NewEmailService(config *Config) EmailService {
    if config.EmailProvider == "sendgrid" {
        return NewSendGridEmailService(config.SendGridKey)
    }
    return NewSESEmailService(config.AWSConfig)
}
```

### Example 2: Cache Service (Multiple Strategies)
```go
type CacheService interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte, ttl time.Duration) error
    Delete(key string) error
}

// Three implementations
type RedisCache struct { /* ... */ }      // Production
type MemoryCache struct { /* ... */ }     // Development
type NoOpCache struct { /* ... */ }       // Caching disabled
```

### Example 3: Payment Service (Multiple Providers)
```go
type PaymentService interface {
    ProcessPayment(amount float64, currency string) (transactionID string, err error)
    RefundPayment(transactionID string) error
}

type StripePaymentService struct { /* ... */ }
type PayPalPaymentService struct { /* ... */ }
type RazorpayPaymentService struct { /* ... */ }
```

---

## Part 6: For Your Todo App

### Your Current Architecture (CORRECT ✅)

```go
// api/handler/user_handler.go
type UserHandler struct {
    userService *service.UserService  // Concrete
}

// internal/service/user_service_impl.go
type UserService struct {
    userRepo repository.UserRepository  // Interface
    jwtUtil  *utils.JWTUtil
}

// domain/repository/user_repository.go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    // ...
}
```

**Why this is correct:**
- ✅ Repository interface provides testability boundary
- ✅ Single service implementation (YAGNI)
- ✅ Test services by mocking repositories
- ✅ Simple, maintainable, idiomatic Go

**Don't add service interfaces unless:**
- You need multiple implementations
- You're wrapping external APIs
- Large team requiring handler/service isolation

---

## Part 7: Testing Patterns

### Pattern 1: Mock Repository (Recommended)

```go
// test/mocks/user_repository_mock.go
type MockUserRepository struct {
    users       map[uuid.UUID]*entity.User
    CreateError error
    FindError   error
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[uuid.UUID]*entity.User),
    }
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
    if m.CreateError != nil {
        return m.CreateError
    }
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    if m.FindError != nil {
        return nil, m.FindError
    }
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("user not found")
    }
    return user, nil
}
```

**Usage:**
```go
func TestUserService_Login_InvalidPassword(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()

    // Pre-populate test data
    user := entity.NewUser("john", "john@test.com", "hashed_pass", "John")
    mockRepo.Create(context.Background(), user)

    mockJWT := mocks.NewMockJWTUtil()
    service := service.NewUserService(mockRepo, mockJWT)

    // Test with wrong password
    _, err := service.Login(ctx, dto.LoginRequest{
        Username: "john",
        Password: "wrong_password",
    })

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid credentials")
}
```

---

## Summary

### What Needs Interfaces?

| Layer | Interface? | Reason |
|-------|-----------|---------|
| Entities | ❌ | No dependencies, test directly |
| Services | ❌ (usually) | Mock repositories instead |
| Repositories | ✅ | Database boundary |
| External APIs | ✅ | Expensive, slow, third-party |

### Testing Strategy

**For your project: Strategy 1 (Mock Repositories)**
- No service interfaces
- Test services with mocked repositories
- Test handlers with real services + mocked repositories
- Simpler, more idiomatic Go

### Key Principles

1. **Accept interfaces, return structs**
2. **Define interfaces at point of use** (not definition)
3. **Keep interfaces small** (1-3 methods)
4. **Don't create interfaces speculatively** (YAGNI)
5. **Interfaces at boundaries** (database, external APIs)

### Your Architecture is Correct! ✅

Your current design already follows Go best practices:
- Repository interfaces for database boundary
- Concrete services with clear dependencies
- Testable without service interfaces
- Simple, maintainable, idiomatic

**Don't over-engineer - your current approach is perfect for this project!**
