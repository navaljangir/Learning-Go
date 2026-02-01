# Go Deep Concepts

Comprehensive guide covering OOP, pointers, struct tags, types, validation, logging, and more.

---

## Table of Contents
1. [OOP in Go](#oop-in-go)
2. [Pointers & References (15 Examples)](#pointers--references)
3. [Struct Tags (json:"email")](#struct-tags)
4. [Number Types (int64, etc.)](#number-types)
5. [Validation Libraries](#validation-libraries)
6. [Logging in Go](#logging-in-go)
7. [gin.New() vs gin.Default()](#ginnew-vs-gindefault)
8. [JWT RegisteredClaims (Embedding)](#jwt-registeredclaims)
9. [Nested Packages & Imports](#nested-packages--imports)
10. [Path Aliases (Not Supported)](#path-aliases)

---

## OOP in Go

Go is **NOT** a traditional OOP language like Java/C++. It doesn't have:
- Classes
- Inheritance
- Constructors
- `this` keyword

But Go achieves similar results through:
- **Structs** (like classes without inheritance)
- **Methods** (functions attached to types)
- **Interfaces** (polymorphism)
- **Embedding** (composition over inheritance)

### Node.js/TypeScript vs Go Comparison

```typescript
// TypeScript - Traditional OOP
class User {
  private id: number;
  public name: string;

  constructor(id: number, name: string) {
    this.id = id;
    this.name = name;
  }

  greet(): string {
    return `Hello, ${this.name}`;
  }
}

class Admin extends User {
  role: string = "admin";

  greet(): string {
    return `Hello Admin ${this.name}`;
  }
}
```

```go
// Go - Composition-based approach
type User struct {
    ID   int    // Public (capital letter)
    Name string // Public
    password string // private (lowercase)
}

// Method attached to User (like class method)
func (u User) Greet() string {
    return "Hello, " + u.Name
}

// "Constructor" pattern (just a function)
func NewUser(id int, name string) *User {
    return &User{ID: id, Name: name}
}

// Admin "inherits" from User via embedding
type Admin struct {
    User       // Embedded - gets all User fields and methods
    Role string
}

// Override method
func (a Admin) Greet() string {
    return "Hello Admin " + a.Name  // Can access User.Name directly
}
```

### The Four Pillars of OOP in Go

| OOP Concept | Go Equivalent |
|-------------|---------------|
| **Encapsulation** | Capital letter = public, lowercase = private |
| **Abstraction** | Interfaces |
| **Inheritance** | Embedding (composition) |
| **Polymorphism** | Interfaces |

### Interfaces (Polymorphism)

```go
// Interface - defines behavior, not data
type Greeter interface {
    Greet() string
}

// Both User and Admin implement Greeter (implicitly!)
// No "implements" keyword needed

func SayHello(g Greeter) {
    fmt.Println(g.Greet())
}

func main() {
    user := User{Name: "Tejas"}
    admin := Admin{User: User{Name: "Admin"}, Role: "super"}

    SayHello(user)  // Works! User has Greet()
    SayHello(admin) // Works! Admin has Greet()
}
```

### Methods with Pointer vs Value Receiver

```go
type Counter struct {
    count int
}

// Value receiver - gets a COPY (can't modify original)
func (c Counter) GetCount() int {
    return c.count
}

// Pointer receiver - gets original (CAN modify)
func (c *Counter) Increment() {
    c.count++  // Modifies the actual Counter
}

func main() {
    counter := Counter{count: 0}
    counter.Increment()  // Go automatically converts to (&counter).Increment()
    fmt.Println(counter.GetCount())  // 1
}
```

---

## Pointers & References

### Quick Refresher
```go
x := 10

&x      // "address of" x → returns a pointer (memory address)
*p      // "value at" pointer → returns the actual value

var p *int = &x   // p is a pointer to int, holding address of x
*p = 20           // change value at that address
fmt.Println(x)    // 20 (x was modified!)
```

### Node.js vs Go Comparison

```javascript
// JavaScript - objects are always passed by reference
function modifyUser(user) {
    user.name = "Changed";  // Modifies original
}

const user = { name: "Original" };
modifyUser(user);
console.log(user.name);  // "Changed"
```

```go
// Go - you choose: value (copy) or pointer (reference)

// Value - gets a copy
func modifyUserValue(user User) {
    user.Name = "Changed"  // Only modifies the copy!
}

// Pointer - gets the original
func modifyUserPointer(user *User) {
    user.Name = "Changed"  // Modifies original
}

func main() {
    user := User{Name: "Original"}

    modifyUserValue(user)
    fmt.Println(user.Name)  // "Original" (unchanged!)

    modifyUserPointer(&user)  // Pass address
    fmt.Println(user.Name)  // "Changed"
}
```

---

## 15 Pointer Examples (5 Intermediate + 10 Hard)

### Intermediate Examples (1-5)

#### Example 1: Swap Two Numbers
```go
// Without pointers - DOESN'T WORK
func swapBroken(a, b int) {
    a, b = b, a  // Only swaps local copies
}

// With pointers - WORKS
func swap(a, b *int) {
    *a, *b = *b, *a  // Swap values at addresses
}

func main() {
    x, y := 10, 20
    swap(&x, &y)
    fmt.Println(x, y)  // 20, 10
}
```

#### Example 2: Modify Struct Field
```go
type Config struct {
    Debug bool
    Port  int
}

func enableDebug(c *Config) {
    c.Debug = true
    c.Port = 8080
}

func main() {
    config := Config{}
    enableDebug(&config)
    fmt.Printf("%+v\n", config)  // {Debug:true Port:8080}
}
```

#### Example 3: Return Pointer from Function
```go
func createUser(name string) *User {
    user := User{Name: name}  // Local variable
    return &user              // Return pointer (Go handles this safely!)
}

func main() {
    u := createUser("Tejas")
    fmt.Println(u.Name)  // "Tejas" - works because Go moves to heap
}
```

#### Example 4: nil Pointer Check
```go
func printUserName(u *User) {
    if u == nil {
        fmt.Println("No user provided")
        return
    }
    fmt.Println(u.Name)
}

func main() {
    var u *User  // nil by default
    printUserName(u)  // "No user provided"

    u = &User{Name: "Tejas"}
    printUserName(u)  // "Tejas"
}
```

#### Example 5: Pointer to Slice Element
```go
func doubleFirst(nums []int) {
    if len(nums) > 0 {
        ptr := &nums[0]  // Pointer to first element
        *ptr *= 2
    }
}

func main() {
    numbers := []int{5, 10, 15}
    doubleFirst(numbers)
    fmt.Println(numbers)  // [10 10 15]
}
```

### Hard Examples (6-15)

#### Example 6: Linked List with Pointers
```go
type Node struct {
    Value int
    Next  *Node  // Pointer to next node
}

func (n *Node) Append(value int) {
    current := n
    for current.Next != nil {
        current = current.Next
    }
    current.Next = &Node{Value: value}
}

func (n *Node) Print() {
    current := n
    for current != nil {
        fmt.Printf("%d -> ", current.Value)
        current = current.Next
    }
    fmt.Println("nil")
}

func main() {
    head := &Node{Value: 1}
    head.Append(2)
    head.Append(3)
    head.Print()  // 1 -> 2 -> 3 -> nil
}
```

#### Example 7: Double Pointer (Pointer to Pointer)
```go
func setToNil(pp **User) {
    *pp = nil  // Set the pointer itself to nil
}

func main() {
    u := &User{Name: "Tejas"}
    fmt.Println(u)  // &{Tejas}

    setToNil(&u)    // Pass address of pointer
    fmt.Println(u)  // <nil>
}
```

#### Example 8: Pointer in Map Value
```go
type Score struct {
    Points int
}

func main() {
    scores := map[string]*Score{
        "player1": {Points: 100},
        "player2": {Points: 200},
    }

    // Can modify through pointer
    scores["player1"].Points += 50
    fmt.Println(scores["player1"].Points)  // 150

    // Compare with value map (can't do this):
    // valueScores := map[string]Score{}
    // valueScores["p1"].Points++  // ERROR: cannot assign to valueScores["p1"].Points
}
```

#### Example 9: Interface with Pointer Receiver
```go
type Incrementer interface {
    Increment()
}

type Counter struct {
    count int
}

// Pointer receiver - only *Counter implements Incrementer
func (c *Counter) Increment() {
    c.count++
}

func doIncrement(i Incrementer) {
    i.Increment()
}

func main() {
    c := Counter{}

    // doIncrement(c)   // ERROR: Counter doesn't implement Incrementer
    doIncrement(&c)     // OK: *Counter implements Incrementer

    fmt.Println(c.count)  // 1
}
```

#### Example 10: Pointer to Embedded Struct
```go
type Address struct {
    City    string
    Country string
}

type Person struct {
    Name    string
    Address *Address  // Pointer to avoid copying
}

func main() {
    addr := &Address{City: "Mumbai", Country: "India"}

    p1 := Person{Name: "Tejas", Address: addr}
    p2 := Person{Name: "Other", Address: addr}  // Same address!

    p1.Address.City = "Delhi"
    fmt.Println(p2.Address.City)  // "Delhi" - p2 also changed!
}
```

#### Example 11: Slice of Pointers vs Pointer to Slice
```go
func main() {
    // Slice of pointers - each element is a pointer
    users := []*User{
        {Name: "User1"},
        {Name: "User2"},
    }

    // Modifying through slice element
    users[0].Name = "Modified"

    // Pointer to slice - the slice itself can be modified
    var slice []int
    appendToSlice(&slice, 1, 2, 3)
    fmt.Println(slice)  // [1 2 3]
}

func appendToSlice(s *[]int, vals ...int) {
    *s = append(*s, vals...)
}
```

#### Example 12: Recursive Tree with Pointers
```go
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

func (t *TreeNode) Insert(value int) {
    if value < t.Value {
        if t.Left == nil {
            t.Left = &TreeNode{Value: value}
        } else {
            t.Left.Insert(value)
        }
    } else {
        if t.Right == nil {
            t.Right = &TreeNode{Value: value}
        } else {
            t.Right.Insert(value)
        }
    }
}

func (t *TreeNode) InOrder() {
    if t == nil {
        return
    }
    t.Left.InOrder()
    fmt.Print(t.Value, " ")
    t.Right.InOrder()
}

func main() {
    root := &TreeNode{Value: 10}
    root.Insert(5)
    root.Insert(15)
    root.Insert(3)
    root.InOrder()  // 3 5 10 15
}
```

#### Example 13: Method That Returns Pointer to Self (Chaining)
```go
type Builder struct {
    query string
}

func (b *Builder) Select(fields string) *Builder {
    b.query += "SELECT " + fields
    return b  // Return pointer for chaining
}

func (b *Builder) From(table string) *Builder {
    b.query += " FROM " + table
    return b
}

func (b *Builder) Where(condition string) *Builder {
    b.query += " WHERE " + condition
    return b
}

func (b *Builder) Build() string {
    return b.query
}

func main() {
    query := (&Builder{}).
        Select("*").
        From("users").
        Where("active = true").
        Build()

    fmt.Println(query)  // SELECT * FROM users WHERE active = true
}
```

#### Example 14: Pointer Comparison
```go
func main() {
    a := 10
    b := 10

    p1 := &a
    p2 := &a
    p3 := &b

    fmt.Println(p1 == p2)  // true (same address)
    fmt.Println(p1 == p3)  // false (different addresses, even though same value)
    fmt.Println(*p1 == *p3) // true (values are equal)
}
```

#### Example 15: Pointer with Goroutines (Concurrent Modification)
```go
import "sync"

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func main() {
    counter := &SafeCounter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()  // Safe because of mutex
        }()
    }

    wg.Wait()
    fmt.Println(counter.Value())  // 1000 (always correct)
}
```

---

## Struct Tags

### What is `json:"email"`?

Struct tags are **metadata** attached to struct fields. They're used by packages (like `encoding/json`) to know how to process fields.

```go
type User struct {
    ID       int64  `json:"id"`           // JSON key will be "id"
    Name     string `json:"name"`         // JSON key will be "name"
    Email    string `json:"email"`        // JSON key will be "email"
    Password string `json:"-"`            // "-" means SKIP this field
    IsAdmin  bool   `json:"is_admin,omitempty"` // omitempty: skip if false/empty
}
```

### Why Can't We Write `type User struct { id int64 }`?

You CAN write lowercase `id`, but:
1. **It becomes private** - not accessible outside the package
2. **JSON encoder can't see it** - won't include in JSON output

```go
type User struct {
    id    int64  // Private - won't appear in JSON!
    Name  string // Public - will appear in JSON
}

func main() {
    u := User{id: 1, Name: "Tejas"}
    json, _ := json.Marshal(u)
    fmt.Println(string(json))  // {"Name":"Tejas"} - id is missing!
}
```

### Node.js vs Go Comparison

```typescript
// TypeScript - decorators/class-transformer
class User {
  @Expose({ name: 'user_id' })
  id: number;

  @Exclude()
  password: string;
}
```

```go
// Go - struct tags
type User struct {
    ID       int64  `json:"user_id"`
    Password string `json:"-"`
}
```

### Common Struct Tags

| Tag | Package | Purpose |
|-----|---------|---------|
| `json:"name"` | encoding/json | JSON field name |
| `json:"-"` | encoding/json | Skip this field |
| `json:",omitempty"` | encoding/json | Skip if empty/zero |
| `binding:"required"` | gin | Validation (required field) |
| `binding:"email"` | gin | Validation (must be email) |
| `db:"column_name"` | sqlx/gorm | Database column name |
| `gorm:"primaryKey"` | gorm | GORM ORM configuration |
| `validate:"min=3"` | go-validator | Custom validation |

### Multiple Tags

```go
type User struct {
    ID    int64  `json:"id" db:"user_id" gorm:"primaryKey"`
    Email string `json:"email" binding:"required,email" db:"email"`
}
```

---

## Number Types

### All Integer Types

| Type | Size | Range | Use Case |
|------|------|-------|----------|
| `int8` | 8 bits | -128 to 127 | Small numbers |
| `int16` | 16 bits | -32,768 to 32,767 | Small numbers |
| `int32` | 32 bits | -2.1B to 2.1B | Regular integers |
| `int64` | 64 bits | -9.2 quintillion to 9.2 quintillion | Large numbers, IDs |
| `int` | 32 or 64 bits | Platform dependent | **Default choice** |
| `uint8` | 8 bits | 0 to 255 | Unsigned (positive only) |
| `uint16` | 16 bits | 0 to 65,535 | Unsigned |
| `uint32` | 32 bits | 0 to 4.2B | Unsigned |
| `uint64` | 64 bits | 0 to 18.4 quintillion | Large positive numbers |
| `uint` | 32 or 64 bits | Platform dependent | Unsigned default |
| `byte` | 8 bits | Same as uint8 | Raw data, characters |
| `rune` | 32 bits | Same as int32 | Unicode characters |

### Float Types

| Type | Size | Precision |
|------|------|-----------|
| `float32` | 32 bits | ~7 decimal digits |
| `float64` | 64 bits | ~15 decimal digits (**default choice**) |

### When to Use What?

```go
// Default choices
var count int         // Most integers
var price float64     // Most decimals
var id int64          // Database IDs (to avoid overflow)

// Specific cases
var age uint8         // Age (0-255, never negative)
var char rune         // Single Unicode character
var data []byte       // Raw binary data

// Database IDs
type User struct {
    ID int64  // Use int64 for database IDs (handles large values)
}
```

### Node.js Comparison

```javascript
// JavaScript - only one number type
let x = 42;        // Number (64-bit float internally)
let big = 9007199254740993n;  // BigInt for huge numbers
```

```go
// Go - explicit types
var x int = 42           // 32/64-bit integer
var y int64 = 9007199254740993  // 64-bit integer
var z float64 = 3.14     // 64-bit float
```

---

## Validation Libraries

### Gin's Built-in Validation (Uses go-playground/validator)

Already in your code! Look at [models/user.go](../06_gin_server/models/user.go):

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### Common Validation Tags

| Tag | Description |
|-----|-------------|
| `required` | Field must not be empty |
| `email` | Must be valid email |
| `min=3` | Minimum length/value |
| `max=50` | Maximum length/value |
| `len=10` | Exact length |
| `oneof=male female` | Must be one of these values |
| `numeric` | Must be numeric string |
| `alphanum` | Alphanumeric only |
| `url` | Must be valid URL |
| `uuid` | Must be valid UUID |
| `gt=0` | Greater than 0 |
| `gte=0` | Greater than or equal to 0 |

### Example with More Validation

```go
type CreateProductRequest struct {
    Name        string  `json:"name" binding:"required,min=3,max=100"`
    Price       float64 `json:"price" binding:"required,gt=0"`
    Category    string  `json:"category" binding:"required,oneof=electronics clothing food"`
    SKU         string  `json:"sku" binding:"required,len=10,alphanum"`
    Description string  `json:"description" binding:"max=500"`
    ImageURL    string  `json:"image_url" binding:"omitempty,url"`
    Stock       int     `json:"stock" binding:"gte=0"`
}
```

### Custom Validation

```go
import "github.com/go-playground/validator/v10"

// Custom validator function
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // No spaces allowed
    return !strings.Contains(username, " ")
}

func main() {
    router := gin.Default()

    // Register custom validation
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("nospace", validateUsername)
    }
}

// Use it
type User struct {
    Username string `binding:"required,nospace"`
}
```

### Node.js Zod vs Go Comparison

```typescript
// Zod (TypeScript)
const UserSchema = z.object({
  username: z.string().min(3).max(50),
  email: z.string().email(),
  password: z.string().min(6),
});
```

```go
// Go struct tags
type UserRequest struct {
    Username string `binding:"required,min=3,max=50"`
    Email    string `binding:"required,email"`
    Password string `binding:"required,min=6"`
}
```

---

## Logging in Go

### Standard Library (log)

```go
import "log"

func main() {
    log.Println("Info message")
    log.Printf("User %s logged in", "tejas")
    log.Fatal("Fatal error - will exit program")  // Calls os.Exit(1)
    log.Panic("Panic - will panic")               // Calls panic()
}
```

### Popular Logging Libraries

| Library | Description | Best For |
|---------|-------------|----------|
| **zerolog** | Zero allocation JSON logger | Production, high performance |
| **zap** | Uber's structured logger | Production, high performance |
| **logrus** | Structured logger (older) | Simple projects |
| **slog** | Go 1.21+ standard library | New projects |

### zerolog Example (Recommended)

```bash
go get github.com/rs/zerolog
```

```go
import (
    "os"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Pretty console output for development
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Basic logging
    log.Info().Msg("Server starting")
    log.Error().Err(err).Msg("Failed to connect")

    // Structured logging with fields
    log.Info().
        Str("user", "tejas").
        Int("status", 200).
        Msg("Request completed")

    // Output: {"level":"info","user":"tejas","status":200,"message":"Request completed"}
}
```

### slog (Go 1.21+ Standard Library)

```go
import "log/slog"

func main() {
    // JSON output
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    logger.Info("User logged in",
        slog.String("username", "tejas"),
        slog.Int("user_id", 123),
    )

    // Output: {"time":"...","level":"INFO","msg":"User logged in","username":"tejas","user_id":123}
}
```

### Gin with Custom Logger

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Setup zerolog
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    router := gin.New()

    // Custom logging middleware
    router.Use(func(c *gin.Context) {
        start := time.Now()
        c.Next()

        log.Info().
            Str("method", c.Request.Method).
            Str("path", c.Request.URL.Path).
            Int("status", c.Writer.Status()).
            Dur("latency", time.Since(start)).
            Msg("Request")
    })
}
```

---

## gin.New() vs gin.Default()

### The Problem with gin.Default()

```go
// What most tutorials show
r := gin.Default()
```

**What `gin.Default()` does internally:**
```go
func Default() *Engine {
    engine := New()
    engine.Use(Logger(), Recovery())  // Adds these automatically
    return engine
}
```

### Why Avoid gin.Default()?

| Issue | Problem |
|-------|---------|
| **Logger output** | Uses default logger that outputs to stdout with its own format |
| **No control** | Can't customize log format, level, or destination |
| **Duplicate logging** | If you add your own logger, you get double logs |
| **Production concerns** | Default logger isn't suitable for production (no JSON, no levels) |

### The Better Approach: gin.New()

```go
// Create router WITHOUT default middleware
router := gin.New()

// Add only what you need
router.Use(gin.Recovery())              // Keep this - prevents crashes
router.Use(middlewares.LoggerMiddleware()) // Your custom logger
router.Use(middlewares.CORSMiddleware())   // Your CORS config
```

### Comparison

```go
// BAD - gin.Default()
func main() {
    r := gin.Default()  // Has Logger() + Recovery()

    // Problem: If you add custom logger, you get DOUBLE logging
    r.Use(customLogger())  // Now you have 2 loggers!
}

// GOOD - gin.New()
func main() {
    r := gin.New()  // Empty - no middleware

    // Add exactly what you need
    r.Use(gin.Recovery())     // Prevent panic crashes
    r.Use(customLogger())     // Your logger only
    r.Use(corsMiddleware())   // Your CORS
}
```

### What Each Middleware Does

| Middleware | Purpose | Keep it? |
|------------|---------|----------|
| `gin.Recovery()` | Catches panics, prevents server crash | **Yes, always** |
| `gin.Logger()` | Logs requests to stdout | **Replace with custom** |

### Production Setup

```go
func NewServer() *gin.Engine {
    // Set mode based on environment
    if os.Getenv("ENV") == "production" {
        gin.SetMode(gin.ReleaseMode)  // Disables debug logs
    }

    router := gin.New()

    // Recovery is essential - catches panics
    router.Use(gin.Recovery())

    // Custom structured logger (JSON for production)
    router.Use(func(c *gin.Context) {
        start := time.Now()
        c.Next()

        // Structured logging (works with log aggregators)
        log.Info().
            Str("method", c.Request.Method).
            Str("path", c.Request.URL.Path).
            Int("status", c.Writer.Status()).
            Dur("latency", time.Since(start)).
            Str("ip", c.ClientIP()).
            Msg("request")
    })

    return router
}
```

### Node.js Comparison

```javascript
// Express.js - same concept
const app = express();

// BAD - using morgan default
app.use(morgan('dev'));       // Default format
app.use(customLogger());      // Now double logging!

// GOOD - only your logger
app.use(customLogger());      // Single, controlled logging
```

### Quick Reference

| Method | Use Case |
|--------|----------|
| `gin.Default()` | Quick prototypes, tutorials |
| `gin.New()` | Production apps, when you need control |

### Our Server Uses gin.New()

See [server/server.go](../06_gin_server/server/server.go):
```go
router := gin.New()                         // No default middleware
router.Use(gin.Recovery())                  // Panic recovery
router.Use(middlewares.LoggerMiddleware())  // Custom logger
router.Use(middlewares.CORSMiddleware())    // Custom CORS
```

---

## JWT RegisteredClaims

### What is `jwt.RegisteredClaims`?

This is **embedding** - Go's way of "inheriting" fields and methods.

```go
type JWTClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims  // ← Embedded! Gets all RegisteredClaims fields
}
```

### What's Inside RegisteredClaims?

```go
// From github.com/golang-jwt/jwt/v5
type RegisteredClaims struct {
    Issuer    string           `json:"iss,omitempty"` // Who created the token
    Subject   string           `json:"sub,omitempty"` // Who the token is about
    Audience  ClaimStrings     `json:"aud,omitempty"` // Who can use the token
    ExpiresAt *NumericDate     `json:"exp,omitempty"` // When token expires
    NotBefore *NumericDate     `json:"nbf,omitempty"` // Token not valid before
    IssuedAt  *NumericDate     `json:"iat,omitempty"` // When token was issued
    ID        string           `json:"jti,omitempty"` // Unique token ID
}
```

### How Embedding Works

```go
type JWTClaims struct {
    UserID   string
    Username string
    jwt.RegisteredClaims  // Embedded
}

func main() {
    claims := JWTClaims{
        UserID:   "123",
        Username: "tejas",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "my-app",
        },
    }

    // Access embedded fields directly
    fmt.Println(claims.ExpiresAt)  // Works! (promoted from RegisteredClaims)
    fmt.Println(claims.Issuer)     // Works!

    // Or through the embedded struct
    fmt.Println(claims.RegisteredClaims.ExpiresAt)
}
```

### Node.js Comparison

```javascript
// JavaScript - spread operator
const baseClaims = {
  iss: "my-app",
  exp: Date.now() + 86400000,
  iat: Date.now()
};

const token = {
  user_id: "123",
  username: "tejas",
  ...baseClaims  // Spread base claims
};
```

```go
// Go - embedding
type JWTClaims struct {
    UserID   string
    Username string
    jwt.RegisteredClaims  // Like spreading baseClaims
}
```

---

## Nested Packages & Imports

### Question: What if we have nested folders?

Yes! You name the **parent folder** in the import path.

### Folder Structure

```
gin_server/
├── go.mod                    ← module gin_server
├── main.go
├── server/
│   └── server.go             ← package server
├── handlers/
│   ├── auth.go               ← package handlers
│   └── user.go               ← package handlers
├── utils/
│   ├── jwt.go                ← package utils
│   ├── hash.go               ← package utils
│   └── validators/           ← NESTED!
│       └── email.go          ← package validators
└── internal/
    └── database/
        └── db.go             ← package database
```

### Importing Nested Packages

```go
// main.go
package main

import (
    "gin_server/handlers"                    // handlers folder
    "gin_server/utils"                       // utils folder
    "gin_server/utils/validators"            // nested validators folder
    "gin_server/internal/database"           // nested database folder
)

func main() {
    handlers.RegisterHandler()
    utils.HashPassword("secret")
    validators.ValidateEmail("test@example.com")
    database.Connect()
}
```

### The Rule

```go
// File: utils/validators/email.go
package validators  // ← Package name = FOLDER name (not full path)

func ValidateEmail(email string) bool {
    // ...
}
```

```go
// Importing it
import "gin_server/utils/validators"  // Full path from module root

// Using it
validators.ValidateEmail("test@test.com")  // Use package name
```

### Real Example

```go
// File: utils/validators/email.go
package validators

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
    return emailRegex.MatchString(email)
}

// File: handlers/auth.go
package handlers

import (
    "gin_server/utils/validators"  // Import nested package
)

func RegisterHandler(c *gin.Context) {
    email := c.PostForm("email")

    if !validators.IsValidEmail(email) {  // Use it
        // ...
    }
}
```

### Import Alias (When Names Conflict)

```go
import (
    "gin_server/utils/validators"
    customValidators "gin_server/custom/validators"  // Alias to avoid conflict
)

func main() {
    validators.IsValidEmail("...")
    customValidators.IsValidEmail("...")
}
```

---

## Path Aliases

### Can we do `~/` or `@/` like in Node.js?

**No**, Go doesn't support path aliases like TypeScript/Webpack.

```typescript
// TypeScript tsconfig.json
{
  "paths": {
    "@/*": ["./src/*"],
    "~/*": ["./src/*"]
  }
}

// Then use
import { User } from "@/models/user";
```

**Go doesn't have this.** You always use the full module path:

```go
import "gin_server/models"
import "gin_server/utils/validators"
```

### Why?

Go prioritizes **explicit over implicit**. The full path tells you exactly where the code comes from.

### Workarounds

1. **Keep module names short**
   ```bash
   go mod init myapp  # Short name
   ```
   ```go
   import "myapp/handlers"  // Not too long
   ```

2. **Use import aliases for long paths**
   ```go
   import (
       auth "github.com/mycompany/myproject/internal/services/authentication"
   )

   func main() {
       auth.Login()
   }
   ```

3. **Restructure for shorter paths**
   ```
   # Instead of deeply nested
   internal/services/auth/handlers/login.go

   # Flatten
   handlers/auth.go
   ```

### Summary Table

| Feature | Node.js/TS | Go |
|---------|------------|-----|
| Path aliases (`@/`, `~/`) | ✅ Yes | ❌ No |
| Relative imports (`./`) | ✅ Yes | ❌ No (use full module path) |
| Import aliases | ✅ Yes | ✅ Yes (`import x "path"`) |
| Full path imports | ✅ Yes | ✅ Yes (required) |

---

## Quick Reference

| Concept | Node.js/TS | Go |
|---------|------------|-----|
| Class | `class User {}` | `type User struct {}` |
| Constructor | `constructor()` | `func NewUser() *User` |
| Method | `user.greet()` | `func (u User) Greet()` |
| Inheritance | `extends` | Embedding |
| Interface | `implements` | Implicit (just match methods) |
| Private | `private` or `#field` | lowercase first letter |
| Public | `public` | Uppercase first letter |
| JSON mapping | Decorators | Struct tags |
| Validation | Zod/Joi | go-playground/validator |
| Logging | Winston/Pino | zerolog/zap/slog |
| Path aliases | `@/`, `~/` | Not supported |
