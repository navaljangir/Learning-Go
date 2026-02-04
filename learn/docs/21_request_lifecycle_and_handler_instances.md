# Request Lifecycle: When Are Handlers Created?

## The Key Question

When you use `&AuthHandler{userService: userService}`, does it:
- Create a new struct for **every request**?
- Create **one struct** at startup and reuse it?

**Answer:** It creates **ONE struct at startup** and **reuses it for all requests**.

Let's trace through the complete flow.

---

## Timeline: Application Startup to Multiple Requests

### Step 1: Application Startup (ONCE)

```go
// cmd/api/main.go
func main() {
    // 1. Load config (happens ONCE)
    cfg := config.Load()

    // 2. Connect to database (happens ONCE)
    db := initDatabase(cfg)

    // 3. Create repositories (happens ONCE)
    userRepo, todoRepo := initRepositories(db)

    // 4. Create services (happens ONCE)
    userService, todoService := initServices(userRepo, todoRepo, jwtUtil)

    // 5. CREATE HANDLERS (happens ONCE) ⭐
    authHandler, userHandler, todoHandler := initHandlers(userService, todoService)

    // 6. Setup router (happens ONCE)
    router := setupRouter(authHandler, userHandler, todoHandler, jwtUtil)

    // 7. Start server
    srv := createServer(cfg, router)
    srv.ListenAndServe()
}
```

### Step 2: Inside `initHandlers` (ONCE at startup)

```go
func initHandlers(
    userService service.UserService,
    todoService service.TodoService,
) (handler.AuthHandlerInterface, handler.UserHandlerInterface, handler.TodoHandlerInterface) {

    // ⭐ THIS LINE EXECUTES ONCE AT STARTUP
    authHandler := handler.NewAuthHandler(userService)
    //             ↑
    //             Creates ONE instance of *AuthHandler
    //             Memory allocated: 8-16 bytes (just a pointer to userService)

    userHandler := handler.NewTodoHandler(userService)
    todoHandler := handler.NewTodoHandler(todoService)

    log.Println("✓ Handlers initialized")

    return authHandler, userHandler, todoHandler
    //     ↑
    //     Returns the SAME pointer that will be reused for all requests
}
```

### Step 3: Inside `NewAuthHandler` (ONCE at startup)

```go
func NewAuthHandler(userService service.UserService) AuthHandlerInterface {
    // ⭐ THIS EXECUTES ONCE AT STARTUP
    return &AuthHandler{userService: userService}
    //     ↑
    //     Step A: Create AuthHandler struct in memory
    //     Step B: Store userService pointer in it
    //     Step C: Return pointer to this struct
    //
    //     Memory location: Let's say 0x1234abcd
}
```

### Step 4: Router Setup (ONCE at startup)

```go
func SetupRouter(
    authHandler handler.AuthHandlerInterface,  // ← Receives 0x1234abcd
    userHandler handler.UserHandlerInterface,
    todoHandler handler.TodoHandlerInterface,
    jwtUtil *utils.JWTUtil,
) *gin.Engine {
    r := gin.New()

    // ⭐ THIS REGISTERS THE HANDLER METHOD (ONCE)
    auth := v1.Group("/auth")
    {
        auth.POST("/register", authHandler.Register)
        //                     ↑
        //                     Stores a reference to the Register method
        //                     of the handler at 0x1234abcd
        //
        //                     Gin stores: "POST /api/v1/auth/register" → authHandler.Register
    }

    return r
}
```

**At this point:**
```
Memory state:
┌────────────────────────────┐
│ AuthHandler instance       │ ← Created ONCE
│ Memory: 0x1234abcd        │
│ Fields:                    │
│   userService: 0x5678ef   │ ← Points to UserService
└────────────────────────────┘

Gin Router:
┌─────────────────────────────────────────────┐
│ Routes table:                               │
│ POST /api/v1/auth/register → 0x1234abcd.Register  │
│ POST /api/v1/auth/login    → 0x1234abcd.Login     │
└─────────────────────────────────────────────┘
```

---

## Request Handling: Multiple Requests (MANY TIMES)

### Request 1: User Registers

