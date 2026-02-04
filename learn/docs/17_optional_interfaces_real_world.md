# Optional Interfaces in Go - Real World Usage

## The Problem: Not All Implementations Have Same Capabilities

In real-world applications, different implementations of the same interface may have different **optional capabilities**.

### Real Example: Repository Pattern with Multiple Backends

```
TodoRepository Interface (Required for all implementations)
├── InMemoryRepository     ✅ Has statistics, fast access
├── PostgresRepository     ✅ Has statistics, persistent
├── RedisRepository        ✅ Has statistics, TTL info
└── FileRepository         ❌ No statistics (just files)
```

**Problem:** Some implementations can provide extra features (like statistics), others can't. How do we handle this?

---

## Solution: Optional Interfaces

### Step 1: Define Core Interface (Required)

```go
// domain/repository/todo_repository.go

// TodoRepository - REQUIRED interface (all implementations must have this)
type TodoRepository interface {
    Create(ctx context.Context, todo *entity.Todo) error
    FindByID(ctx context.Context, id string) (*entity.Todo, error)
    FindAll(ctx context.Context) ([]*entity.Todo, error)
    Update(ctx context.Context, todo *entity.Todo) error
    Delete(ctx context.Context, id string) error
}
```

### Step 2: Define Optional Interface (Nice-to-Have)

```go
// StorageInfo - OPTIONAL interface (only some implementations will have this)
type StorageInfo interface {
    GetStorageType() string
    GetStats() map[string]interface{}
}

// CacheCapable - OPTIONAL interface (only cache implementations)
type CacheCapable interface {
    ClearCache() error
    GetCacheHitRate() float64
}

// BatchCapable - OPTIONAL interface (only DB implementations)
type BatchCapable interface {
    BatchCreate(ctx context.Context, todos []*entity.Todo) error
    BatchDelete(ctx context.Context, ids []string) error
}
```

---

## Implementation Examples

### Implementation 1: In-Memory (Has StorageInfo)

```go
package memory

type InMemoryTodoRepository struct {
    mu          sync.RWMutex
    todos       map[string]*entity.Todo
    accessCount int
    lastAccess  time.Time
}

// Required: TodoRepository methods
func (r *InMemoryTodoRepository) Create(ctx context.Context, todo *entity.Todo) error { 
    // implementation
}
func (r *InMemoryTodoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
    // implementation
}
// ... other required methods

// Optional: StorageInfo methods
func (r *InMemoryTodoRepository) GetStorageType() string {
    return "in-memory"
}

func (r *InMemoryTodoRepository) GetStats() map[string]interface{} {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    return map[string]interface{}{
        "storage_type": "in-memory",
        "total_todos":  len(r.todos),
        "access_count": r.accessCount,
        "last_access":  r.lastAccess.Format(time.RFC3339),
    }
}

// Constructor returns only TodoRepository (required interface)
func NewInMemoryTodoRepository() repository.TodoRepository {
    return &InMemoryTodoRepository{
        todos: make(map[string]*entity.Todo),
    }
}

// Compile-time checks (both required and optional)
var _ repository.TodoRepository = (*InMemoryTodoRepository)(nil)  // Required
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)     // Optional (we chose to implement)
```

### Implementation 2: PostgreSQL (Has StorageInfo + BatchCapable)

