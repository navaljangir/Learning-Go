# Handler Interfaces in Go Clean Architecture

## Overview

Handler interfaces define the contract for HTTP request handlers, allowing for better testability, flexibility, and adherence to clean architecture principles.

## What Are Handler Interfaces?

In our TODO application, handlers are the **Presentation Layer** that process HTTP requests. By defining interfaces for handlers, we can:

1. **Test more easily** - Mock handlers in tests
2. **Follow SOLID principles** - Depend on abstractions, not implementations
3. **Enable flexibility** - Swap implementations without changing dependent code
4. **Maintain consistency** - Apply the same pattern across all architectural layers

## File Structure

```
api/
└── handler/
    ├── interfaces.go        # Handler interface definitions
    ├── auth_handler.go      # Concrete AuthHandler implementation
    ├── user_handler.go      # Concrete UserHandler implementation
    └── todo_handler.go      # Concrete TodoHandler implementation
```

## Interface Definitions

### api/handler/interfaces.go

```go
package handler

import "github.com/gin-gonic/gin"

// AuthHandlerInterface defines methods for authentication handlers
type AuthHandlerInterface interface {
    Register(c *gin.Context)
    Login(c *gin.Context)
}

// UserHandlerInterface defines methods for user profile handlers
type UserHandlerInterface interface {
    GetProfile(c *gin.Context)
    UpdateProfile(c *gin.Context)
}

// TodoHandlerInterface defines methods for todo handlers
type TodoHandlerInterface interface {
    Create(c *gin.Context)
    List(c *gin.Context)
    GetByID(c *gin.Context)
    Update(c *gin.Context)
    ToggleComplete(c *gin.Context)
    Delete(c *gin.Context)
}
```

## Implementation Pattern

### Concrete Handler (Implements Interface)

```go
// TodoHandler is the concrete implementation
type TodoHandler struct {
    todoService service.TodoService
}

// NewTodoHandler creates a new handler instance
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
    return &TodoHandler{todoService: todoService}
}

// Create implements TodoHandlerInterface.Create
func (h *TodoHandler) Create(c *gin.Context) {
    // Handler logic...
}
```

**Key Points:**
- The concrete struct (`TodoHandler`) automatically implements the interface if it has all required methods
- No explicit "implements" keyword needed in Go
- Constructor returns concrete type (`*TodoHandler`), which is then used as interface type

## Usage in Router

### Before (Using Concrete Types)

```go
func SetupRouter(
    authHandler *handler.AuthHandler,
    userHandler *handler.UserHandler,
    todoHandler *handler.TodoHandler,
    jwtUtil *utils.JWTUtil,
) *gin.Engine {
    // ...
}
```

### After (Using Interfaces)

```go
func SetupRouter(
    authHandler handler.AuthHandlerInterface,
    userHandler handler.UserHandlerInterface,
    todoHandler handler.TodoHandlerInterface,
    jwtUtil *utils.JWTUtil,
) *gin.Engine {
    r := gin.New()

    // Use handlers via interface methods
    auth := v1.Group("/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }

    return r
}
```

## Usage in main.go

```go
// initHandlers returns interfaces, not concrete types
func initHandlers(
    userService service.UserService,
    todoService service.TodoService,
) (handler.AuthHandlerInterface, handler.UserHandlerInterface, handler.TodoHandlerInterface) {
    // Create concrete implementations
    authHandler := handler.NewAuthHandler(userService)
    userHandler := handler.NewUserHandler(userService)
    todoHandler := handler.NewTodoHandler(todoService)

    // Return as interface types (automatic conversion)
    return authHandler, userHandler, todoHandler
}
```

## Benefits

### 1. Easier Testing

You can create mock handlers for testing router logic:

```go
type MockTodoHandler struct{}

func (m *MockTodoHandler) Create(c *gin.Context) {
    c.JSON(200, gin.H{"message": "mock create"})
}

func (m *MockTodoHandler) List(c *gin.Context) { /* ... */ }
// Implement all interface methods...

// In test:
func TestRouter(t *testing.T) {
    mockHandler := &MockTodoHandler{}
    router := SetupRouter(nil, nil, mockHandler, nil)
    // Test router without real dependencies
}
```

### 2. Dependency Inversion Principle