```go
// Client sends:
// POST /api/v1/auth/register
// Body: {"username": "naval", "email": "naval@gmail.com", "password": "Pass123"}

// Step 1: Gin receives request
// Step 2: Gin looks up route: "POST /api/v1/auth/register"
// Step 3: Gin finds: authHandler.Register (pointer 0x1234abcd)
// Step 4: Gin creates NEW gin.Context for this request ⭐

ctx := &gin.Context{
    Request:  req,        // ← New HTTP request
    Writer:   w,          // ← New response writer
    Params:   []Param{},  // ← Empty for this route
    // ... other fields
}

// Step 5: Gin calls the handler method ⭐
authHandler.Register(ctx)
//          ↑
//          Calls method on EXISTING handler instance (0x1234abcd)
//          Does NOT create new AuthHandler!
```

### Inside Handler Method (for Request 1)

```go
// api/handler/auth_handler.go
func (h *AuthHandler) Register(c *gin.Context) {
//     ↑
//     h is the SAME AuthHandler instance created at startup (0x1234abcd)
//     c is a NEW gin.Context created for THIS request

    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    // h.userService is the SAME service from startup
    response, err := h.userService.Register(c.Request.Context(), req)
    if err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    utils.Created(c, response)
}
// After this function returns:
// - gin.Context (c) is discarded/reused
// - AuthHandler (h) remains in memory, unchanged
```

### Request 2: User Logs In (Different User)

```go
// Client sends:
// POST /api/v1/auth/login
// Body: {"username": "john", "password": "Pass456"}

// Step 1: Gin creates ANOTHER NEW gin.Context ⭐
ctx2 := &gin.Context{
    Request:  req2,       // ← Different HTTP request
    Writer:   w2,         // ← Different response writer
    // ...
}

// Step 2: Gin calls the SAME handler instance ⭐
authHandler.Login(ctx2)
//          ↑
//          SAME AuthHandler instance (0x1234abcd)
//          Different gin.Context (ctx2)
```

### Request 3: Another User Registers (Concurrent!)

```go
// Client sends:
// POST /api/v1/auth/register
// Body: {"username": "alice", "password": "Pass789"}

// This happens AT THE SAME TIME as Request 2 above!

// Step 1: Gin creates ANOTHER NEW gin.Context (in different goroutine) ⭐
ctx3 := &gin.Context{
    Request:  req3,
    Writer:   w3,
    // ...
}

// Step 2: Gin calls the SAME handler instance ⭐
//         But in a DIFFERENT GOROUTINE!
go authHandler.Register(ctx3)
//             ↑
//             SAME AuthHandler instance (0x1234abcd)
//             Different gin.Context (ctx3)
//             Different goroutine!
```

---

## Memory Visualization: Timeline

### At Startup (t=0)

```
┌─────────────────────────┐
│ AuthHandler             │ ← Created ONCE
│ Address: 0x1234abcd     │
│ Fields:                 │
│   userService: pointer  │
└─────────────────────────┘
```

### During Request 1 (t=1s)

```
┌─────────────────────────┐
│ AuthHandler             │ ← SAME instance
│ Address: 0x1234abcd     │    (reused)
│ Fields:                 │
│   userService: pointer  │
└─────────────────────────┘
         ↑
         │ method called with
         ↓
┌─────────────────────────┐
│ gin.Context (Request 1) │ ← Created NEW for this request
│ Address: 0xaaaa1111     │
│ Fields:                 │
│   Request: req1         │
│   Writer: w1            │
└─────────────────────────┘
```

### During Request 2 (t=2s)

```
┌─────────────────────────┐
│ AuthHandler             │ ← SAME instance
│ Address: 0x1234abcd     │    (reused again)
│ Fields:                 │
│   userService: pointer  │
└─────────────────────────┘
         ↑
         │ method called with
         ↓
┌─────────────────────────┐
│ gin.Context (Request 2) │ ← Created NEW for this request
│ Address: 0xbbbb2222     │
│ Fields:                 │
│   Request: req2         │
│   Writer: w2            │
└─────────────────────────┘

(gin.Context from Request 1 is gone/garbage collected)
```

### During Concurrent Requests 3 & 4 (t=3s)

```
┌─────────────────────────┐
│ AuthHandler             │ ← SAME instance
│ Address: 0x1234abcd     │    (shared by both goroutines!)
│ Fields:                 │
│   userService: pointer  │ ← Read-only, safe for concurrent access
└─────────────────────────┘
         ↑                 ↑
         │                 │
    Goroutine 1       Goroutine 2
         │                 │
         ↓                 ↓
┌───────────────┐   ┌───────────────┐
│ gin.Context 3 │   │ gin.Context 4 │ ← Two NEW contexts
│ Addr: 0xcccc  │   │ Addr: 0xdddd  │
│ Request: req3 │   │ Request: req4 │
└───────────────┘   └───────────────┘
```