```go
package postgres

type PostgresTodoRepository struct {
    db *sql.DB
}

// Required: TodoRepository methods
func (r *PostgresTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
    // SQL insert
}
func (r *PostgresTodoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
    // SQL query
}
// ... other required methods

// Optional: StorageInfo
func (r *PostgresTodoRepository) GetStorageType() string {
    return "postgresql"
}

func (r *PostgresTodoRepository) GetStats() map[string]interface{} {
    var count int
    r.db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
    
    return map[string]interface{}{
        "storage_type":  "postgresql",
        "total_todos":   count,
        "database_size": r.getDatabaseSize(),
    }
}

// Optional: BatchCapable (only DB has this)
func (r *PostgresTodoRepository) BatchCreate(ctx context.Context, todos []*entity.Todo) error {
    // Use SQL transaction with batch insert
    tx, _ := r.db.BeginTx(ctx, nil)
    // ... batch insert logic
}

func (r *PostgresTodoRepository) BatchDelete(ctx context.Context, ids []string) error {
    // DELETE FROM todos WHERE id IN (...)
}

// Constructor
func NewPostgresTodRepository(db *sql.DB) repository.TodoRepository {
    return &PostgresTodoRepository{db: db}
}

// Compile-time checks
var _ repository.TodoRepository = (*PostgresTodoRepository)(nil)
var _ repository.StorageInfo = (*PostgresTodoRepository)(nil)
var _ repository.BatchCapable = (*PostgresTodoRepository)(nil)  // Extra capability
```

### Implementation 3: File-Based (NO Optional Interfaces)

```go
package file

type FileTodoRepository struct {
    filePath string
}

// Required: TodoRepository methods (only these!)
func (r *FileTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
    // Write to file
}
func (r *FileTodoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
    // Read from file
}
// ... other required methods

// NO StorageInfo methods (doesn't make sense for files)
// NO BatchCapable methods (files can't do batch operations efficiently)

// Constructor
func NewFileTodoRepository(path string) repository.TodoRepository {
    return &FileTodoRepository{filePath: path}
}

// Only check required interface
var _ repository.TodoRepository = (*FileTodoRepository)(nil)
// No optional interface checks - this implementation doesn't have them!
```

---

## Real-World Usage: How to Handle Switching

### Scenario: Health Check Endpoint

**Use Case:** Admin dashboard that shows repository statistics

```go
// api/handler/health_handler.go
package handler

type HealthHandler struct {
    todoRepo repository.TodoRepository  // We only know it's a TodoRepository
}

func (h *HealthHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    }
    
    // Check if repository supports statistics (OPTIONAL feature)
    if info, ok := h.todoRepo.(repository.StorageInfo); ok {
        // This repository HAS optional StorageInfo interface
        health["storage"] = info.GetStats()
    } else {
        // This repository DOESN'T have StorageInfo
        health["storage"] = "not available"
    }
    
    json.NewEncoder(w).Encode(health)
}
```

**Result:**
```json
// With InMemoryRepository or PostgresRepository (have StorageInfo):
{
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "storage": {
        "storage_type": "in-memory",
        "total_todos": 42,
        "access_count": 1537
    }
}

// With FileRepository (no StorageInfo):
{
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "storage": "not available"
}
```

### Scenario: Bulk Import Feature

**Use Case:** Import 10,000 todos from CSV file

```go
// service/todo_service.go
package service

type TodoService struct {
    repo repository.TodoRepository
}

func (s *TodoService) BulkImport(ctx context.Context, todos []*entity.Todo) error {
    // Check if repository supports batch operations (OPTIONAL feature)
    if batchRepo, ok := s.repo.(repository.BatchCapable); ok {
        // Use efficient batch insert (PostgreSQL, MySQL)
        log.Println("Using batch insert (efficient)")
        return batchRepo.BatchCreate(ctx, todos)
    }
    
    // Fallback: Insert one by one (InMemory, File)
    log.Println("Using individual inserts (slower)")
    for _, todo := range todos {
        if err := s.repo.Create(ctx, todo); err != nil {
            return err
        }
    }
    return nil
}
```

**What happens:**
- PostgreSQL: Uses `BatchCreate()` → Fast (single transaction)
- InMemory/File: Falls back to loop → Slower but works

### Scenario: Cache Management Endpoint

```go
// api/handler/cache_handler.go
package handler

type CacheHandler struct {
    todoRepo repository.TodoRepository
}

func (h *CacheHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
    // Check if repository supports cache operations
    if cacheRepo, ok := h.todoRepo.(repository.CacheCapable); ok {
        err := cacheRepo.ClearCache()
        if err != nil {
            http.Error(w, "Failed to clear cache", 500)
            return
        }
        w.Write([]byte("Cache cleared"))
    } else {
        // This repository doesn't have cache
        http.Error(w, "Cache not supported for this storage", 400)
    }
}
```

