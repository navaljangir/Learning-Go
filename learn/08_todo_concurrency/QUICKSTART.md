# Quick Start Guide

Get up and running in 5 minutes!

## Prerequisites

- Go 1.21 or higher
- curl (for testing)

## Installation

```bash
# Navigate to project
cd learn/08_todo_concurrency

# Download dependencies
go mod download

# Build (optional)
make build
```

## Run the Server

```bash
# Option 1: Development mode with hot reload (RECOMMENDED)
# Automatically restarts when you change code!
make dev

# Option 2: Run directly (manual restart needed)
go run cmd/api/main.go

# Option 3: Use Makefile
make run

# Option 4: Run with race detector (detects concurrency bugs!)
make race
```

**First time?** Install air for hot reload:
```bash
go install github.com/air-verse/air@latest
```

You should see:
```
╔════════════════════════════════════════════════════════════════╗
║           TODO APP - Learning Concurrency & Interfaces        ║
╚════════════════════════════════════════════════════════════════╝

🚀 Server starting on http://localhost:8080
```

## Quick Test

Open another terminal and try these commands:

### 1. Create a Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go Interfaces","description":"Understand interface contracts","priority":3}'
```

### 2. List Todos
```bash
curl http://localhost:8080/api/v1/todos
```

### 3. Batch Create (Watch Console!)
```bash
curl -X POST http://localhost:8080/api/v1/todos/batch \
  -H "Content-Type: application/json" \
  -d '{
    "todos": [
      {"title":"Task 1","priority":2},
      {"title":"Task 2","priority":1},
      {"title":"Task 3","priority":3}
    ]
  }'
```

**Watch your server console - you'll see workers processing concurrently!**

### 4. Send Async Notification
```bash
curl -X POST http://localhost:8080/api/v1/todos/1/notify \
  -H "Content-Type: application/json" \
  -d '{"message":"Don't forget!","delay_seconds":5}'
```

**The API responds immediately, but watch console - notification appears after 5 seconds!**

### 5. View Statistics
```bash
curl http://localhost:8080/api/v1/stats | jq
```

## Learning Path

1. **Start here:** Read `README.md` for concepts overview
2. **Hands-on:** Follow `EXAMPLES.md` for detailed walkthroughs
3. **Deep dive:** Study the code with comments in:
   - `domain/repository/todo_repository.go` - Interface design
   - `internal/repository/memory/todo_memory.go` - Mutex usage
   - `internal/service/batch_processor.go` - Goroutines & channels
   - `internal/service/notifier.go` - Async operations
   - `internal/service/stats_service.go` - Thread-safe counters

## Key Concepts Demonstrated

### Interfaces (Polymorphism)
```
TodoRepository (interface)
    ├── InMemoryTodoRepository
    └── CachedTodoRepository
```
Same interface, different implementations!

### Goroutines (Concurrency)
```go
go func() {
    // This runs concurrently!
}()
```

### Channels (Communication)
```go
jobs := make(chan Todo, 10)  // Buffered channel
jobs <- todo                  // Send
todo := <-jobs                // Receive
```

### Mutex (Thread Safety)
```go
mu.Lock()                     // Exclusive access
defer mu.Unlock()
count++                       // Safe increment
```

## Common Commands

```bash
make run          # Run the server
make race         # Run with race detector
make build        # Build binary
make test-batch   # Test batch endpoint
make test-stats   # View statistics
```

## Troubleshooting

### Port already in use
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

### Race detector warnings
**This is good!** It means you're learning. The warnings show potential concurrency bugs.

### Module errors
```bash
go mod tidy
go clean -modcache
```

## What's Next?

1. Try all examples in `EXAMPLES.md`
2. Experiment with worker count in `main.go` (line 29)
3. Add your own endpoints
4. Implement file storage (see Exercise 1 in EXAMPLES.md)
5. Learn more in Go docs: https://go.dev/doc/

## Need Help?

- Check `README.md` for architecture details
- See `EXAMPLES.md` for detailed examples
- Study code comments in each file
- Run with `-race` to catch concurrency bugs

Happy learning! 🚀