```
High-level module (Router) → Depends on → Interface (TodoHandlerInterface)
                                              ↑
                                              |
                            Implements by TodoHandler (Low-level module)
```

The router depends on the interface, not the concrete implementation. This is the **Dependency Inversion Principle** from SOLID.

### 3. Flexibility

You can swap implementations without changing router code:

```go
// Development handler with extra logging
type DevTodoHandler struct {
    *TodoHandler
}

func (d *DevTodoHandler) Create(c *gin.Context) {
    log.Println("DEV: Creating todo")
    d.TodoHandler.Create(c)
}

// Still satisfies TodoHandlerInterface!
func main() {
    devHandler := &DevTodoHandler{TodoHandler: handler.NewTodoHandler(svc)}
    router := SetupRouter(nil, nil, devHandler, nil)
}
```

### 4. Consistency Across Layers

Our architecture now uses interfaces consistently:

```
Presentation Layer: handler.TodoHandlerInterface
         ↓
Application Layer:  service.TodoService (interface)
         ↓
Domain Layer:       repository.TodoRepository (interface)
         ↓
Infrastructure:     sqlc_impl.TodoRepository (concrete)
```

## Go Interface Implementation Rules

### Implicit Implementation

Go interfaces are satisfied **implicitly**:

```go
// If TodoHandler has all methods from TodoHandlerInterface,
// it automatically implements the interface - no declaration needed!

var handler TodoHandlerInterface = NewTodoHandler(svc) // ✅ Works automatically
```

### Empty Interface

```go
interface{} // Accepts any type (like TypeScript's 'any')
```

### Type Assertion

```go
// If you need the concrete type:
concreteHandler := handler.(*TodoHandler)  // Type assertion
```

## When NOT to Use Handler Interfaces

You might skip handler interfaces if:

1. **Very simple application** - Single handler, no testing needed
2. **Prototyping** - Quick proof of concept
3. **No testing requirements** - If you're not writing tests (not recommended!)

However, even for simple apps, using interfaces is a good practice and doesn't add much overhead.

## Comparison with Other Layers

| Layer | Interface Location | Implementation |
|-------|-------------------|----------------|
| **Handler** | `api/handler/interfaces.go` | `api/handler/todo_handler.go` |
| **Service** | `domain/service/todo_service.go` | `internal/service/todo_service.go` |
| **Repository** | `domain/repository/todo_repository.go` | `internal/repository/sqlc_impl/todo_repository.go` |

All follow the same pattern: **interface in domain/public layer, implementation in internal/infrastructure layer**.

## Node.js Analogy

### TypeScript Interface (Similar concept)

```typescript
// TypeScript
interface TodoHandler {
  create(req: Request, res: Response): Promise<void>;
  list(req: Request, res: Response): Promise<void>;
}

class TodoHandlerImpl implements TodoHandler {
  async create(req: Request, res: Response): Promise<void> {
    // Implementation
  }

  async list(req: Request, res: Response): Promise<void> {
    // Implementation
  }
}

// Use interface type
function setupRouter(todoHandler: TodoHandler) {
  router.post('/todos', todoHandler.create);
  router.get('/todos', todoHandler.list);
}
```

### Go (This application)

```go
// Go - Interface
type TodoHandlerInterface interface {
    Create(c *gin.Context)
    List(c *gin.Context)
}

// Go - Implementation (implicit, no "implements" keyword)
type TodoHandler struct {
    todoService service.TodoService
}

func (h *TodoHandler) Create(c *gin.Context) {
    // Implementation
}

func (h *TodoHandler) List(c *gin.Context) {
    // Implementation
}

// Use interface type
func SetupRouter(todoHandler TodoHandlerInterface) *gin.Engine {
    r := gin.New()
    r.POST("/todos", todoHandler.Create)
    r.GET("/todos", todoHandler.List)
    return r
}
```

**Key Difference:** Go uses **implicit** implementation (no `implements` keyword), TypeScript uses **explicit** implementation.

## Summary

Handler interfaces provide:
- ✅ Better testability through mocking
- ✅ Dependency inversion (depend on abstractions)
- ✅ Flexibility to swap implementations
- ✅ Consistency with other architectural layers
- ✅ Cleaner, more maintainable code

They're a small addition that brings significant architectural benefits, especially as your application grows.
