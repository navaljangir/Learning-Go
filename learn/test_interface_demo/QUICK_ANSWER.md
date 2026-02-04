# Quick Answer: Where to Define and Implement Optional Interfaces

## Your Question:
> "Where to implement different optional fields in different files and then use with the `var _ = ()(nil)` something like this?"

---

## The Answer:

### 1. Define ALL interfaces in ONE place

**File:** `domain/repository/todo_repository.go`

```go
package repository

// Required interface
type TodoRepository interface {
    Create()
    FindByID()
    Update()
    Delete()
}

// Optional interface #1
type StorageInfo interface {
    GetStorageType() string
    GetStats() map[string]interface{}
}

// Optional interface #2
type BatchCapable interface {
    BatchCreate()
    BatchDelete()
}
```

---

### 2. Implement in DIFFERENT files (one file per implementation)

#### File 1: `internal/repository/memory/todo_memory.go`

```go
package memory

type InMemoryTodoRepository struct { ... }

// Required methods
func (r *InMemoryTodoRepository) Create() { ... }
func (r *InMemoryTodoRepository) FindByID() { ... }
// etc...

// Optional methods (we chose to implement StorageInfo)
func (r *InMemoryTodoRepository) GetStorageType() string { 
    return "memory" 
}
func (r *InMemoryTodoRepository) GetStats() map[string]interface{} { 
    return map[string]interface{}{ "total": len(r.todos) }
}

// ⭐ Compile-time checks at the BOTTOM of the file:
var _ repository.TodoRepository = (*InMemoryTodoRepository)(nil)  // Check required
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)     // Check optional (we implement it)
// No check for BatchCapable because we DON'T implement it
```

#### File 2: `internal/repository/postgres/todo_postgres.go`

```go
package postgres

type PostgresTodoRepository struct { ... }

// Required methods
func (r *PostgresTodoRepository) Create() { ... }
func (r *PostgresTodoRepository) FindByID() { ... }
// etc...

// Optional methods (we implement BOTH StorageInfo AND BatchCapable)
func (r *PostgresTodoRepository) GetStorageType() string { 
    return "postgres" 
}
func (r *PostgresTodoRepository) GetStats() map[string]interface{} { 
    // Query database for stats
}
func (r *PostgresTodoRepository) BatchCreate() { 
    // Use SQL transaction for efficient batch insert
}
func (r *PostgresTodoRepository) BatchDelete() { 
    // Use SQL DELETE WHERE IN (...)
}

// ⭐ Compile-time checks at the BOTTOM of the file:
var _ repository.TodoRepository = (*PostgresTodoRepository)(nil)  // Check required
var _ repository.StorageInfo = (*PostgresTodoRepository)(nil)     // Check optional #1
var _ repository.BatchCapable = (*PostgresTodoRepository)(nil)    // Check optional #2
```

#### File 3: `internal/repository/file/todo_file.go`

```go
package file

type FileTodoRepository struct { ... }

// Required methods ONLY
func (r *FileTodoRepository) Create() { ... }
func (r *FileTodoRepository) FindByID() { ... }
// etc...

// NO optional methods (file storage doesn't support them)

// ⭐ Compile-time check at the BOTTOM of the file:
var _ repository.TodoRepository = (*FileTodoRepository)(nil)  // Check required ONLY
// NO checks for optional interfaces - we don't implement any!
```

---

### 3. Use with runtime type assertions

**File:** `main.go` or any service

```go
package main

func main() {
    // Get any implementation
    var repo repository.TodoRepository
    repo = memory.NewInMemoryTodoRepository()  // Or postgres, or file
    
    // Use required methods (always works)
    repo.Create(todo)
    
    // Check for optional StorageInfo
    if info, ok := repo.(repository.StorageInfo); ok {
        stats := info.GetStats()  // Only runs if implemented
        fmt.Println(stats)
    } else {
        fmt.Println("Stats not available")  // Graceful fallback
    }
    
    // Check for optional BatchCapable
    if batch, ok := repo.(repository.BatchCapable); ok {
        batch.BatchCreate(todos)  // Fast batch insert
    } else {
        for _, t := range todos {
            repo.Create(t)  // Slower but works
        }
    }
}
```