---

## The Critical Question: Is This Safe?

### Yes, It's Safe! Here's Why:

#### 1. **Handler Instance is Read-Only**

```go
type AuthHandler struct {
    userService service.UserService  // ← This NEVER changes after startup
}

func (h *AuthHandler) Register(c *gin.Context) {
    // h.userService is only READ, never modified
    response, err := h.userService.Register(c.Request.Context(), req)
    //               ↑
    //               Reading h.userService (safe for concurrent access)
}
```

#### 2. **Request Data is in gin.Context (Not Handler)**

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    c.ShouldBindJSON(&req)  // ← Request data stored in LOCAL variable
    //                          Each goroutine has its own 'req'

    // req is on this goroutine's stack, not shared
}
```

#### 3. **Each Request Gets Its Own gin.Context**

```go
// Request 1 (Goroutine 1):
func (h *AuthHandler) Register(c1 *gin.Context) {
    var req dto.RegisterRequest  // ← Stack memory for Goroutine 1
    c1.ShouldBindJSON(&req)      // ← Reading from c1
}

// Request 2 (Goroutine 2) - happens at same time:
func (h *AuthHandler) Register(c2 *gin.Context) {
    var req dto.RegisterRequest  // ← Different stack memory for Goroutine 2
    c2.ShouldBindJSON(&req)      // ← Reading from c2
}

// No shared mutable state!
```

---

## What If Handler Had Mutable State? (DANGEROUS!)

### ❌ UNSAFE Example

```go
// DON'T DO THIS!
type BadAuthHandler struct {
    userService  service.UserService
    requestCount int  // ← Mutable state!
}

func (h *BadAuthHandler) Register(c *gin.Context) {
    h.requestCount++  // ❌ RACE CONDITION!
    //               Multiple goroutines modifying same variable

    fmt.Printf("Request #%d\n", h.requestCount)  // ❌ Unpredictable output
}
```

**Why is this unsafe?**

```
Time    Goroutine 1           Goroutine 2         h.requestCount
────────────────────────────────────────────────────────────────
t=0                                               0
t=1     Read: h.requestCount=0
t=2                           Read: h.requestCount=0
t=3     Increment: 0+1=1
t=4                           Increment: 0+1=1
t=5     Write: h.requestCount=1
t=6                           Write: h.requestCount=1
t=7                                               1 (should be 2!)
```

### ✅ SAFE Example (Our Current Design)

```go
type AuthHandler struct {
    userService service.UserService  // ← Immutable after startup
    // No mutable state!
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest  // ← Local variable (stack)
    c.ShouldBindJSON(&req)       // ← Request-specific context

    // All data is either:
    // - Read-only (h.userService)
    // - Local to this goroutine (req, c)
}
```

---

## Complete Request Flow Diagram

```
Application Startup (ONCE):
═══════════════════════════

main()
  │
  ├─> initDatabase()
  ├─> initRepositories()
  ├─> initServices()
  │
  ├─> initHandlers()                    ┌────────────────────┐
  │     │                               │ AuthHandler        │
  │     └─> NewAuthHandler(userService) │ Memory: 0x1234abcd │
  │           │                         │ Created ONCE       │
  │           └─> return &AuthHandler{...} ───────────────>  └────────────────────┘
  │                                                            ↑
  ├─> setupRouter(authHandler, ...)                          │
  │     │                                                      │
  │     └─> auth.POST("/register", authHandler.Register) ─────┘
  │           (Stores reference to method)                    Stored in
  │                                                            Gin's route table
  └─> srv.ListenAndServe()


Request Handling (MANY TIMES):
═══════════════════════════════

