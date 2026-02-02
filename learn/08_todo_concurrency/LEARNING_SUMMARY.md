# Learning Summary

This document explains what each file teaches and where to find specific concepts.

---

## 📚 What You'll Learn

### 1. Interfaces (Polymorphism & Dependency Injection)

**Core Concept:** Write code that depends on behavior (interface), not implementation.

**Files to Study:**
- `domain/repository/todo_repository.go:10-30` - Interface definition
- `internal/repository/memory/todo_memory.go:15-25` - Implementation 1
- `internal/repository/cache/todo_cache.go:18-28` - Implementation 2
- `internal/service/todo_service.go:15-20` - Using the interface

**Key Learning Points:**
- Interface defines a contract (methods)
- Any type implementing those methods satisfies the interface
- Services depend on interfaces, not concrete types
- Swap implementations without changing business logic

**Try This:**
```bash
# Both use TodoRepository interface, but different implementations
curl http://localhost:8080/api/v1/admin/storage-info
```

---

### 2. Goroutines (Lightweight Concurrency)

**Core Concept:** Launch concurrent operations with `go` keyword.

**Files to Study:**
- `internal/service/batch_processor.go:45-65` - Worker goroutines
- `internal/service/notifier.go:60-85` - Background worker
- `cmd/api/main.go:88-116` - Server goroutine

**Key Learning Points:**
- `go func()` launches a goroutine
- Goroutines are lightweight (can create thousands)
- Don't wait for goroutine to finish unless you want to
- Use `sync.WaitGroup` to wait for multiple goroutines

**Try This:**
```bash
# Watch console for concurrent worker activity
curl -X POST http://localhost:8080/api/v1/todos/batch \
  -d '{"todos":[{"title":"T1","priority":1},...10 more...]}'
```

---

### 3. Channels (Goroutine Communication)

**Core Concept:** Channels are typed pipes for sending/receiving data between goroutines.

**Files to Study:**
- `internal/service/batch_processor.go:30-90` - Worker pool with channels
- `internal/service/notifier.go:35-50` - Notification queue
- `internal/service/batch_processor.go:130-165` - Select statement

**Key Learning Points:**
- Unbuffered channel: `make(chan T)` - blocks until received
- Buffered channel: `make(chan T, N)` - can hold N items
- Send: `ch <- value`
- Receive: `value := <-ch`
- Close: `close(ch)` - signals no more values
- Range: `for v := range ch` - receives until closed

**Channel Patterns Demonstrated:**
1. **Worker Pool** (`batch_processor.go:30-90`)
   - Jobs sent to channel
   - Workers receive and process
   - Results sent to results channel

2. **Producer-Consumer** (`notifier.go:60-85`)
   - Producer queues notifications
   - Consumer processes in background

3. **Fan-out/Fan-in** (`batch_processor.go:45-90`)
   - Distribute work to multiple workers (fan-out)
   - Collect results from all workers (fan-in)

**Try This:**
```bash
# Queue fills up (buffered channel)
curl http://localhost:8080/api/v1/notifications/stats
```

---

### 4. Mutex (Thread-Safe Shared State)

**Core Concept:** Protect shared data from concurrent access with locks.

**Files to Study:**
- `internal/repository/memory/todo_memory.go:30-60` - RWMutex usage
- `internal/service/stats_service.go:30-70` - Thread-safe counters
- `internal/repository/cache/todo_cache.go:60-80` - Cache mutex

**Key Learning Points:**
- **Race condition:** Multiple goroutines access shared data = bugs
- **Mutex:** Mutual exclusion lock
- `mu.Lock()` - Exclusive access (write)
- `mu.RLock()` - Shared access (read)
- Always `defer mu.Unlock()` to prevent deadlocks
- Use `go run -race` to detect races

**Common Patterns:**
```go
// Write operation - exclusive lock
func (s *Service) Update() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.counter++  // Safe!
}

// Read operation - shared lock (multiple readers OK)
func (s *Service) Get() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.counter  // Safe!
}
```

**Try This:**
```bash
# Run with race detector
make race

# Hit stats endpoint from multiple terminals
while true; do curl http://localhost:8080/api/v1/stats; done
```

---

## 🎯 Learning by Endpoint

### Basic CRUD → Interfaces
```
POST   /api/v1/todos       → todo_service.go uses TodoRepository interface
GET    /api/v1/todos       → Works with any implementation
DELETE /api/v1/todos/:id   → No knowledge of storage details
```

### Batch Operations → Goroutines + Channels
```
POST /api/v1/todos/batch   → Worker pool pattern
                            → Multiple goroutines process concurrently
                            → Channels communicate jobs and results
```

