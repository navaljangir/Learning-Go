# Interfaces in Go

## Definition

**An interface in Go is a type that specifies a set of method signatures (behavior), but doesn't implement them.**

Any type that implements all the methods in an interface automatically satisfies that interface - **no explicit declaration needed** (implicit implementation).

---

## Simple Analogy

**From Node.js/TypeScript:**
```typescript
// TypeScript: Explicit interface
interface Writer {
  write(data: string): void;
}

// Must explicitly implement
class FileWriter implements Writer {
  write(data: string) { /* ... */ }
}
```

**In Go:**
```go
// Go: Just define the contract
type Writer interface {
    Write(data []byte) (int, error)
}

// Any type with this method automatically satisfies Writer
type FileWriter struct{}

func (f FileWriter) Write(data []byte) (int, error) {
    // Implementation
    return len(data), nil
}
// FileWriter is now a Writer - no "implements" keyword needed!
```

---

## Basic Syntax

```go
type InterfaceName interface {
    MethodName1(param type) returnType
    MethodName2(param type) (returnType1, returnType2)
}
```

**Key points:**
- Interfaces only declare method signatures, not implementations
- Interface names often end in `-er` (e.g., `Reader`, `Writer`, `Stringer`)
- Methods must match **exactly** (name, parameters, return types)

---

## Simple Example

```go
package main

import "fmt"

// 1. Define the interface
type Greeter interface {
    Greet() string
}

// 2. Create types that implement it
type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Hello, I'm " + p.Name
}

type Robot struct {
    ID int
}

func (r Robot) Greet() string {
    return fmt.Sprintf("Beep boop, Robot #%d", r.ID)
}

// 3. Use the interface as a parameter type
func SayHello(g Greeter) {
    fmt.Println(g.Greet())
}

func main() {
    p := Person{Name: "Alice"}
    r := Robot{ID: 42}

    SayHello(p)  // Hello, I'm Alice
    SayHello(r)  // Beep boop, Robot #42
}
```

**Why this matters:** `SayHello()` accepts ANY type that has a `Greet()` method. No inheritance needed!

---

## Empty Interface `interface{}`

The empty interface has **zero methods**, so every type satisfies it.

```go
func PrintAnything(v interface{}) {
    fmt.Println(v)
}

PrintAnything(42)
PrintAnything("hello")
PrintAnything([]int{1, 2, 3})
```

**Modern Go (1.18+):** Use `any` instead:
```go
func PrintAnything(v any) {  // any is an alias for interface{}
    fmt.Println(v)
}
```

---

## Common Standard Library Interfaces

### 1. `io.Reader` (reading data)
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Examples:** Files, network connections, strings all implement `Reader`.

```go
import (
    "io"
    "os"
    "strings"
)

func processData(r io.Reader) {
    data, _ := io.ReadAll(r)
    fmt.Println(string(data))
}

func main() {
    file, _ := os.Open("data.txt")
    processData(file)  // Works with file

    str := strings.NewReader("Hello")
    processData(str)  // Works with string
}
```

### 2. `io.Writer` (writing data)
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

```go
func saveData(w io.Writer, data string) {
    w.Write([]byte(data))
}

func main() {
    file, _ := os.Create("output.txt")
    saveData(file, "Hello")  // Writes to file

    saveData(os.Stdout, "Hello")  // Writes to console
}
```

### 3. `fmt.Stringer` (string representation)
```go
type Stringer interface {
    String() string
}
```

If a type implements `String()`, `fmt.Println()` will use it automatically.

```go
type User struct {
    Name  string
    Email string
}

func (u User) String() string {
    return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

func main() {
    u := User{"Alice", "alice@example.com"}
    fmt.Println(u)  // Output: Alice <alice@example.com>
}
```

### 4. `error` interface
```go
type error interface {
    Error() string
}
```

Any type with an `Error()` method is an error!

```go
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func doSomething() error {
    return MyError{Code: 404, Message: "Not found"}
}
```

---

## Why Use Interfaces?

### 1. **Polymorphism**
Different types can be used interchangeably:

