# Async Operations in Go

## What Does "Warm Connection" Mean?

**"Warm" = Already connected and ready to use**

### Cold Connection (Slow)
```
Request comes in
    ↓
Need DB connection
    ↓
Create TCP socket           ← Takes time (~10-50ms)
    ↓
TCP handshake (3-way)      ← Network round trip
    ↓
MySQL authentication       ← Another round trip
    ↓
Ready to query
```

### Warm Connection (Fast)
```
Request comes in
    ↓
Need DB connection
    ↓
Grab from pool (already connected)  ← Instant! (~0.1ms)
    ↓
Ready to query
```

### Visual Comparison

```go
// After handling request:

// WITHOUT SetMaxIdleConns (closes all connections)
┌─────────────────────────────────┐
│  Connection Pool                │
│  [empty] [empty] [empty]        │  ← All closed
└─────────────────────────────────┘
Next request: Must create new connection (slow!)

// WITH SetMaxIdleConns(5)
┌─────────────────────────────────┐
│  Connection Pool                │
│  [Conn1] [Conn2] [Conn3]        │  ← 5 stay open
│  [Conn4] [Conn5]                │     (warm & ready)
└─────────────────────────────────┘
Next request: Grabs Conn1 instantly (fast!)
```

**Trade-off:**
- More idle connections = faster response, but uses DB resources
- Fewer idle connections = slower response, but saves DB resources

---

## When to Use Goroutines for "Fire and Forget"

Use goroutines when you don't want to wait for an operation:
- Cache updates (Redis)
- Analytics events
- Logging to external services
- Sending emails
- Webhook notifications

---

## Example: Redis Cache Update (Async)

```go
func (s *TodoServiceImpl) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
    priority := entity.Priority(req.Priority)
    todo := entity.NewTodo(userID, req.Title, req.Description, priority, req.DueDate)

    // 1. Save to database (MUST wait for this - critical)
    if err := s.todoRepo.Create(ctx, todo); err != nil {
        return nil, err
    }

    // 2. Update Redis cache (fire and forget - don't block response)
    go func() {
        // Use context.Background() because original ctx might be cancelled
        cacheCtx := context.Background()

        cacheKey := fmt.Sprintf("user:%s:todos", userID)
        if err := s.redisClient.Del(cacheCtx, cacheKey).Err(); err != nil {
            log.Printf("Failed to invalidate cache: %v", err)
            // Don't crash - cache failure shouldn't break the app
        }
    }()

    // 3. Send analytics event (fire and forget)
    go func() {
        s.analytics.Track("todo_created", map[string]interface{}{
            "user_id": userID,
            "todo_id": todo.ID,
        })
    }()

    response := dto.TodoToResponse(todo)
    return &response, nil
}
```

**What happens:**
```
Request timeline:
0ms   → Handler receives request
1ms   → Service.Create() called
2ms   → todoRepo.Create() starts (BLOCKS, waits for DB)
52ms  → DB responds, row inserted
53ms  → Launch goroutine #1 for Redis (DOESN'T WAIT)
53ms  → Launch goroutine #2 for analytics (DOESN'T WAIT)
54ms  → Return response to client ← Client gets response fast!

Meanwhile (in background):
53ms  → Goroutine #1 runs: Redis cache invalidation
100ms → Goroutine #1 completes (client already got response)

53ms  → Goroutine #2 runs: Send analytics
120ms → Goroutine #2 completes (client already got response)
```

---

## Common Pattern: Background Worker

For truly critical async work (like sending emails), use a proper pattern:

```go
// Option 1: Simple fire-and-forget with error logging
go func() {
    if err := sendEmail(user.Email, "Welcome!"); err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}()

// Option 2: With timeout (better)
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := sendEmailWithContext(ctx, user.Email, "Welcome!"); err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}()

// Option 3: Queue-based (production recommended)
// Use a proper job queue like:
// - Redis-based: asynq, gocraft/work
// - Database-based: river (PostgreSQL)
// - External: RabbitMQ, AWS SQS
s.jobQueue.Enqueue("send_email", map[string]interface{}{
    "to": user.Email,
    "template": "welcome",
})
```

---

## Real Example: Audit Logging (Don't Block User)

```go
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Capture old values for audit
    oldFullName := user.FullName

    // Update user
    user.Update(req.FullName)

    // Save to DB (MUST wait)
    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }

    // Audit log (fire and forget)
    go func() {
        auditCtx := context.Background()
        s.auditLogger.Log(auditCtx, AuditLog{
            UserID:    userID,
            Action:    "profile_updated",
            OldValue:  oldFullName,
            NewValue:  user.FullName,
            Timestamp: time.Now(),
        })
    }()

    response := dto.UserToResponse(user)
    return &response, nil
}
```