### Notifications → Async Goroutines
```
POST /api/v1/todos/:id/notify → Non-blocking operation
                               → Returns immediately
                               → Background worker processes later
```

### Statistics → Mutex
```
GET /api/v1/stats             → Thread-safe counter reads
GET /api/v1/stats/detailed    → Multiple fields read atomically
POST /api/v1/stats/reset      → Safe counter reset
```

---

## 📖 Recommended Learning Order

### Day 1: Interfaces
1. Read `domain/repository/todo_repository.go`
2. Study `internal/repository/memory/todo_memory.go`
3. Compare with `internal/repository/cache/todo_cache.go`
4. See how `internal/service/todo_service.go` uses the interface
5. Try: Create a todo, switch storage, see it still works

### Day 2: Goroutines & Channels
1. Read `internal/service/batch_processor.go` comments
2. Try batch endpoint, watch console output
3. Experiment with worker count (main.go:29)
4. Study the channel operations
5. Try: Create 100 todos in batch, observe concurrent processing

### Day 3: Async Operations
1. Study `internal/service/notifier.go`
2. Send notifications with delays
3. Watch how API returns immediately
4. Check notification stats
5. Try: Queue 50 notifications, monitor progress

### Day 4: Mutex & Thread Safety
1. Read `internal/service/stats_service.go`
2. Run with `-race` flag
3. Hit stats endpoint concurrently
4. Try removing mutex, see races
5. Understand read vs write locks

---

## 🔍 Code Reading Guide

### Understand an Interface
```
1. Look at interface definition (domain/repository/)
2. Find implementations (internal/repository/*/
3. See who uses it (internal/service/)
4. Trace a request through handler → service → repository
```

### Understand a Goroutine
```
1. Find `go func()` or `go worker()`
2. See what data it needs
3. Check how it communicates (channels? shared state?)
4. Find where it's started and stopped
```

### Understand a Channel
```
1. Find channel creation: make(chan T, N)
2. Find senders: ch <- value
3. Find receivers: v := <-ch
4. Check if/when it's closed: close(ch)
```

### Understand Mutex Usage
```
1. Find the mutex: sync.Mutex or sync.RWMutex
2. Find Lock/RLock calls
3. Check what data is protected
4. Verify defer Unlock patterns
```

---

## 🧪 Experiments to Try

### 1. Break the Mutex
Remove mutex from `todo_memory.go`, run with `-race`, observe crashes.

### 2. Change Worker Count
Edit `main.go:29`, set to 1, 3, 10 workers. See performance differences.

### 3. Channel Buffer Size
Edit `notifier.go:40`, change buffer size, see queue behavior.

### 4. Implement File Storage
Create `FileStorageTodoRepository` implementing `TodoRepository`.

### 5. Add Logging
Add a logging goroutine that receives events from channels.

---

## 📊 Project Statistics

- **Total Files:** 18
- **Lines of Code:** ~2,000
- **Concepts Covered:** 4 major (Interfaces, Goroutines, Channels, Mutex)
- **API Endpoints:** 18
- **Learning Examples:** 10+ in EXAMPLES.md

---

## 🎓 After This Project

You'll be able to:
- ✅ Design interface-based architectures
- ✅ Write concurrent code with goroutines
- ✅ Use channels for goroutine communication
- ✅ Protect shared state with mutexes
- ✅ Detect and fix race conditions
- ✅ Build scalable Go backends

### Next Steps:
1. Add database persistence (PostgreSQL + pgx)
2. Implement context cancellation throughout
3. Add distributed locking (Redis)
4. Build with gRPC instead of REST
5. Add observability (Prometheus metrics)

---

## 🐛 Debugging Tips

### Data Race Detected
```bash
go run -race cmd/api/main.go
# Shows exactly where the race is!
```

### Goroutine Leak
```bash
# Check goroutine count over time
curl http://localhost:8080/api/v1/stats/goroutines
# Should stabilize, not keep growing
```

### Deadlock
```
fatal error: all goroutines are asleep - deadlock!
```
- Check for unclosed channels
- Verify all sends have receivers
- Look for missing mutex unlocks

---

## 📚 Further Reading

- **Effective Go:** https://go.dev/doc/effective_go
- **Go Concurrency Patterns:** https://go.dev/blog/pipelines
- **Concurrency is not Parallelism:** https://go.dev/blog/waza-talk
- **Go Memory Model:** https://go.dev/ref/mem

---

Happy learning! Remember: The best way to learn is by doing. Try breaking things, see what happens! 🚀