```go
type Shape interface {
    Area() float64
}

type Circle struct { Radius float64 }
func (c Circle) Area() float64 { return 3.14 * c.Radius * c.Radius }

type Rectangle struct { Width, Height float64 }
func (r Rectangle) Area() float64 { return r.Width * r.Height }

func printArea(s Shape) {
    fmt.Println("Area:", s.Area())
}

func main() {
    printArea(Circle{Radius: 5})
    printArea(Rectangle{Width: 3, Height: 4})
}
```

### 2. **Decoupling / Dependency Injection**
Depend on behavior, not concrete types:

```go
// Bad: Tightly coupled to PostgreSQL
func GetUser(db *PostgresDB, id int) User {
    return db.QueryUser(id)
}

// Good: Can use any database that implements the interface
type Database interface {
    QueryUser(id int) User
}

func GetUser(db Database, id int) User {
    return db.QueryUser(id)
}

// Now you can swap PostgreSQL, MySQL, MockDB, etc.
```

### 3. **Testing**
Easy to create mocks:

```go
type EmailSender interface {
    Send(to, subject, body string) error
}

type RealEmailSender struct{}
func (r RealEmailSender) Send(to, subject, body string) error {
    // Actually send email
}

type MockEmailSender struct{}
func (m MockEmailSender) Send(to, subject, body string) error {
    fmt.Println("Mock: Email sent to", to)
    return nil
}

func NotifyUser(sender EmailSender, user string) {
    sender.Send(user, "Hello", "Welcome!")
}

// In production: use RealEmailSender
// In tests: use MockEmailSender
```

---

## Interface Composition

Interfaces can be combined:

```go
type Reader interface {
    Read(p []byte) (int, error)
}

type Writer interface {
    Write(p []byte) (int, error)
}

type Closer interface {
    Close() error
}

// ReadWriter combines Reader and Writer
type ReadWriter interface {
    Reader
    Writer
}

// ReadWriteCloser combines all three
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

**Example:** `os.File` implements `ReadWriteCloser` because it has `Read()`, `Write()`, and `Close()` methods.

---

## Type Assertion

Check if an interface value holds a specific type:

```go
var i interface{} = "hello"

// Type assertion
s := i.(string)
fmt.Println(s)  // hello

// Safe type assertion (with ok check)
s, ok := i.(string)
if ok {
    fmt.Println("It's a string:", s)
}

// Panic if wrong type
n := i.(int)  // panic: interface conversion: interface {} is string, not int
```

---

## Type Switch

Handle multiple types:

```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Println("Integer:", v)
    case string:
        fmt.Println("String:", v)
    case bool:
        fmt.Println("Boolean:", v)
    default:
        fmt.Println("Unknown type")
    }
}

describe(42)       // Integer: 42
describe("hello")  // String: hello
describe(true)     // Boolean: true
```

---

## Best Practices

1. **Keep interfaces small** - Prefer single-method interfaces (e.g., `io.Reader`, `io.Writer`)
2. **Define interfaces where they're used** - Not where types are defined
3. **Accept interfaces, return structs** - Functions should accept interfaces (flexible) but return concrete types (clear)
4. **Name with `-er` suffix** - `Reader`, `Writer`, `Handler`, `Validator`

```go
// Good: Small, focused interface
type Validator interface {
    Validate() error
}

// Less ideal: Large, kitchen-sink interface
type UserService interface {
    Create(user User) error
    Update(user User) error
    Delete(id int) error
    FindByID(id int) (User, error)
    FindByEmail(email string) (User, error)
    List() ([]User, error)
}
```

---

## Real-World Example: HTTP Handler

```go
// The http.Handler interface
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Any type with this method can handle HTTP requests
type MyHandler struct{}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func main() {
    http.ListenAndServe(":8080", MyHandler{})
}
```

---

## Summary

| Concept | Description |
|---------|-------------|
| **Interface** | Contract specifying methods a type must have |
| **Implicit** | No `implements` keyword - automatic satisfaction |
| **Empty interface** | `interface{}` or `any` - accepts any type |
| **Type assertion** | Check/extract concrete type from interface |
| **Composition** | Combine interfaces to create larger ones |
| **Polymorphism** | Different types used interchangeably via shared interface |

**Golden Rule:** Design your code to accept interfaces (flexible) and return concrete types (clear).