---

## The Pattern: `var _ Interface = (*Type)(nil)`

### What it does:
```go
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)
│   │  │                       │
│   │  │                       └─ Create nil pointer of your type
│   │  └─ Interface to check against
│   └─ Blank identifier (throw away the variable)
└─ Declare variable
```

### Why use it:
1. **Compile-time safety** - Catches missing methods BEFORE running
2. **Documentation** - Shows which interfaces your type implements
3. **Optional but recommended** - Especially for optional interfaces

### When it fails:
```
Compile error if method missing:
  *InMemoryTodoRepository does not implement repository.StorageInfo
  (missing method GetStats)
```

### When to use it:
- ✅ At the **bottom** of each implementation file
- ✅ For the **required** interface (always)
- ✅ For **optional** interfaces you chose to implement
- ❌ NOT for optional interfaces you didn't implement

---

## Summary Table

| File Location | Defines What | Uses `var _ = ()(nil)` |
|--------------|--------------|------------------------|
| `domain/repository/todo_repository.go` | ALL interfaces (required + optional) | ❌ No (interfaces don't need it) |
| `internal/repository/memory/todo_memory.go` | Implements TodoRepository + StorageInfo | ✅ Yes (2 checks) |
| `internal/repository/postgres/todo_postgres.go` | Implements TodoRepository + StorageInfo + BatchCapable | ✅ Yes (3 checks) |
| `internal/repository/file/todo_file.go` | Implements TodoRepository only | ✅ Yes (1 check) |
| `main.go` | Uses repositories with type assertions | ❌ No (this is usage code) |

---

## Real Example from Your Code

In your `08_todo_concurrency/internal/repository/memory/todo_memory.go`:

```go
// Line 180 (currently commented out):
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)
```

**What it's doing:**
- Checking if `*InMemoryTodoRepository` implements `StorageInfo` interface
- This is a **compile-time safety check**
- It's **optional** (that's why it can be commented out)
- Should be uncommented when `GetStats()` is implemented

**Why it's commented:**
- The `GetStats()` method (lines 166-176) is also commented out
- If you uncomment the check but not the method → compile error
- Both should be uncommented together

---

## The Key Insight

**One interface definition → Multiple different implementations → Each checks what it implements**

```
                    domain/repository/todo_repository.go
                    ┌────────────────────────────────┐
                    │ Define: TodoRepository         │
                    │ Define: StorageInfo (optional) │
                    │ Define: BatchCapable (optional)│
                    └────────────────┬───────────────┘
                                     │
                        ┌────────────┼────────────┐
                        │            │            │
                        ▼            ▼            ▼
         memory/todo_memory.go  postgres/todo_postgres.go  file/todo_file.go
         ┌──────────────────┐  ┌───────────────────────┐  ┌─────────────────┐
         │ Implement:       │  │ Implement:            │  │ Implement:      │
         │ • TodoRepository │  │ • TodoRepository      │  │ • TodoRepository│
         │ • StorageInfo    │  │ • StorageInfo         │  │                 │
         │                  │  │ • BatchCapable        │  │ (nothing else)  │
         │                  │  │                       │  │                 │
         │ Check with:      │  │ Check with:           │  │ Check with:     │
         │ var _ = ...      │  │ var _ = ...           │  │ var _ = ...     │
         │ var _ = ...      │  │ var _ = ...           │  │ (only 1 check)  │
         │ (2 checks)       │  │ var _ = ...           │  │                 │
         │                  │  │ (3 checks)            │  │                 │
         └──────────────────┘  └───────────────────────┘  └─────────────────┘
```

Each file checks ONLY what IT implements!
