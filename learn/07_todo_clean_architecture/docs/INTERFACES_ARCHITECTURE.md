# Interfaces Architecture

## Overview

This document explains how interfaces are used for dependency injection and testing in our clean architecture TODO app.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────────────┐
│                    API / Presentation Layer                      │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐                │
│  │  Auth      │  │   User     │  │   Todo     │                │
│  │  Handler   │  │  Handler   │  │  Handler   │                │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘                │
│        │                │                │                        │
│        └────────────────┼────────────────┘                        │
│                         │                                         │
└─────────────────────────┼─────────────────────────────────────────┘
                          │ Depends on
                          │ (via interfaces)
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Domain Layer                                │
│                   (Interfaces Only)                              │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  domain/service/user_service.go                         │    │
│  │                                                          │    │
│  │  type UserService interface {                           │    │
│  │      Register(...)                                      │    │
│  │      Login(...)                                         │    │
│  │      GetProfile(...)                                    │    │
│  │      UpdateProfile(...)                                 │    │
│  │  }                                                       │    │
│  └────────────────────────────────────────────────────────┘    │
│                                                                   │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  domain/service/todo_service.go                         │    │
│  │                                                          │    │
│  │  type TodoService interface {                           │    │
│  │      Create(...)                                        │    │
│  │      GetByID(...)                                       │    │
│  │      List(...)                                          │    │
│  │      Update(...)                                        │    │
│  │      ToggleComplete(...)                                │    │
│  │      Delete(...)                                        │    │
│  │  }                                                       │    │
│  └────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                          ▲
                          │ Implements
                          │ (concrete types)
                          │
┌─────────────────────────┼─────────────────────────────────────────┐
│                   Application Layer                               │
│                (Concrete Implementations)                         │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  internal/service/user_service_impl.go                  │    │
│  │                                                          │    │
│  │  type UserServiceImpl struct {                          │    │
│  │      userRepo repository.UserRepository                 │    │
│  │      jwtUtil  *utils.JWTUtil                            │    │
│  │  }                                                       │    │
│  │                                                          │    │
│  │  var _ domainService.UserService = (*UserServiceImpl)(nil)   │
│  │  ↑ Compile-time check                                  │    │
│  └────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

## Example: Todo Creation Flow

Let's trace how a TODO creation request flows through the system:

### 1. **Handler receives HTTP request** (Presentation Layer)

```go
// api/handler/todo_handler.go

type TodoHandler struct {
    todoService service.TodoService  // ← Interface (not concrete type!)
}

func (h *TodoHandler) Create(c *gin.Context) {
    userID := c.MustGet(constants.ContextUserID).(uuid.UUID)
    var req dto.CreateTodoRequest
    c.ShouldBindJSON(&req)

    // Call interface method
    response, err := h.todoService.Create(c.Request.Context(), userID, req)

    utils.Created(c, response)
}
```

**Key point:** Handler doesn't know if `todoService` is a real implementation, a mock, or a test double. It only knows the interface contract.

### 2. **Interface defines the contract** (Domain Layer)

```go
// domain/service/todo_service.go

type TodoService interface {
    Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error)
    GetByID(ctx context.Context, todoID, userID uuid.UUID) (*dto.TodoResponse, error)
    // ... other methods
}
```

**Key point:** This is the contract that any implementation must satisfy.

### 3. **Concrete implementation** (Application Layer)

```go
// internal/service/todo_service_impl.go

type TodoServiceImpl struct {
    todoRepo repository.TodoRepository
}

// Compile-time verification that TodoServiceImpl implements TodoService
var _ domainService.TodoService = (*TodoServiceImpl)(nil)

func (s *TodoServiceImpl) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
    priority := entity.Priority(req.Priority)
    todo := entity.NewTodo(userID, req.Title, req.Description, priority, req.DueDate)

    if err := s.todoRepo.Create(ctx, todo); err != nil {
        return nil, err
    }

    response := dto.TodoToResponse(todo)
    return &response, nil
}
```

**Key point:** The concrete implementation satisfies the interface contract. The compile-time check ensures we haven't missed any methods.

### 4. **Dependency Injection** (Main)

