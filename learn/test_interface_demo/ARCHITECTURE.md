# Optional Interfaces Architecture

## Visual Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    domain/repository/todo_repository.go                 │
│                                                                         │
│  ┌──────────────────────┐    ┌──────────────────────────────┐         │
│  │ TodoRepository       │    │ OPTIONAL INTERFACES:         │         │
│  │ (REQUIRED)           │    │                               │         │
│  │                      │    │ • StorageInfo                │         │
│  │ • Create()           │    │   - GetStorageType()         │         │
│  │ • FindByID()         │    │   - GetStats()               │         │
│  │ • FindAll()          │    │                               │         │
│  │ • Update()           │    │ • BatchCapable               │         │
│  │ • Delete()           │    │   - BatchCreate()            │         │
│  │                      │    │   - BatchDelete()            │         │
│  └──────────────────────┘    │                               │         │
│                               │ • CacheCapable               │         │
│                               │   - ClearCache()             │         │
│                               │   - GetCacheHitRate()        │         │
│                               │   - GetTTL()                 │         │
│                               └──────────────────────────────┘         │
└─────────────────────────────────────────────────────────────────────────┘
                                     │
                                     │ implements
                    ┌────────────────┼────────────────┐
                    │                │                │
                    ▼                ▼                ▼
    ┌───────────────────────┐  ┌───────────────────────┐  ┌───────────────────────┐
    │  InMemoryRepository   │  │  PostgresRepository   │  │   FileRepository      │
    │  (memory/todo.go)     │  │  (postgres/todo.go)   │  │   (file/todo.go)      │
    ├───────────────────────┤  ├───────────────────────┤  ├───────────────────────┤
    │ ✅ TodoRepository     │  │ ✅ TodoRepository     │  │ ✅ TodoRepository     │
    │ ✅ StorageInfo        │  │ ✅ StorageInfo        │  │ ❌ StorageInfo        │
    │ ❌ BatchCapable       │  │ ✅ BatchCapable       │  │ ❌ BatchCapable       │
    │ ❌ CacheCapable       │  │ ❌ CacheCapable       │  │ ❌ CacheCapable       │
    └───────────────────────┘  └───────────────────────┘  └───────────────────────┘

                    ┌───────────────────────┐
                    │   RedisRepository     │
                    │   (redis/todo.go)     │
                    ├───────────────────────┤
                    │ ✅ TodoRepository     │
                    │ ✅ StorageInfo        │
                    │ ❌ BatchCapable       │
                    │ ✅ CacheCapable       │
                    └───────────────────────┘
```

---

## Type Assertion Flow

```
                           ┌─────────────────┐
                           │  Your Code      │
                           │  (main.go)      │
                           └────────┬────────┘
                                    │
                                    │ Has reference to:
                                    │ var repo TodoRepository
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
                    ▼                               ▼
        ┌───────────────────────┐       ┌───────────────────────┐
        │  Use REQUIRED         │       │  Check for OPTIONAL   │
        │  methods directly     │       │  interfaces at runtime│
        │                       │       │                       │
        │  repo.Create()        │       │  if info, ok :=       │
        │  repo.FindByID()      │       │    repo.(StorageInfo) │
        │  repo.Update()        │       │  {                    │
        │                       │       │    info.GetStats()    │
        │  ✅ Always works      │       │  }                    │
        │                       │       │                       │
        └───────────────────────┘       │  ✅ Safe: checks first│
                                        └───────────────────────┘
```

---

## Compile-Time vs Runtime Checks

### Compile-Time Check (In Implementation File)

```go
// At bottom of memory/todo_memory.go:
var _ repository.TodoRepository = (*InMemoryTodoRepository)(nil)
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)
```

**Purpose:** Verify implementation is correct **before running**

**When it fails:**
```
Error at compile time:
  *InMemoryTodoRepository does not implement repository.StorageInfo
  (missing method GetStats)
```

### Runtime Check (In Usage Code)

```go
// In main.go or service:
if info, ok := repo.(repository.StorageInfo); ok {
    // This code only runs if StorageInfo is implemented
    stats := info.GetStats()
}
```

**Purpose:** Safely use optional features **while running**

**When it fails:**
- Doesn't fail! Just skips the `if` block

---

## Decision Tree: Which Interfaces to Implement?

```
Start: Creating a new repository implementation
│
├─ ALWAYS implement: TodoRepository (required)
│
└─ Can you provide statistics?
   ├─ YES → Implement StorageInfo
   │        • GetStorageType()
   │        • GetStats()
   │
   └─ NO → Don't implement StorageInfo
   