Request 1 (t=1s):
  HTTP POST /api/v1/auth/register
    │
    ├─> Gin receives request
    ├─> Gin creates gin.Context #1 (0xaaaa1111) ───────┐
    │                                                   │
    ├─> Gin looks up route table:                      │
    │   "POST /auth/register" → authHandler.Register   │
    │                            (0x1234abcd.Register)  │
    │                                                   │
    └─> Gin calls: authHandler.Register(ctx1) ─────────┼──> func (h *AuthHandler) Register(c *gin.Context)
              Uses existing handler ─────┐             │         │ h = 0x1234abcd (reused)
                                         │             │         │ c = 0xaaaa1111 (new)
                                         ↓             │         │
                              ┌────────────────────┐  │         ├─> var req dto.RegisterRequest
                              │ AuthHandler        │  │         ├─> c.ShouldBindJSON(&req)
                              │ Memory: 0x1234abcd │◄─┘         ├─> h.userService.Register(...)
                              │ Created at startup │            └─> utils.Created(c, response)
                              └────────────────────┘
                                       ↑
Request 2 (t=2s):                     │
  HTTP POST /api/v1/auth/login        │ SAME handler instance!
    │                                  │
    ├─> Gin creates gin.Context #2 (0xbbbb2222) ───────┐
    │                                                   │
    └─> Gin calls: authHandler.Login(ctx2) ────────────┼──> func (h *AuthHandler) Login(c *gin.Context)
              Uses SAME handler ──────┘                 │         │ h = 0x1234abcd (reused!)
                                                       │         │ c = 0xbbbb2222 (new, different from ctx1)
                                                       │         │
                                                       │         └─> Process login...
                                                       │
                                            gin.Context #1 (0xaaaa1111) is gone
                                            (garbage collected after Request 1 finished)


Request 3 & 4 (t=3s) - CONCURRENT:
════════════════════════════════════

Request 3 (Goroutine 1):              Request 4 (Goroutine 2):
  HTTP POST /auth/register              HTTP POST /auth/register
    │                                     │
    ├─> gin.Context #3 (0xcccc3333)      ├─> gin.Context #4 (0xdddd4444)
    │                                     │
    └─> authHandler.Register(ctx3)       └─> authHandler.Register(ctx4)
              │                                     │
              │         ┌────────────────────┐     │
              └────────>│ AuthHandler        │<────┘
                        │ Memory: 0x1234abcd │
                        │ Shared by both!    │
                        └────────────────────┘
                                 │
                         Only reads userService
                         (safe for concurrent access)
```

---

## Summary Table

| What | When Created | How Many | Lifetime | Shared? |
|------|-------------|----------|----------|---------|
| **AuthHandler** | Startup (once) | 1 instance | Entire application | ✅ Yes (safe - read-only) |
| **gin.Context** | Per request | 1 per request | During request only | ❌ No (each request gets own) |
| **Local variables** | Per request | 1 per request | During request only | ❌ No (goroutine stack) |
| **userService** | Startup (once) | 1 instance | Entire application | ✅ Yes (safe - read-only) |

---

## Key Takeaways

1. **Handler created ONCE at startup**
   ```go
   authHandler := handler.NewAuthHandler(userService)  // Executes ONCE
   ```

2. **Handler REUSED for all requests**
   ```go
   authHandler.Register(ctx1)  // Request 1
   authHandler.Register(ctx2)  // Request 2 - SAME handler
   authHandler.Register(ctx3)  // Request 3 - SAME handler
   ```

3. **gin.Context created NEW for each request**
   ```go
   // Each request gets its own gin.Context
   // Request data stored here, not in handler
   ```

4. **Safe because handler is read-only**
   ```go
   type AuthHandler struct {
       userService service.UserService  // ← Never modified
   }
   // All mutable data is in gin.Context or local variables
   ```

5. **Concurrent requests safe**
   ```go
   // Multiple goroutines can safely call methods on same handler
   // Because they only READ handler fields, never WRITE
   ```

---

## Node.js/Express Comparison

### Express (JavaScript)

```javascript
// Startup (once)
const authHandler = new AuthHandler(userService);  // Created ONCE

app.post('/auth/register', authHandler.register);  // Store reference

// Per request (many times)
// Express calls: authHandler.register(req, res)
//                    ↑                 ↑
//              SAME handler      NEW req/res for each request
```

**Same pattern!** Handler created once, request objects created per-request.

### Key Difference

- **Go**: Each request handled in separate goroutine (concurrent)
- **Node.js**: Single-threaded event loop (not concurrent)

But the handler instance lifecycle is the same!

---

This is the foundation of how HTTP servers work efficiently - create expensive objects (handlers, services, db connections) once at startup, create cheap objects (request contexts) per-request.