---

## ⚠️ Important Gotchas

### 1. Don't use request context in goroutines
```go
// ❌ BAD - ctx might be cancelled when goroutine runs
go func() {
    s.redisClient.Set(ctx, key, value, 0)  // ctx is from HTTP request
}()

// ✅ GOOD - use fresh context
go func() {
    bgCtx := context.Background()
    s.redisClient.Set(bgCtx, key, value, 0)
}()
```

**Why?** When the HTTP response is sent, the request context is cancelled. Your goroutine might fail if it tries to use that context.

### 2. Don't access request/response objects
```go
// ❌ BAD - gin.Context is not safe for concurrent use
go func() {
    c.Set("cache_updated", true)  // RACE CONDITION!
}()

// ✅ GOOD - copy data you need
emailAddr := req.Email
go func() {
    sendEmail(emailAddr)  // Use copied value
}()
```

### 3. Handle panics in goroutines
```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Goroutine panic: %v", r)
            // Optionally: send to error tracking service
        }
    }()

    // Risky operation that might panic
    s.riskyOperation()
}()
```

---

## When NOT to Use Goroutines

Don't use goroutines for:

### 1. Database writes - You need to know if they succeeded
```go
// ❌ BAD - What if this fails?
go s.todoRepo.Create(ctx, todo)
return &response, nil  // User thinks it saved, but might have failed!
```

### 2. Operations affecting the response
```go
// ❌ BAD - Response sent before calculation
go calculateTotals(&response)
return response, nil  // Totals will be empty!
```

### 3. Auth/authorization checks
```go
// ❌ BAD - Race condition
go checkPermissions(userID)
return secretData, nil  // Returned before permission check!
```

---

## Production Pattern: Worker Pool

For high-volume async work, use a worker pool:

```go
type AsyncWorker struct {
    jobQueue chan Job
    workers  int
}

type Job struct {
    Type string
    Data interface{}
}

func NewAsyncWorker(workers int) *AsyncWorker {
    w := &AsyncWorker{
        jobQueue: make(chan Job, 1000),  // Buffer 1000 jobs
        workers:  workers,
    }
    w.start()
    return w
}

func (w *AsyncWorker) start() {
    for i := 0; i < w.workers; i++ {
        go func(workerID int) {
            for job := range w.jobQueue {
                w.processJob(job, workerID)
            }
        }(i)
    }
}

func (w *AsyncWorker) Submit(job Job) {
    select {
    case w.jobQueue <- job:
        // Job submitted
    default:
        log.Printf("Job queue full, dropping job: %+v", job)
    }
}

func (w *AsyncWorker) processJob(job Job, workerID int) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Worker %d panic: %v", workerID, r)
        }
    }()

    log.Printf("Worker %d processing %s", workerID, job.Type)

    switch job.Type {
    case "cache_invalidate":
        // Handle cache invalidation
    case "send_email":
        // Handle email
    case "analytics":
        // Handle analytics
    }
}

// Usage:
asyncWorker := NewAsyncWorker(10)  // 10 concurrent workers

func (s *TodoService) Create(...) {
    // ... save todo ...

    asyncWorker.Submit(Job{
        Type: "cache_invalidate",
        Data: map[string]interface{}{"user_id": userID},
    })
}
```

---

## Summary Table

| Operation | Pattern | Why |
|-----------|---------|-----|
| Database writes | `result, err := db.Create()` | Must know if succeeded |
| Database reads | `result, err := db.Find()` | Need the data for response |
| Redis cache update | `go redis.Del()` | Don't need to wait, failure OK |
| Send email | `go sendEmail()` or queue | Don't block, retry on failure |
| Analytics tracking | `go analytics.Track()` | Don't block, loss acceptable |
| Audit logging | Queue preferred | Must not lose data |
| Webhooks | Queue preferred | Need retries |

**Golden Rule:** If failure matters or you need the result, DON'T use `go`. If it's fire-and-forget and failure is acceptable, DO use `go`.

---

## Key Takeaways

1. **DB calls are synchronous** - They block the current goroutine (but that's OK because goroutines are cheap)
2. **Connection pooling handles efficiency** - You don't manually create goroutines for DB calls
3. **"Warm connections"** - Keeping idle connections open speeds up subsequent requests
4. **Use `go` for async work** - Cache updates, logging, analytics that don't need to block the response
5. **Be careful with contexts** - Don't use request context in background goroutines
6. **For critical async work** - Use proper job queues instead of simple goroutines