└─ Can you do efficient batch operations?
   ├─ YES (Database) → Implement BatchCapable
   │                   • BatchCreate()
   │                   • BatchDelete()
   │
   └─ NO (File/Memory) → Don't implement BatchCapable
   
└─ Is this a cache?
   ├─ YES (Redis) → Implement CacheCapable
   │                • ClearCache()
   │                • GetCacheHitRate()
   │                • GetTTL()
   │
   └─ NO → Don't implement CacheCapable
```

---

## File Template

### For Each New Implementation:

```go
// internal/repository/YOUR_IMPL/todo_YOUR_IMPL.go
package YOUR_IMPL

import "demo/domain/repository"

// Your concrete type
type YourTodoRepository struct {
    // fields
}

// Constructor
func NewYourTodoRepository() repository.TodoRepository {
    return &YourTodoRepository{}
}

// ============================================================
// REQUIRED INTERFACE - TodoRepository
// ============================================================

func (r *YourTodoRepository) Create(...) { ... }
func (r *YourTodoRepository) FindByID(...) { ... }
func (r *YourTodoRepository) FindAll(...) { ... }
func (r *YourTodoRepository) Update(...) { ... }
func (r *YourTodoRepository) Delete(...) { ... }

// ============================================================
// OPTIONAL INTERFACES (only if you implement them)
// ============================================================

// StorageInfo (if you have it):
// func (r *YourTodoRepository) GetStorageType() string { ... }
// func (r *YourTodoRepository) GetStats() map[string]interface{} { ... }

// BatchCapable (if you have it):
// func (r *YourTodoRepository) BatchCreate(...) { ... }
// func (r *YourTodoRepository) BatchDelete(...) { ... }

// CacheCapable (if you have it):
// func (r *YourTodoRepository) ClearCache() error { ... }
// func (r *YourTodoRepository) GetCacheHitRate() float64 { ... }
// func (r *YourTodoRepository) GetTTL(...) { ... }

// ============================================================
// COMPILE-TIME CHECKS
// ============================================================

// Always check required interface
var _ repository.TodoRepository = (*YourTodoRepository)(nil)

// Only check optional interfaces you actually implemented:
// var _ repository.StorageInfo = (*YourTodoRepository)(nil)
// var _ repository.BatchCapable = (*YourTodoRepository)(nil)
// var _ repository.CacheCapable = (*YourTodoRepository)(nil)
```

---

## Common Patterns in Usage Code

### Pattern 1: Try Advanced, Fall Back to Basic

```go
func BulkImport(repo TodoRepository, todos []*Todo) {
    // Try batch operation (fast)
    if batch, ok := repo.(BatchCapable); ok {
        return batch.BatchCreate(todos)
    }
    
    // Fall back to loop (slower but works)
    for _, todo := range todos {
        repo.Create(todo)
    }
}
```

### Pattern 2: Optional Feature for Monitoring

```go
func HealthCheck(repo TodoRepository) map[string]interface{} {
    health := map[string]interface{}{
        "status": "healthy",
    }
    
    // Add stats if available
    if info, ok := repo.(StorageInfo); ok {
        health["stats"] = info.GetStats()
    }
    
    return health
}
```

### Pattern 3: Cache-Specific Operations

```go
func ClearAllCaches(repos []TodoRepository) {
    for _, repo := range repos {
        // Only clear if it's a cache
        if cache, ok := repo.(CacheCapable); ok {
            cache.ClearCache()
        }
    }
}
```

### Pattern 4: Feature Detection

```go
func GetCapabilities(repo TodoRepository) []string {
    caps := []string{"TodoRepository"}
    
    if _, ok := repo.(StorageInfo); ok {
        caps = append(caps, "StorageInfo")
    }
    if _, ok := repo.(BatchCapable); ok {
        caps = append(caps, "BatchCapable")
    }
    if _, ok := repo.(CacheCapable); ok {
        caps = append(caps, "CacheCapable")
    }
    
    return caps
}
```

---

## Summary Checklist

When creating a new repository implementation:

- [ ] Create file in `internal/repository/YOUR_IMPL/`
- [ ] Implement ALL methods from `TodoRepository` (required)
- [ ] Decide which optional interfaces to implement
- [ ] Add compile-time check for required interface: `var _ repository.TodoRepository = (*YourType)(nil)`
- [ ] Add compile-time checks for optional interfaces you implement
- [ ] Test with runtime type assertions in usage code

When using repositories:

- [ ] Use required methods directly (always available)
- [ ] Check for optional interfaces with type assertion before using
- [ ] Provide fallback behavior if optional interface not available
- [ ] Don't assume any repository has optional features