---

## Switching Implementations (The Power of Interfaces)

### Configuration-Based Selection

```go
// cmd/api/main.go
package main

func main() {
    config := loadConfig()
    
    // Create repository based on config (EASY SWITCH!)
    var todoRepo repository.TodoRepository
    
    switch config.StorageType {
    case "memory":
        todoRepo = memory.NewInMemoryTodoRepository()
        
    case "postgres":
        db := connectPostgres(config.DatabaseURL)
        todoRepo = postgres.NewPostgresTodRepository(db)
        
    case "file":
        todoRepo = file.NewFileTodoRepository(config.FilePath)
        
    default:
        log.Fatal("Unknown storage type")
    }
    
    // Create service (doesn't care which implementation!)
    todoService := service.NewTodoService(todoRepo)
    
    // Create handlers (doesn't care which implementation!)
    healthHandler := handler.NewHealthHandler(todoRepo)
    
    // Setup routes
    http.HandleFunc("/health", healthHandler.GetHealthStatus)
    http.HandleFunc("/todos", todoService.HandleTodos)
    
    log.Println("Server starting...")
    http.ListenAndServe(":8080", nil)
}
```

### Environment-Based Selection

```bash
# .env for development
STORAGE_TYPE=memory

# .env for production
STORAGE_TYPE=postgres
DATABASE_URL=postgres://user:pass@localhost/todos
```

**Result:** Change one environment variable, entire app switches storage backend!

---

## Benefits of Optional Interfaces

### 1. **Flexibility**
- Core functionality required for all
- Extra features available when implementation supports it
- No breaking changes when adding new capabilities

### 2. **Graceful Degradation**
```go
// Try to use advanced feature, fallback if not available
if batchRepo, ok := repo.(BatchCapable); ok {
    batchRepo.BatchCreate(todos)  // Fast path
} else {
    for _, todo := range todos {
        repo.Create(todo)  // Slow path, but works
    }
}
```

### 3. **Clean Architecture**
```
Service Layer → Only knows TodoRepository interface
                ↓
             Doesn't care if it's:
             - InMemory
             - PostgreSQL  
             - Redis
             - File
             
At runtime → Checks for optional features when needed
```

---

## Type Assertion Pattern (How to Check)

### Pattern 1: Check and Use

```go
if info, ok := repo.(repository.StorageInfo); ok {
    // repo HAS StorageInfo interface
    stats := info.GetStats()
    fmt.Println(stats)
} else {
    // repo DOESN'T have StorageInfo
    fmt.Println("Statistics not available")
}
```

### Pattern 2: Early Check and Store

```go
type TodoService struct {
    repo      repository.TodoRepository
    hasStats  bool
    statsImpl repository.StorageInfo
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
    s := &TodoService{repo: repo}
    
    // Check once at initialization
    if info, ok := repo.(repository.StorageInfo); ok {
        s.hasStats = true
        s.statsImpl = info
    }
    
    return s
}

func (s *TodoService) GetStats() map[string]interface{} {
    if s.hasStats {
        return s.statsImpl.GetStats()  // No type assertion needed
    }
    return map[string]interface{}{"error": "not supported"}
}
```

### Pattern 3: Feature Detection Helper

```go
// Helper function to detect capabilities
func GetRepositoryCapabilities(repo repository.TodoRepository) map[string]bool {
    caps := make(map[string]bool)
    
    _, caps["storage_info"] = repo.(repository.StorageInfo)
    _, caps["batch_capable"] = repo.(repository.BatchCapable)
    _, caps["cache_capable"] = repo.(repository.CacheCapable)
    
    return caps
}

// Usage
caps := GetRepositoryCapabilities(todoRepo)
if caps["batch_capable"] {
    log.Println("This repository supports batch operations")
}
```

---

## Real-World Example: Monitoring Dashboard

