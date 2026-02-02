# Learning Examples & Exercises

This guide walks you through practical examples to understand interfaces, goroutines, channels, and mutex.

---

## Setup

1. Start the server:
```bash
cd learn/08_todo_concurrency
go mod download
go run cmd/api/main.go
```

2. Open another terminal for running commands

---

## Part 1: Understanding Interfaces

### Example 1: Create a Todo (Interface in Action)

```bash
# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go Interfaces",
    "description": "Understand interface contracts",
    "priority": 3
  }'
```

**What's happening:**
- `TodoHandler` calls `TodoService`
- `TodoService` calls `TodoRepository` **interface**
- The **actual implementation** is `InMemoryTodoRepository`
- But the service doesn't know or care which implementation it is!

**Key Learning:** The service depends on the interface, not the concrete type. This is dependency inversion!

---

### Example 2: Check Storage Type

```bash
curl http://localhost:8080/api/v1/admin/storage-info
```

**Response:**
```json
{
  "implements_storage_info": true,
  "storage_type": "in-memory",
  "stats": {
    "storage_type": "in-memory",
    "total_todos": 1,
    "access_count": 5
  }
}
```

**Key Learning:** We use type assertion `repo.(repository.StorageInfo)` to check if the repository implements optional methods.

---

### Exercise 1: Create Your Own Implementation

Challenge: Create a `FileStorageTodoRepository` that saves todos to a JSON file.

Steps:
1. Create `internal/repository/file/todo_file.go`
2. Implement all methods of `TodoRepository` interface
3. Use a mutex to protect file access
4. Update `main.go` to use your implementation

**Hint:** Your struct should look like:
```go
type FileStorageTodoRepository struct {
    mu       sync.RWMutex
    filename string
}
```

---

## Part 2: Understanding Goroutines & Channels

### Example 3: Batch Processing (Worker Pool Pattern)

```bash
# Create 10 todos concurrently
curl -X POST http://localhost:8080/api/v1/todos/batch \
  -H "Content-Type: application/json" \
  -d '{
    "todos": [
      {"title": "Task 1", "priority": 2},
      {"title": "Task 2", "priority": 1},
      {"title": "Task 3", "priority": 3},
      {"title": "Task 4", "priority": 2},
      {"title": "Task 5", "priority": 1}
    ]
  }'
```

**Watch the console output:**
```
🚀 Starting batch processing of 5 todos...
📤 Sent job 1 to queue
📤 Sent job 2 to queue
🔨 Worker 1 processing job: Task 1
🔨 Worker 2 processing job: Task 2
🔨 Worker 3 processing job: Task 3
✅ Worker 1 completed: Task 1 (ID: 1)
🔨 Worker 1 processing job: Task 4
...
```

**Key Learning:**
- Jobs sent to buffered channel
- 3 workers receive jobs from channel
- Workers process concurrently
- Results sent back through results channel

**Open `internal/service/batch_processor.go:30-60` to see the code!**

---

### Example 4: Async Notifications (Non-blocking Operations)

```bash
# First, create a todo and note its ID
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Important Task","priority":3}'

# Send notification with 5 second delay
curl -X POST http://localhost:8080/api/v1/todos/1/notify \
  -H "Content-Type: application/json" \
  -d '{"message":"Don't forget this!","delay_seconds":5}'
```

**What happens:**
1. API responds **immediately** (status 202)
2. Watch console - notification sent after 5 seconds!
3. The HTTP handler didn't wait - it returned right away

**Key Learning:**
- Notification queued in channel (line: `internal/service/notifier.go:68`)
- Background worker goroutine processes queue
- Main goroutine returns immediately
- This is **non-blocking async** operation!

---

### Example 5: Monitor Goroutines

While batch processing is running, check goroutine count:

```bash
# Terminal 1: Start batch processing
curl -X POST http://localhost:8080/api/v1/todos/batch \
  -H "Content-Type: application/json" \
  -d '{"todos": [/* 20 todos */]}'

# Terminal 2: Check goroutines (run this immediately)
watch -n 0.5 'curl -s http://localhost:8080/api/v1/stats/goroutines'
```

**You'll see:**
- Goroutine count increases during processing
- Count decreases when workers finish
- Base goroutines: ~5-10 (server + notification worker)
- During batch: +3 (worker goroutines)

**Key Learning:** Goroutines are lightweight. You can create thousands without issues!

---

### Exercise 2: Create a Pipeline

Challenge: Build a pipeline that:
1. Reads todos from input channel
2. Filters by priority (>= 2)
3. Transforms titles to uppercase
4. Sends to output channel

**Hint:**
```go
func pipeline(input <-chan *entity.Todo) <-chan *entity.Todo {
    output := make(chan *entity.Todo)

    go func() {
        defer close(output)
        for todo := range input {
            if todo.Priority >= 2 {
                todo.Title = strings.ToUpper(todo.Title)
                output <- todo
            }
        }
    }()

    return output
}
```

---

## Part 3: Understanding Mutex

### Example 6: Race Condition Detection

Run the server with race detector:

```bash
go run -race cmd/api/main.go
```

Now, hit the stats endpoint from multiple terminals simultaneously:

```bash
# Terminal 1
while true; do curl -s http://localhost:8080/api/v1/stats > /dev/null; done

# Terminal 2
while true; do curl -s http://localhost:8080/api/v1/stats > /dev/null; done

# Terminal 3
while true; do curl -s http://localhost:8080/api/v1/stats > /dev/null; done
```

