# Todo App - Learning Concurrency & Interfaces

An educational TODO application demonstrating Go's core features: **Interfaces**, **Goroutines**, **Channels**, and **Mutex**.

## Learning Objectives

### 1. **Interfaces**
Learn how Go's interfaces enable polymorphism and dependency injection.
- Multiple storage implementations (in-memory, cache)
- Interface-based architecture
- Strategy pattern in Go

### 2. **Goroutines**
Understand lightweight concurrent execution.
- Background task processing
- Async notifications
- Batch operations

### 3. **Channels**
Master communication between goroutines.
- Job queues
- Result aggregation
- Pipeline patterns

### 4. **Mutex**
Learn thread-safe shared state management.
- Concurrent counters
- Statistics tracking
- Cache synchronization

---

## Project Structure

```
08_todo_concurrency/
├── cmd/api/main.go              # Entry point
├── domain/
│   ├── entity/todo.go           # Todo entity
│   └── repository/
│       └── todo_repository.go   # Repository INTERFACE (key learning!)
├── internal/
│   ├── repository/
│   │   ├── memory/              # In-memory implementation
│   │   │   └── todo_memory.go
│   │   └── cache/               # Cached implementation
│   │       └── todo_cache.go
│   ├── service/
│   │   ├── todo_service.go      # Uses interface
│   │   ├── batch_processor.go   # Goroutines + Channels
│   │   ├── notifier.go          # Async notifications
│   │   └── stats.go             # Mutex for counters
│   └── dto/todo_dto.go
├── api/
│   ├── handler/
│   │   ├── todo_handler.go      # Standard CRUD
│   │   ├── batch_handler.go     # Batch operations
│   │   ├── stats_handler.go     # Statistics
│   │   └── notify_handler.go    # Notifications
│   └── router/router.go
└── pkg/utils/response.go
```

---

## API Endpoints

### Basic CRUD (Demonstrates Interfaces)
```
POST   /api/v1/todos           - Create todo
GET    /api/v1/todos           - List todos
GET    /api/v1/todos/:id       - Get specific todo
PUT    /api/v1/todos/:id       - Update todo
DELETE /api/v1/todos/:id       - Delete todo
```

### Concurrency Features

#### 1. Batch Processing (Goroutines + Channels)
```
POST /api/v1/todos/batch
```
**Learns:** Fan-out pattern, worker pools, result aggregation
**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/todos/batch \
  -H "Content-Type: application/json" \
  -d '{
    "todos": [
      {"title": "Task 1", "description": "First task"},
      {"title": "Task 2", "description": "Second task"},
      {"title": "Task 3", "description": "Third task"}
    ]
  }'
```

#### 2. Async Notifications (Goroutines + Channels)
```
POST /api/v1/todos/:id/notify
```
**Learns:** Non-blocking operations, buffered channels
**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/todos/1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Don't forget this task!",
    "delay_seconds": 5
  }'
```

#### 3. Statistics (Mutex)
```
GET /api/v1/stats
```
**Learns:** Thread-safe counters, read/write locks
**Returns:**
```json
{
  "total_todos": 150,
  "completed": 75,
  "pending": 75,
  "requests_processed": 1234
}
```

#### 4. Cache Statistics (Interface + Mutex)
```
GET /api/v1/cache/stats
```
**Learns:** Multiple interface implementations
**Example:**
```bash
curl http://localhost:8080/api/v1/cache/stats
```

#### 5. Switch Storage Backend (Interface Switching)
```
POST /api/v1/admin/switch-storage
```
**Learns:** Runtime interface switching
**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/switch-storage \
  -H "Content-Type: application/json" \
  -d '{"backend": "cache"}'
```

---

## Running the Application

```bash
# Navigate to project directory
cd learn/08_todo_concurrency

# Install dependencies
go mod download

# Option 1: Development mode with hot reload (RECOMMENDED)
# Automatically restarts when you change code!
make dev

# Option 2: Run directly (manual restart needed)
go run cmd/api/main.go

# Option 3: Run with race detector (detects concurrency bugs)
make race
```

**First time?** Install air for hot reload:
```bash
go install github.com/air-verse/air@latest
```

The server starts on `http://localhost:8080`

> See [DEV_TOOLS.md](./DEV_TOOLS.md) for detailed explanation of Makefile vs Air vs other tools

---

## Code Examples & Learning Points

### 1. Interface Design (domain/repository/todo_repository.go)