```go
// cmd/api/main.go

// Create concrete implementation
func initServices(
    userRepo repository.UserRepository,
    todoRepo repository.TodoRepository,
    jwtUtil *utils.JWTUtil,
) (service.UserService, service.TodoService) {  // ← Return interfaces!
    userService := serviceImpl.NewUserService(userRepo, jwtUtil)
    todoService := serviceImpl.NewTodoService(todoRepo)
    return userService, todoService
}

// Inject into handlers
func initHandlers(
    userService service.UserService,  // ← Accept interfaces!
    todoService service.TodoService,
) (*handler.AuthHandler, *handler.UserHandler, *handler.TodoHandler) {
    authHandler := handler.NewAuthHandler(userService)
    userHandler := handler.NewUserHandler(userService)
    todoHandler := handler.NewTodoHandler(todoService)
    return authHandler, userHandler, todoHandler
}
```

**Key point:** We create concrete implementations but pass them as interfaces. This allows swapping implementations without changing the handlers.

## Benefits of This Pattern

### 1. **Testability**
You can easily create mock implementations for testing:

```go
// In tests
type MockTodoService struct {
    CreateFunc func(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error)
}

func (m *MockTodoService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
    return m.CreateFunc(ctx, userID, req)
}

// ... implement other interface methods

// Use in test
func TestTodoHandler_Create(t *testing.T) {
    mockService := &MockTodoService{
        CreateFunc: func(...) (*dto.TodoResponse, error) {
            return &dto.TodoResponse{ID: "123", Title: "Test"}, nil
        },
    }

    handler := handler.NewTodoHandler(mockService)  // ← Inject mock!
    // ... test handler
}
```

### 2. **Loose Coupling**
Handlers depend on abstractions (interfaces), not concrete implementations. This means:
- You can change the implementation without touching handlers
- You can swap MySQL for PostgreSQL without changing handlers
- You can add caching, logging, or other features by creating wrapper implementations

### 3. **Clean Architecture**
The dependency flow is correct:
- **Presentation → Domain (interfaces)**
- **Application → Domain (implements interfaces)**
- Domain layer has no dependencies (pure business logic)

## Common Mistakes

### [BAD] Pointer to Interface
```go
type TodoHandler struct {
    todoService *service.TodoService  // WRONG!
}
```

This is a "pointer to an interface" which is almost never what you want. Interfaces are already reference types internally.

### [BAD] Using Concrete Types in Handlers
```go
import "todo_app/internal/service"

type TodoHandler struct {
    todoService *service.TodoServiceImpl  // WRONG! Tight coupling
}
```

This defeats the purpose of interfaces - now you can't swap implementations.

### ** Correct Pattern
```go
import "todo_app/domain/service"

type TodoHandler struct {
    todoService service.TodoService  // CORRECT! Depend on abstraction
}
```

## Testing Example

Here's how you'd write a unit test for the handler:

```go
package handler_test

import (
    "testing"
    "todo_app/api/handler"
    "todo_app/domain/service"
    "todo_app/internal/dto"

    "github.com/google/uuid"
)

// Mock implementation of TodoService
type MockTodoService struct {
    CreateFn func(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error)
}

func (m *MockTodoService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
    return m.CreateFn(ctx, userID, req)
}

// Implement other methods to satisfy interface...
func (m *MockTodoService) GetByID(...) { return nil, nil }
func (m *MockTodoService) List(...) { return nil, nil }
func (m *MockTodoService) Update(...) { return nil, nil }
func (m *MockTodoService) ToggleComplete(...) { return nil, nil }
func (m *MockTodoService) Delete(...) { return nil }

func TestTodoHandler_Create(t *testing.T) {
    // Arrange
    mockService := &MockTodoService{
        CreateFn: func(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
            return &dto.TodoResponse{
                ID:    uuid.New(),
                Title: req.Title,
            }, nil
        },
    }

    handler := handler.NewTodoHandler(mockService)

    // Act & Assert
    // ... test the handler
}
```

## Summary

1. **Interfaces define contracts** in the domain layer
2. **Concrete types implement** these interfaces in the application layer
3. **Handlers depend on interfaces**, not concrete types
4. **Main.go wires everything together** using dependency injection
5. **Tests can inject mocks** instead of real implementations
6. **Never use pointers to interfaces** - interfaces are already reference types