**If mutex was NOT used:**
```
==================
WARNING: DATA RACE
Write at 0x... by goroutine 47:
  main.StatsService.RecordRequest()
      /path/to/stats.go:28

Previous read at 0x... by goroutine 46:
  main.StatsService.GetStats()
      /path/to/stats.go:45
==================
```

**With mutex:**
- No race detected!
- Counters increment correctly
- Multiple goroutines safely access shared data

**Key Learning:** Without mutex, concurrent reads/writes cause data races. Go's race detector catches these!

---

### Example 7: Watch Stats Update

```bash
# Terminal 1: Create todos continuously
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/todos \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"Task $i\",\"priority\":2}" &
  sleep 0.1
done

# Terminal 2: Watch stats update (thread-safe!)
watch -n 0.5 'curl -s http://localhost:8080/api/v1/stats/detailed | jq'
```

**You'll see:**
```json
{
  "request_count": 47,
  "create_requests": 47,
  "read_requests": 12,
  "active_goroutines": 8,
  "memory_alloc_mb": 3.45
}
```

**Key Learning:**
- `request_count` always accurate (no lost updates)
- Mutex ensures atomic increments
- Multiple counters updated safely

---

### Example 8: Cache Statistics

Switch to cached backend and see cache-specific stats:

```bash
# Switch to cache backend
curl -X POST http://localhost:8080/api/v1/admin/switch-storage \
  -H "Content-Type: application/json" \
  -d '{"backend":"cache"}'

# Add some todos
for i in {1..20}; do
  curl -X POST http://localhost:8080/api/v1/todos \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"Task $i\",\"priority\":2}"
done

# Check cache stats
curl http://localhost:8080/api/v1/stats/storage | jq
```

**Response:**
```json
{
  "storage_type": "cached",
  "total_todos": 20,
  "max_size": 100,
  "hits": 15,
  "misses": 5,
  "evictions": 0,
  "hit_rate": "75.00%"
}
```

**Key Learning:**
- Different interface implementations expose different features
- Cache tracks hits/misses with mutex
- Mutex protects both map AND counters

---

### Exercise 3: Break the Mutex (Intentionally)

Challenge: Remove the mutex from `InMemoryTodoRepository` and see what happens.

Steps:
1. Copy `internal/repository/memory/todo_memory.go` to `todo_memory_broken.go`
2. Remove all `mu.Lock()` and `mu.RLock()` calls
3. Run with `-race` flag
4. Hit the API concurrently
5. Watch race detector catch the bugs!

**Expected output:**
```
WARNING: DATA RACE
```

**Lesson:** Mutex is NOT optional for concurrent access to shared data!

---

## Part 4: Advanced Patterns

### Example 9: Context Cancellation

```bash
# This will be implemented in the notification timeout example
# See internal/service/notifier.go:SendWithTimeout
```

**Key Concepts:**
- Context carries deadlines
- Can cancel long-running operations
- Propagates across goroutine boundaries

---

### Example 10: Select Statement

Open `internal/service/batch_processor.go:140-160` and study the `select` statement:

```go
select {
case result, ok := <-successChan:
    // Handle success
case result, ok := <-errorChan:
    // Handle error
}
```

**Key Learning:**
- `select` waits on multiple channels
- Whichever is ready first executes
- Similar to `switch` but for channels

---

## Testing Scenarios

### Scenario 1: Load Test

Create 1000 todos concurrently:

```bash
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/todos/batch \
    -H "Content-Type: application/json" \
    -d '{"todos":[{"title":"T1","priority":1},{"title":"T2","priority":2}]}' &
done
wait

# Check stats
curl http://localhost:8080/api/v1/stats | jq
```

**Monitor:**
- Response times
- Goroutine count
- Memory usage
- All statistics accurate? (Thanks to mutex!)

---

### Scenario 2: Stress Test Notifications

```bash
# Queue 50 notifications
for i in {1..50}; do
  curl -X POST "http://localhost:8080/api/v1/todos/$i/notify" \
    -H "Content-Type: application/json" \
    -d '{"message":"Reminder '$i'","delay_seconds":1}' &
done

# Check notification stats
curl http://localhost:8080/api/v1/notifications/stats | jq
```

**Key Learning:**
- Buffered channel prevents blocking
- Worker processes queue in background
- Main goroutines don't wait

---

## Quiz Questions

After going through examples, test your understanding:

1. **Interfaces:**
   - Q: Can a service depend on multiple repository interfaces?
   - A: Yes! Each repository (user, todo, etc.) has its own interface

2. **Goroutines:**
   - Q: What happens if you forget to close a channel?
   - A: Goroutines waiting to receive will block forever (goroutine leak)

3. **Channels:**
   - Q: Buffered vs unbuffered channel - when to use which?
   - A: Buffered when producer/consumer speeds differ, unbuffered for synchronization

4. **Mutex:**
   - Q: When to use RLock vs Lock?
   - A: RLock for reads (multiple OK), Lock for writes (exclusive)

---

## Next Steps

1. Implement file storage (Exercise 1)
2. Build a pipeline (Exercise 2)
3. Add your own endpoint using these concepts
4. Study the code comments in each file
5. Experiment with different worker counts
6. Try removing mutex and observe races

---

## Key Files to Study

| File | Concept | Lines to Focus On |
|------|---------|------------------|
| `domain/repository/todo_repository.go` | Interfaces | 10-30 |
| `internal/repository/memory/todo_memory.go` | Mutex | 30-60, 70-80 |
| `internal/service/batch_processor.go` | Goroutines & Channels | 30-90 |
| `internal/service/notifier.go` | Async Operations | 50-100 |
| `internal/service/stats_service.go` | Thread-safe Counters | 30-70 |

Happy learning! 🚀