```go
// admin/dashboard.go
package admin

type Dashboard struct {
    repos map[string]repository.TodoRepository
}

func (d *Dashboard) GetAllRepositoryStats() map[string]interface{} {
    stats := make(map[string]interface{})
    
    for name, repo := range d.repos {
        repoStats := map[string]interface{}{
            "name": name,
            "required_features": "✓ TodoRepository",
        }
        
        // Check optional features
        if info, ok := repo.(repository.StorageInfo); ok {
            repoStats["storage_info"] = info.GetStats()
            repoStats["storage_type"] = info.GetStorageType()
        }
        
        if batch, ok := repo.(repository.BatchCapable); ok {
            repoStats["supports_batch"] = true
        }
        
        if cache, ok := repo.(repository.CacheCapable); ok {
            repoStats["cache_hit_rate"] = cache.GetCacheHitRate()
        }
        
        stats[name] = repoStats
    }
    
    return stats
}
```

**Output:**
```json
{
    "primary": {
        "name": "primary",
        "required_features": "✓ TodoRepository",
        "storage_info": { "total_todos": 1523 },
        "storage_type": "postgresql",
        "supports_batch": true
    },
    "cache": {
        "name": "cache",
        "required_features": "✓ TodoRepository",
        "storage_info": { "total_todos": 150 },
        "storage_type": "redis",
        "cache_hit_rate": 0.89
    },
    "backup": {
        "name": "backup",
        "required_features": "✓ TodoRepository"
    }
}
```

---

## When to Use Optional Interfaces

| Scenario | Use Optional Interface? | Why |
|----------|------------------------|-----|
| **Statistics/Monitoring** | ✅ Yes (`StorageInfo`) | Not all storage backends have stats |
| **Batch Operations** | ✅ Yes (`BatchCapable`) | Only DBs support efficient batching |
| **Caching** | ✅ Yes (`CacheCapable`) | Only cache implementations need this |
| **Transactions** | ✅ Yes (`Transactional`) | Only DBs have transactions |
| **CRUD Operations** | ❌ No (required) | Every implementation must have these |
| **Basic Queries** | ❌ No (required) | Core functionality for all |

---

## Key Takeaways

1. **Required Interface** = Core functionality all implementations MUST have
   ```go
   type TodoRepository interface { Create, Read, Update, Delete }
   ```

2. **Optional Interface** = Extra capabilities only some implementations have
   ```go
   type StorageInfo interface { GetStats, GetStorageType }
   ```

3. **Type Assertion** = Check at runtime if optional feature exists
   ```go
   if info, ok := repo.(StorageInfo); ok { /* use it */ }
   ```

4. **Easy Switching** = Change implementation without changing service code
   ```go
   var repo TodoRepository = memory.New()  // or postgres.New() or file.New()
   service := NewService(repo)  // Service doesn't care!
   ```

5. **Compile-Time Check** = Optional but recommended for safety
   ```go
   var _ StorageInfo = (*InMemoryRepo)(nil)  // Verify at compile time
   ```

---

## Comparison with Node.js

### Node.js (Duck Typing)
```javascript
// No compile-time checks - runtime only
class TodoService {
    constructor(repo) {
        this.repo = repo;
    }
    
    getStats() {
        if (this.repo.getStats) {  // Check at runtime
            return this.repo.getStats();
        }
        return { error: 'not supported' };
    }
}

// Error only appears when code runs!
const service = new TodoService(new FileRepo());
service.getStats();  // Oops! getStats doesn't exist - runtime error
```

### Go (Static Typing + Interfaces)
```go
// Compile-time checks + runtime type assertions
type TodoService struct {
    repo repository.TodoRepository
}

func (s *TodoService) GetStats() map[string]interface{} {
    if info, ok := s.repo.(repository.StorageInfo); ok {  // Safe check
        return info.GetStats()
    }
    return map[string]interface{}{"error": "not supported"}
}

// No runtime error - gracefully handles missing feature
service := NewTodoService(file.NewFileRepo())
service.GetStats()  // Returns "not supported" - no crash!
```

**Go Advantage:** Type safety + flexibility with optional features!
