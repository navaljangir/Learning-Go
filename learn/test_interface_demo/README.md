# Optional Interfaces Demo - File Organization

This demo shows **where to define and implement optional interfaces** across multiple files.

## Folder Structure

```
demo/
├── domain/
│   └── repository/
│       └── todo_repository.go         ← Define ALL interfaces (required + optional) HERE
│
├── internal/
│   └── repository/
│       ├── memory/
│       │   └── todo_memory.go         ← Implements: TodoRepository + StorageInfo
│       │
│       ├── postgres/
│       │   └── todo_postgres.go       ← Implements: TodoRepository + StorageInfo + BatchCapable
│       │
│       ├── file/
│       │   └── todo_file.go           ← Implements: TodoRepository ONLY
│       │
│       └── redis/
│           └── todo_redis.go          ← Implements: TodoRepository + StorageInfo + CacheCapable
│
└── main.go                            ← Uses all implementations, checks at runtime
```

---

## Step-by-Step Guide

### Step 1: Define ALL Interfaces in ONE Place

**File:** `domain/repository/todo_repository.go`

```go
package repository

// Required interface (ALL implementations MUST have)
type TodoRepository interface {
    Create()
    FindByID()
    // ...
}

// Optional interface #1 (only some implementations)
type StorageInfo interface {
    GetStorageType() string
    GetStats() map[string]interface{}
}

// Optional interface #2 (only some implementations)
type BatchCapable interface {
    BatchCreate()
    BatchDelete()
}

// Optional interface #3 (only some implementations)
type CacheCapable interface {
    ClearCache() error
    GetCacheHitRate() float64
}
```

**Key Point:** All interfaces defined in ONE file for easy reference.

---

### Step 2: Implement in Separate Files

Each implementation goes in its own file with the appropriate compile-time checks.

#### Implementation 1: Memory (Has StorageInfo)

**File:** `internal/repository/memory/todo_memory.go`

```go
package memory

type InMemoryTodoRepository struct { ... }

// Required interface methods
func (r *InMemoryTodoRepository) Create() { ... }
func (r *InMemoryTodoRepository) FindByID() { ... }
// ... other required methods

// Optional interface methods (StorageInfo)
func (r *InMemoryTodoRepository) GetStorageType() string { ... }
func (r *InMemoryTodoRepository) GetStats() map[string]interface{} { ... }

// Compile-time checks
var _ repository.TodoRepository = (*InMemoryTodoRepository)(nil)  // Required
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)     // Optional: We implement this
// No check for BatchCapable or CacheCapable - we don't implement those
```

#### Implementation 2: Postgres (Has StorageInfo + BatchCapable)

**File:** `internal/repository/postgres/todo_postgres.go`

```go
package postgres

type PostgresTodoRepository struct { ... }

// Required interface methods
func (r *PostgresTodoRepository) Create() { ... }
func (r *PostgresTodoRepository) FindByID() { ... }
// ... other required methods

// Optional interface methods (StorageInfo)
func (r *PostgresTodoRepository) GetStorageType() string { ... }
func (r *PostgresTodoRepository) GetStats() map[string]interface{} { ... }

// Optional interface methods (BatchCapable)
func (r *PostgresTodoRepository) BatchCreate() { ... }
func (r *PostgresTodoRepository) BatchDelete() { ... }

// Compile-time checks
var _ repository.TodoRepository = (*PostgresTodoRepository)(nil)  // Required
var _ repository.StorageInfo = (*PostgresTodoRepository)(nil)     // Optional: We implement this
var _ repository.BatchCapable = (*PostgresTodoRepository)(nil)    // Optional: We implement this
// No check for CacheCapable - we don't implement that
```

#### Implementation 3: File (NO Optional Interfaces)

**File:** `internal/repository/file/todo_file.go`

```go
package file

type FileTodoRepository struct { ... }

// Required interface methods ONLY
func (r *FileTodoRepository) Create() { ... }
func (r *FileTodoRepository) FindByID() { ... }
// ... other required methods

// NO optional interface methods

// Compile-time checks
var _ repository.TodoRepository = (*FileTodoRepository)(nil)  // Required ONLY
// No checks for optional interfaces - we don't implement any
```

#### Implementation 4: Redis (Has StorageInfo + CacheCapable)

**File:** `internal/repository/redis/todo_redis.go`

```go
package redis

type RedisTodoRepository struct { ... }

// Required interface methods
func (r *RedisTodoRepository) Create() { ... }
func (r *RedisTodoRepository) FindByID() { ... }
// ... other required methods

// Optional interface methods (StorageInfo)
func (r *RedisTodoRepository) GetStorageType() string { ... }
func (r *RedisTodoRepository) GetStats() map[string]interface{} { ... }

// Optional interface methods (CacheCapable)
func (r *RedisTodoRepository) ClearCache() error { ... }
func (r *RedisTodoRepository) GetCacheHitRate() float64 { ... }
func (r *RedisTodoRepository) GetTTL() { ... }

// Compile-time checks
var _ repository.TodoRepository = (*RedisTodoRepository)(nil)  // Required
var _ repository.StorageInfo = (*RedisTodoRepository)(nil)     // Optional: We implement this
var _ repository.CacheCapable = (*RedisTodoRepository)(nil)    // Optional: We implement this
// No check for BatchCapable - we don't implement that
```