```go
// The repository interface - ANY storage can implement this
type TodoRepository interface {
    Create(ctx context.Context, todo *entity.Todo) error
    FindByID(ctx context.Context, id string) (*entity.Todo, error)
    FindAll(ctx context.Context) ([]*entity.Todo, error)
    Update(ctx context.Context, todo *entity.Todo) error
    Delete(ctx context.Context, id string) error
}
```

**Key Learning:** Interfaces define behavior, not implementation. This allows swapping storage backends without changing business logic.

### 2. Goroutines (internal/service/notifier.go)

```go
// Non-blocking notification
func (n *Notifier) SendAsync(todoID string, message string) {
    go func() {
        // This runs in background, doesn't block the API response
        time.Sleep(5 * time.Second)
        fmt.Printf("📧 Notification sent for todo %s: %s\n", todoID, message)
    }()
}
```

**Key Learning:** `go` keyword launches a lightweight thread. The HTTP response returns immediately while notification processes in background.

### 3. Channels (internal/service/batch_processor.go)

```go
// Process multiple todos concurrently
func (bp *BatchProcessor) ProcessBatch(todos []CreateTodoRequest) []Result {
    jobs := make(chan CreateTodoRequest, len(todos))
    results := make(chan Result, len(todos))

    // Start 3 worker goroutines
    for w := 0; w < 3; w++ {
        go worker(jobs, results)
    }

    // Send jobs
    for _, todo := range todos {
        jobs <- todo
    }
    close(jobs)

    // Collect results
    var allResults []Result
    for range todos {
        allResults = append(allResults, <-results)
    }
    return allResults
}
```

**Key Learning:** Channels are pipes for communication. Jobs go in, results come out. Multiple workers process concurrently.

### 4. Mutex (internal/service/stats.go)

```go
type StatsService struct {
    mu           sync.RWMutex  // Protects the fields below
    totalTodos   int
    completedTodos int
}

func (s *StatsService) IncrementTotal() {
    s.mu.Lock()           // Acquire write lock
    defer s.mu.Unlock()   // Release lock when function returns
    s.totalTodos++
}

func (s *StatsService) GetStats() Stats {
    s.mu.RLock()          // Acquire read lock (multiple readers OK)
    defer s.mu.RUnlock()
    return Stats{
        TotalTodos: s.totalTodos,
        Completed:  s.completedTodos,
    }
}
```

**Key Learning:** Mutex prevents data races when multiple goroutines access shared data. RWMutex allows multiple readers OR one writer.

---

## Testing Concurrency

### Test Race Conditions
```bash
# Run with race detector to catch concurrency bugs
go run -race cmd/api/main.go
```

### Load Test Batch Endpoint
```bash
# Create 100 todos concurrently
for i in {1..10}; do
    curl -X POST http://localhost:8080/api/v1/todos/batch \
      -H "Content-Type: application/json" \
      -d '{"todos": [{"title": "Task '$i'"}]}' &
done
wait
```

### Monitor Goroutines
The `/api/v1/stats/goroutines` endpoint shows active goroutines count.

---

## Learning Path

1. **Start with interfaces** - Create a todo, see how the same service uses different storage backends
2. **Try batch operations** - Watch how goroutines process multiple items concurrently
3. **Send notifications** - See non-blocking operations in action
4. **Check statistics** - Understand how mutex protects shared data
5. **Switch backends** - Experience interface polymorphism at runtime
6. **Run with `-race`** - Learn to detect race conditions

---

## Common Patterns Demonstrated

| Pattern | Location | Purpose |
|---------|----------|---------|
| Strategy Pattern | Repository implementations | Swap algorithms (storage) at runtime |
| Worker Pool | BatchProcessor | Process tasks concurrently with limited workers |
| Fan-Out/Fan-In | BatchProcessor | Distribute work, collect results |
| Producer-Consumer | Notifier | Decouple task creation from execution |
| Thread-Safe Singleton | StatsService | Global state with synchronized access |

---

## Next Steps

After mastering these concepts:
1. Add database persistence (PostgreSQL)
2. Implement distributed locks (Redis)
3. Add message queues (RabbitMQ/Kafka)
4. Build microservices with gRPC
5. Add observability (metrics, tracing)

---

## Key Takeaways

- **Interfaces** = Contracts for behavior, enable loose coupling
- **Goroutines** = Lightweight threads, easy concurrency
- **Channels** = Typed pipes for goroutine communication
- **Mutex** = Locks for protecting shared state
- **Context** = Propagate cancellation signals (bonus!)

Happy learning! 🚀