---

### Step 3: Use and Check at Runtime

**File:** `main.go` or any service file

```go
package main

func main() {
    // Create any implementation
    var repo repository.TodoRepository
    repo = memory.NewInMemoryTodoRepository()
    // OR: repo = postgres.NewPostgresTodRepository(db)
    // OR: repo = file.NewFileTodoRepository(path)
    // OR: repo = redis.NewRedisTodoRepository(client)
    
    // Use required methods (always available)
    repo.Create(todo)
    repo.FindByID(id)
    
    // Check for optional StorageInfo interface
    if info, ok := repo.(repository.StorageInfo); ok {
        stats := info.GetStats()  // Available!
    } else {
        // Not available, handle gracefully
    }
    
    // Check for optional BatchCapable interface
    if batch, ok := repo.(repository.BatchCapable); ok {
        batch.BatchCreate(todos)  // Available!
    } else {
        // Fall back to loop
        for _, todo := range todos {
            repo.Create(todo)
        }
    }
    
    // Check for optional CacheCapable interface
    if cache, ok := repo.(repository.CacheCapable); ok {
        cache.ClearCache()  // Available!
    } else {
        // Not a cache, nothing to clear
    }
}
```

---

## Implementation Matrix

| Implementation | File Location | Required: TodoRepository | Optional: StorageInfo | Optional: BatchCapable | Optional: CacheCapable |
|---------------|---------------|-------------------------|----------------------|----------------------|----------------------|
| **Memory** | `memory/todo_memory.go` | ✅ YES | ✅ YES | ❌ NO | ❌ NO |
| **Postgres** | `postgres/todo_postgres.go` | ✅ YES | ✅ YES | ✅ YES | ❌ NO |
| **File** | `file/todo_file.go` | ✅ YES | ❌ NO | ❌ NO | ❌ NO |
| **Redis** | `redis/todo_redis.go` | ✅ YES | ✅ YES | ❌ NO | ✅ YES |

---

## Compile-Time Check Pattern

### For Each Implementation File:

```go
// At the bottom of the file:

// ALWAYS check the required interface
var _ repository.TodoRepository = (*YourType)(nil)

// ONLY check optional interfaces you implement
var _ repository.StorageInfo = (*YourType)(nil)    // If you implement GetStorageType() and GetStats()
var _ repository.BatchCapable = (*YourType)(nil)   // If you implement BatchCreate() and BatchDelete()
var _ repository.CacheCapable = (*YourType)(nil)   // If you implement ClearCache(), GetCacheHitRate(), etc.
```

### What Happens:

1. **If all methods exist:** Code compiles ✅
2. **If a method is missing:** Compile error ❌

```
Error: *YourType does not implement repository.StorageInfo
       (missing method GetStats)
```

---

## Key Rules

1. **Define all interfaces in ONE place:** `domain/repository/todo_repository.go`
2. **Each implementation in separate file:** `memory/`, `postgres/`, `file/`, `redis/`
3. **Compile-time checks at bottom of each implementation file:**
   ```go
   var _ RequiredInterface = (*YourType)(nil)     // Always
   var _ OptionalInterface = (*YourType)(nil)     // Only if you implement it
   ```
4. **Runtime type assertion when using optional features:**
   ```go
   if optImpl, ok := repo.(OptionalInterface); ok {
       // Use it
   }
   ```

---

## Running the Demo

```bash
cd learn/test_interface_demo

# Fix import paths first (change "demo/" to actual module path)

# Run the demo
go run main.go
```

**Expected Output:**
```
=== Testing Different Implementations ===

--- Testing Memory Repository ---
✓ Created todo
✓ Found todo: Test Todo from Memory
✓ Storage type: in-memory
✓ Stats: map[access_count:2 ...]

--- Testing File Repository ---
✓ Created todo
✓ Found todo: Test Todo from File
✗ StorageInfo not supported

--- Testing Redis Repository ---
✓ Created todo
✓ Found todo: Test Todo from Redis
✓ Storage type: redis-cache
✓ Stats: map[cache_hits:1 ...]

=== Bulk Import Demo ===
✗ BatchCapable not supported, using loop (slower)
✗ BatchCapable not supported, using loop (slower)

=== Cache Management Demo ===
✗ CacheCapable not supported
✓ Cache operations available
  - Hit rate: 100.00%
  - Cache cleared
```

---

## Summary

- **ONE file** for all interface definitions
- **Multiple files** for different implementations
- **Compile-time checks** (`var _ = ()(nil)`) at the bottom of each implementation
- **Runtime type assertions** (`if x, ok := repo.(Interface); ok`) when using optional features
- **Easy switching** between implementations without changing business logic
