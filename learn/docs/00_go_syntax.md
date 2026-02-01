# Go Syntax Reference

A quick reference for Go syntax, organized by topic.

---

## Table of Contents

### Basics
1. [Packages & Imports](#packages--imports)
2. [Variables & Types](#variables--types)
3. [Printing & Formatting](#printing--formatting)
4. [Control Flow](#control-flow)
5. [Functions](#functions)
6. [Pointers](#pointers)
7. [make vs new](#make-vs-new)

### Data Structures
8. [Arrays, Slices & Maps](#arrays-slices--maps)
9. [Strings](#strings)
10. [Structs & Methods](#structs--methods)
11. [Interfaces](#interfaces)

### Concurrency
12. [Goroutines](#goroutines)
13. [Channels](#channels)
14. [sync Package](#sync-package)
15. [Context](#context)

### Standard Library
16. [Time & Timestamps](#time--timestamps)
17. [File I/O](#file-io)
18. [HTTP Client & Server](#http-client--server)
19. [JSON & Encoding](#json--encoding)
20. [Strings Package](#strings-package)
21. [Regular Expressions](#regular-expressions)
22. [Sorting](#sorting)
23. [Random Numbers](#random-numbers)
24. [Logging](#logging)

### Development
25. [Error Handling](#error-handling)
26. [Testing](#testing)
27. [Go Modules & Commands](#go-modules--commands)
28. [Generics](#generics)
29. [Command Line Args](#command-line-args)
30. [Environment Variables](#environment-variables)
31. [Common Patterns](#common-patterns)

---

## Packages & Imports

### Package Declaration

Every Go file must start with a package declaration.

```go
package main        // Executable program (has main function)
package mypackage   // Library package (imported by others)
```

### Import Statements

```go
// Single import
import "fmt"

// Multiple imports (grouped)
import (
    "fmt"
    "strings"
    "time"
)

// Import with alias
import (
    f "fmt"                           // Use as f.Println()
    "math/rand"                       // Last part is package name
    _ "github.com/lib/pq"             // Blank import (side effects only)
    . "fmt"                           // Dot import (Println instead of fmt.Println) - avoid!
)

// Using aliased import
f.Println("Hello")
```

### Package Visibility (Exported vs Unexported)

```go
package mypackage

// EXPORTED (uppercase first letter) - accessible from other packages
var PublicVar = "visible"
func PublicFunc() {}
type PublicStruct struct {
    PublicField  string   // Exported field
    privateField string   // Unexported field (lowercase)
}

// UNEXPORTED (lowercase first letter) - only within this package
var privateVar = "hidden"
func privateFunc() {}
type privateStruct struct {}
```

### init() Function

`init()` runs automatically when package is imported, before `main()`.

```go
package main

import "fmt"

var config string

func init() {
    // Runs before main()
    // Common uses: setup, config loading, registering drivers
    config = "initialized"
    fmt.Println("init() called")
}

func main() {
    fmt.Println("main() called")
    fmt.Println(config)
}

// Output:
// init() called
// main() called
// initialized
```

**Multiple init functions:**
```go
func init() { fmt.Println("first init") }
func init() { fmt.Println("second init") }  // Both run, in order
```

### Package Organization

```
myproject/
├── go.mod              # Module definition
├── main.go             # package main
├── utils/
│   └── helpers.go      # package utils
└── models/
    └── user.go         # package models
```

```go
// main.go
package main

import (
    "myproject/utils"    // Import by module path
    "myproject/models"
)

func main() {
    utils.DoSomething()
    u := models.User{}
}
```

### Internal Packages

```
myproject/
├── internal/           # Only importable within myproject
│   └── secret.go
└── cmd/
    └── myapp/
        └── main.go     # Can import internal/
```

Packages in `internal/` can only be imported by code in the parent directory tree.

---

## Variables & Types

### Declaration

```go
// Explicit type
var name string = "Go"
var age int = 10

// Type inference
var name = "Go"       // string inferred
age := 10             // short declaration (inside functions only)

// Multiple variables
var x, y int = 1, 2
a, b := "hello", 42

// Zero values (default)
var i int      // 0
var f float64  // 0.0
var b bool     // false
var s string   // "" (empty string)
var p *int     // nil
```

### Basic Types

```go
// Numbers
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
complex64, complex128

// Other
bool        // true, false
string      // "hello"
byte        // alias for uint8
rune        // alias for int32 (Unicode code point)
```

### Type Conversion

```go
i := 42
f := float64(i)    // int to float
s := string(i)     // NOT what you expect! Use strconv

// Proper conversions
import "strconv"
s := strconv.Itoa(42)           // int to string: "42"
i, err := strconv.Atoi("42")    // string to int: 42
f, err := strconv.ParseFloat("3.14", 64)  // string to float
```

### Constants

```go
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)

// iota - auto-incrementing
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
)
```

---

## Printing & Formatting

### Basic Print Functions

```go
import "fmt"

fmt.Print("no newline")
fmt.Println("with newline")
fmt.Printf("formatted: %s\n", "hello")

// Print to string (not stdout)
s := fmt.Sprintf("Hello %s", "World")  // returns string
```

### Format Verbs

```go
// General
%v      // default format
%+v     // with field names (structs)
%#v     // Go syntax representation
%T      // type of value
%%      // literal percent sign

// Strings
%s      // string
%q      // quoted string "hello"
%x      // hex encoding

// Numbers
%d      // decimal integer
%b      // binary
%o      // octal
%x, %X  // hexadecimal (lowercase/uppercase)
%f      // float (default precision)
%.2f    // float with 2 decimal places
%e      // scientific notation
%9d     // width 9, right-aligned
%-9d    // width 9, left-aligned
%09d    // width 9, zero-padded

// Boolean
%t      // true or false

// Pointer
%p      // pointer address
```

### Examples

```go
name := "Go"
version := 1.21
count := 42

fmt.Printf("Language: %s\n", name)           // Language: Go
fmt.Printf("Version: %.2f\n", version)       // Version: 1.21
fmt.Printf("Count: %d\n", count)             // Count: 42
fmt.Printf("Binary: %b\n", count)            // Binary: 101010
fmt.Printf("Hex: %x\n", count)               // Hex: 2a
fmt.Printf("Padded: %05d\n", count)          // Padded: 00042
fmt.Printf("Type: %T\n", count)              // Type: int

// Struct printing
type User struct {
    Name string
    Age  int
}
u := User{"Alice", 30}
fmt.Printf("%v\n", u)    // {Alice 30}
fmt.Printf("%+v\n", u)   // {Name:Alice Age:30}
fmt.Printf("%#v\n", u)   // main.User{Name:"Alice", Age:30}
```

---

## Time & Timestamps

### Import

```go
import "time"
```

### Current Time

```go
now := time.Now()                    // current local time
utc := time.Now().UTC()              // current UTC time
unix := time.Now().Unix()            // Unix timestamp (seconds)
unixMilli := time.Now().UnixMilli()  // Unix timestamp (milliseconds)
unixNano := time.Now().UnixNano()    // Unix timestamp (nanoseconds)
```

### Time Formatting

Go uses a **reference time** for formatting: `Mon Jan 2 15:04:05 MST 2006`

```go
now := time.Now()

// Common formats
now.Format("2006-01-02")                    // 2024-01-15
now.Format("2006-01-02 15:04:05")           // 2024-01-15 14:30:45
now.Format("15:04:05")                      // 14:30:45
now.Format("15:04:05.000")                  // 14:30:45.123 (with milliseconds)
now.Format("15:04:05.000000")               // 14:30:45.123456 (with microseconds)
now.Format("Mon, 02 Jan 2006")              // Mon, 15 Jan 2024
now.Format("January 2, 2006")               // January 15, 2024
now.Format(time.RFC3339)                    // 2024-01-15T14:30:45Z
now.Format(time.Kitchen)                    // 2:30PM

// Reference time breakdown:
// Mon    = day of week
// Jan    = month (short)
// 2      = day
// 15     = hour (24h)
// 04     = minute
// 05     = second
// 2006   = year
// MST    = timezone
// -0700  = timezone offset
```

### Logging with Timestamps

```go
// Simple timestamp in logs
fmt.Printf("[%s] Something happened\n", time.Now().Format("15:04:05.000"))
// Output: [14:30:45.123] Something happened

// With date
fmt.Printf("[%s] Starting server\n", time.Now().Format("2006-01-02 15:04:05"))
// Output: [2024-01-15 14:30:45] Starting server

// For debugging concurrent code
func logWithTime(msg string) {
    fmt.Printf("[%s] %s\n", time.Now().Format("05.000"), msg)
}

// Usage in goroutines
go func() {
    logWithTime("Goroutine 1: started")
    // ... work ...
    logWithTime("Goroutine 1: done")
}()
```

### Durations

```go
// Creating durations
d := 5 * time.Second
d := 100 * time.Millisecond
d := 2 * time.Minute
d := time.Duration(500) * time.Millisecond

// Duration constants
time.Nanosecond
time.Microsecond
time.Millisecond
time.Second
time.Minute
time.Hour

// Sleep
time.Sleep(2 * time.Second)

// Measure elapsed time
start := time.Now()
// ... do work ...
elapsed := time.Since(start)
fmt.Printf("Took: %v\n", elapsed)           // Took: 1.234567s
fmt.Printf("Took: %s\n", elapsed)           // Took: 1.234567s
fmt.Printf("Took: %d ms\n", elapsed.Milliseconds())  // Took: 1234 ms
```

### Parsing Time

```go
// Parse string to time
t, err := time.Parse("2006-01-02", "2024-01-15")
t, err := time.Parse("2006-01-02 15:04:05", "2024-01-15 14:30:45")
t, err := time.Parse(time.RFC3339, "2024-01-15T14:30:45Z")

// Parse with timezone
loc, _ := time.LoadLocation("America/New_York")
t, err := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-15 14:30:45", loc)
```

### Time Operations

```go
now := time.Now()

// Add/subtract duration
future := now.Add(24 * time.Hour)        // tomorrow
past := now.Add(-1 * time.Hour)          // 1 hour ago

// Compare times
now.Before(future)  // true
now.After(past)     // true
now.Equal(now)      // true

// Difference between times
diff := future.Sub(now)  // returns Duration
fmt.Println(diff)        // 24h0m0s

// Extract components
now.Year()      // 2024
now.Month()     // January
now.Day()       // 15
now.Hour()      // 14
now.Minute()    // 30
now.Second()    // 45
now.Weekday()   // Monday
```

### Timers and Tickers

```go
// One-shot timer
timer := time.NewTimer(2 * time.Second)
<-timer.C  // blocks for 2 seconds
fmt.Println("Timer fired!")

// Or simpler:
<-time.After(2 * time.Second)

// Repeating ticker
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop()

for i := 0; i < 5; i++ {
    <-ticker.C
    fmt.Println("Tick")
}
```

### time.After() - Timeout Channel

`time.After(d)` returns a **channel** that sends a value after duration `d`.

```go
// What time.After returns:
time.After(2 * time.Second)  // Returns: <-chan time.Time
                              // Sends current time after 2 seconds
```

**Simple delay (like Sleep):**
```go
fmt.Println("Starting...")
<-time.After(2 * time.Second)  // Blocks for 2 seconds
fmt.Println("Done!")

// Same as:
time.Sleep(2 * time.Second)
```

**Timeout with select (MOST COMMON USE):**
```go
select {
case result := <-apiChannel:
    fmt.Println("Got result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")  // Fires if apiChannel doesn't respond in 5s
}
```

**Visual - How it works:**
```
select starts
    │
    ├── case result := <-apiChannel   (waiting for data...)
    │
    ├── case <-time.After(5s)         (timer: 0s...1s...2s...3s...4s...5s FIRE!)
    │                                                                  ↓
    │                                                        This case wins if
    │                                                        apiChannel is slow
    └── First one ready wins!
```

**time.After vs time.Sleep:**

| | `time.Sleep` | `time.After` |
|--|--------------|--------------|
| Blocks | Yes | Yes (when you receive) |
| Works with select | No | Yes |
| Returns | Nothing | Channel |
| Use case | Simple delay | Timeouts, racing |

```go
// time.Sleep - can't race with other operations
time.Sleep(5 * time.Second)  // Just waits, no escape

// time.After - can race with other channels
select {
case data := <-ch:
    // Got data before timeout
case <-time.After(5 * time.Second):
    // Timeout - ch was too slow
}
```

### time.NewTimer vs time.After

```go
// time.After - simple, but can't cancel
<-time.After(5 * time.Second)

// time.NewTimer - can be stopped/reset
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()  // Clean up if not needed

select {
case <-timer.C:
    fmt.Println("Timer fired")
case <-done:
    timer.Stop()  // Cancel the timer
    fmt.Println("Cancelled")
}

// Reset timer for reuse
timer.Reset(3 * time.Second)
```

### Complete Time Quick Reference

| Operation | Code |
|-----------|------|
| Current time | `time.Now()` |
| Unix timestamp | `time.Now().Unix()` |
| Unix nanoseconds | `time.Now().UnixNano()` |
| Format time | `t.Format("2006-01-02 15:04:05")` |
| Sleep | `time.Sleep(2 * time.Second)` |
| Measure elapsed | `time.Since(start)` |
| Delay channel | `<-time.After(2 * time.Second)` |
| Timeout in select | `case <-time.After(5 * time.Second):` |
| Repeating ticker | `ticker := time.NewTicker(1 * time.Second)` |
| Stoppable timer | `timer := time.NewTimer(5 * time.Second)` |
| Add duration | `t.Add(24 * time.Hour)` |
| Time difference | `t2.Sub(t1)` |

---

## Control Flow

### If/Else

```go
if x > 0 {
    fmt.Println("positive")
} else if x < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}

// With initialization statement
if err := doSomething(); err != nil {
    fmt.Println("Error:", err)
}

// err is only available inside the if block
```

### For Loop

```go
// Classic for
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style
for x < 100 {
    x *= 2
}

// Infinite loop
for {
    // break to exit
}

// Range over slice
nums := []int{1, 2, 3}
for index, value := range nums {
    fmt.Printf("%d: %d\n", index, value)
}

// Range - index only
for i := range nums {
    fmt.Println(i)
}

// Range - value only
for _, v := range nums {
    fmt.Println(v)
}

// Range over map
m := map[string]int{"a": 1, "b": 2}
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}

// Range over string (runes)
for i, r := range "hello" {
    fmt.Printf("%d: %c\n", i, r)
}

// Range over channel
for msg := range ch {
    fmt.Println(msg)
}
```

### Switch

```go
// Basic switch
switch day {
case "Mon":
    fmt.Println("Monday")
case "Tue", "Wed":  // multiple values
    fmt.Println("Midweek")
default:
    fmt.Println("Other day")
}

// No condition (like if-else chain)
switch {
case x < 0:
    fmt.Println("negative")
case x > 0:
    fmt.Println("positive")
default:
    fmt.Println("zero")
}

// Type switch
switch v := i.(type) {
case int:
    fmt.Printf("int: %d\n", v)
case string:
    fmt.Printf("string: %s\n", v)
default:
    fmt.Printf("unknown type: %T\n", v)
}

// Fallthrough (rarely used)
switch n {
case 1:
    fmt.Println("one")
    fallthrough  // continues to next case
case 2:
    fmt.Println("one or two")
}
```

### Defer

```go
// Executes when function returns
func example() {
    defer fmt.Println("cleanup")  // runs last
    fmt.Println("work")
}
// Output: work, cleanup

// Multiple defers - LIFO order
defer fmt.Println("1")
defer fmt.Println("2")
defer fmt.Println("3")
// Output: 3, 2, 1

// Common use: close resources
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()  // guaranteed to run
```

---

## Functions

### Basic Functions

```go
// Simple function
func greet(name string) {
    fmt.Println("Hello", name)
}

// With return value
func add(a, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // "naked return" - returns x and y
}

// Variadic function
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
sum(1, 2, 3)       // 6
sum(1, 2, 3, 4, 5) // 15

// Spread slice into variadic
nums := []int{1, 2, 3}
sum(nums...)  // 6
```

### Anonymous Functions & Closures

```go
// Anonymous function
func() {
    fmt.Println("Anonymous!")
}()  // immediately invoked

// Assign to variable
greet := func(name string) {
    fmt.Println("Hello", name)
}
greet("World")

// Closure - captures outer variable
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
c := counter()
c()  // 1
c()  // 2
c()  // 3
```

### Function Types

```go
// Function as parameter
func apply(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

double := func(x int) int { return x * 2 }
apply([]int{1, 2, 3}, double)  // [2, 4, 6]

// Define function type
type Transformer func(int) int

func apply(nums []int, f Transformer) []int {
    // ...
}
```

---

## Arrays, Slices & Maps

### Arrays (fixed size)

```go
var a [5]int                    // [0 0 0 0 0]
b := [3]int{1, 2, 3}           // [1 2 3]
c := [...]int{1, 2, 3, 4, 5}   // size inferred: [1 2 3 4 5]

len(a)  // 5
a[0]    // first element
a[4]    // last element
```

### Slices (dynamic size)

```go
// Create slice
var s []int                     // nil slice
s := []int{1, 2, 3}            // literal
s := make([]int, 5)            // length 5, capacity 5
s := make([]int, 5, 10)        // length 5, capacity 10

// Slice from array
arr := [5]int{1, 2, 3, 4, 5}
s := arr[1:4]                  // [2 3 4] (index 1 to 3)
s := arr[:3]                   // [1 2 3] (start to 2)
s := arr[2:]                   // [3 4 5] (index 2 to end)
s := arr[:]                    // [1 2 3 4 5] (all)

// Operations
len(s)                         // length
cap(s)                         // capacity
s = append(s, 6)               // add element
s = append(s, 7, 8, 9)         // add multiple
s = append(s, other...)        // append another slice

// Copy
dst := make([]int, len(src))
copy(dst, src)

// Slice tricks
s = s[:0]                      // empty slice (keep capacity)
s = s[1:]                      // remove first element
s = s[:len(s)-1]               // remove last element
s = append(s[:i], s[i+1:]...)  // remove element at index i
```

### Maps

```go
// Create map
var m map[string]int           // nil map (can't write!)
m := make(map[string]int)      // empty map
m := map[string]int{           // literal
    "one": 1,
    "two": 2,
}

// Operations
m["three"] = 3                 // set
value := m["one"]              // get (0 if not exists)
value, ok := m["one"]          // get with existence check
delete(m, "one")               // delete
len(m)                         // number of keys

// Check if key exists
if val, ok := m["key"]; ok {
    fmt.Println("Found:", val)
} else {
    fmt.Println("Not found")
}

// Iterate (order is random!)
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}
```

---

## Structs & Methods

### Struct Definition

```go
type Person struct {
    Name    string
    Age     int
    Email   string
}

// Create instances
p1 := Person{"Alice", 30, "alice@example.com"}     // positional
p2 := Person{Name: "Bob", Age: 25}                 // named (Email = "")
p3 := Person{}                                      // zero value
var p4 Person                                       // zero value

// Access fields
fmt.Println(p1.Name)
p1.Age = 31

// Pointer to struct
p := &Person{Name: "Charlie", Age: 35}
p.Name = "Charles"  // automatic dereference (same as (*p).Name)
```

### Methods

```go
type Rectangle struct {
    Width, Height float64
}

// Value receiver (doesn't modify original)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (can modify original)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// Usage
rect := Rectangle{10, 5}
fmt.Println(rect.Area())  // 50
rect.Scale(2)
fmt.Println(rect.Area())  // 200
```

### Embedding (Composition)

```go
type Animal struct {
    Name string
}

func (a Animal) Speak() {
    fmt.Println(a.Name, "makes a sound")
}

type Dog struct {
    Animal  // embedded
    Breed string
}

// Dog inherits Animal's fields and methods
d := Dog{Animal: Animal{Name: "Buddy"}, Breed: "Labrador"}
fmt.Println(d.Name)   // Buddy (promoted field)
d.Speak()             // Buddy makes a sound (promoted method)
```

### Tags

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"-"`  // ignored
}

// Used by encoding/json, database drivers, etc.
```

---

## JSON & Encoding

### Import

```go
import "encoding/json"
```

### Marshal (Go → JSON)

**Marshal** converts a Go struct/map to a JSON byte slice.

```go
// Signature:
// func json.Marshal(v any) ([]byte, error)

type User struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email,omitempty"`  // omitted if empty
}

user := User{Name: "John", Age: 30}

// Marshal to JSON
jsonBytes, err := json.Marshal(user)
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsonBytes))  // {"name":"John","age":30}
```

### Unmarshal (JSON → Go)

**Unmarshal** converts JSON bytes into a Go struct/map.

```go
// Signature:
// func json.Unmarshal(data []byte, v any) error

jsonStr := `{"name":"John","age":30}`

var user User
err := json.Unmarshal([]byte(jsonStr), &user)  // Pass POINTER!
if err != nil {
    log.Fatal(err)
}

fmt.Println(user.Name)  // John
fmt.Println(user.Age)   // 30
```

### Visual: Marshal vs Unmarshal

```
         Marshal (serialize)
Go struct ──────────────────────► JSON bytes
User{                              []byte(`{"name":"John"}`)
  Name: "John"                            │
}                                         │ string() to print
         Unmarshal (deserialize)          ▼
Go struct ◄────────────────────── JSON string
         (needs pointer &user)    `{"name":"John"}`
```

### Why Pointer for Unmarshal?

```go
// WRONG - won't work
var user User
json.Unmarshal(data, user)   // user is a COPY, changes are lost

// CORRECT - pass pointer
var user User
json.Unmarshal(data, &user)  // &user lets function modify original
```

### Struct Tags for JSON

```go
type User struct {
    Name     string `json:"name"`           // JSON key is "name"
    Age      int    `json:"age"`            // JSON key is "age"
    Email    string `json:"email,omitempty"` // Omit if empty string
    Password string `json:"-"`               // Never include in JSON
    IsAdmin  bool   `json:"is_admin"`        // snake_case in JSON
}
```

| Tag | Meaning |
|-----|---------|
| `json:"name"` | Use "name" as JSON key |
| `json:"name,omitempty"` | Omit field if zero value |
| `json:"-"` | Never include in JSON |
| `json:",omitempty"` | Keep field name, omit if empty |

### Marshal Map (Dynamic JSON)

```go
// When you don't know the structure
data := map[string]any{
    "name":   "John",
    "age":    30,
    "active": true,
}

jsonBytes, _ := json.Marshal(data)
fmt.Println(string(jsonBytes))  // {"active":true,"age":30,"name":"John"}
```

### Unmarshal to Map (Unknown Structure)

```go
jsonStr := `{"name":"John","age":30,"scores":[85,90,78]}`

var result map[string]any
json.Unmarshal([]byte(jsonStr), &result)

name := result["name"].(string)     // Type assertion needed
age := result["age"].(float64)      // JSON numbers are float64!
scores := result["scores"].([]any)  // Arrays are []any
```

### Pretty Print JSON

```go
// MarshalIndent for readable output
jsonBytes, _ := json.MarshalIndent(user, "", "  ")
fmt.Println(string(jsonBytes))
// {
//   "name": "John",
//   "age": 30
// }
```

### JSON Quick Reference

| Operation | Code |
|-----------|------|
| Struct → JSON | `json.Marshal(user)` |
| JSON → Struct | `json.Unmarshal(data, &user)` |
| Pretty JSON | `json.MarshalIndent(user, "", "  ")` |
| Read JSON file | `os.ReadFile()` then `json.Unmarshal()` |
| Write JSON file | `json.Marshal()` then `os.WriteFile()` |
| Decode from Reader | `json.NewDecoder(reader).Decode(&user)` |
| Encode to Writer | `json.NewEncoder(writer).Encode(user)` |

### Common Mistakes

```go
// 1. Forgetting pointer in Unmarshal
json.Unmarshal(data, user)   // WRONG
json.Unmarshal(data, &user)  // CORRECT

// 2. Unexported fields (lowercase) are ignored
type User struct {
    Name string  // Exported - included in JSON
    age  int     // unexported - IGNORED by json package
}

// 3. JSON numbers are float64
var m map[string]any
json.Unmarshal([]byte(`{"count":42}`), &m)
count := m["count"].(float64)  // NOT int!
countInt := int(count)         // Convert if needed
```

---

## Strings Package

### Import

```go
import "strings"
```

### Common String Functions

```go
s := "Hello, World!"

// Check contents
strings.Contains(s, "World")     // true
strings.HasPrefix(s, "Hello")    // true
strings.HasSuffix(s, "!")        // true
strings.Count(s, "l")            // 3

// Find position
strings.Index(s, "World")        // 7
strings.LastIndex(s, "o")        // 8

// Transform
strings.ToUpper(s)               // "HELLO, WORLD!"
strings.ToLower(s)               // "hello, world!"
strings.Title("hello world")     // "Hello World" (deprecated, use cases.Title)
strings.TrimSpace("  hi  ")      // "hi"
strings.Trim("!!hi!!", "!")      // "hi"
strings.TrimPrefix("Hello", "He") // "llo"
strings.TrimSuffix("Hello", "lo") // "Hel"

// Replace
strings.Replace(s, "World", "Go", 1)   // "Hello, Go!" (replace first)
strings.ReplaceAll(s, "l", "L")        // "HeLLo, WorLd!"

// Split and Join
strings.Split("a,b,c", ",")      // []string{"a", "b", "c"}
strings.SplitN("a,b,c", ",", 2)  // []string{"a", "b,c"} (max 2 parts)
strings.Fields("  a  b  c  ")    // []string{"a", "b", "c"} (split on whitespace)
strings.Join([]string{"a","b"}, "-")  // "a-b"

// Repeat
strings.Repeat("ab", 3)          // "ababab"
```

### strings.Builder - Efficient String Building

```go
var b strings.Builder

b.WriteString("Hello")
b.WriteString(", ")
b.WriteString("World!")
b.WriteByte('!')
b.WriteRune('🎉')

result := b.String()  // "Hello, World!!🎉"

// Much faster than: s = s + "more" in a loop
```

### strings.Reader - Read String as io.Reader

```go
r := strings.NewReader("Hello, World!")

buf := make([]byte, 5)
r.Read(buf)           // buf = "Hello"
r.Read(buf)           // buf = ", Wor"
```

### Strings Quick Reference

| Operation | Code |
|-----------|------|
| Contains | `strings.Contains(s, "sub")` |
| Starts with | `strings.HasPrefix(s, "pre")` |
| Ends with | `strings.HasSuffix(s, "suf")` |
| Split | `strings.Split(s, ",")` |
| Join | `strings.Join(parts, ",")` |
| Replace all | `strings.ReplaceAll(s, "old", "new")` |
| Trim spaces | `strings.TrimSpace(s)` |
| Upper/Lower | `strings.ToUpper(s)` / `strings.ToLower(s)` |
| Build efficiently | `strings.Builder` |

---

## File I/O

### Import

```go
import (
    "os"
    "io"
    "bufio"
)
```

### Read Entire File

```go
// Simple - read entire file into memory
data, err := os.ReadFile("file.txt")
if err != nil {
    log.Fatal(err)
}
content := string(data)
```

### Write Entire File

```go
data := []byte("Hello, World!")
err := os.WriteFile("file.txt", data, 0644)  // 0644 = permissions
if err != nil {
    log.Fatal(err)
}
```

### Open, Read, Close (More Control)

```go
// Open for reading
f, err := os.Open("file.txt")
if err != nil {
    log.Fatal(err)
}
defer f.Close()  // ALWAYS close!

// Read into buffer
buf := make([]byte, 1024)
n, err := f.Read(buf)
if err != nil && err != io.EOF {
    log.Fatal(err)
}
fmt.Println(string(buf[:n]))
```

### Read Line by Line

```go
f, err := os.Open("file.txt")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

scanner := bufio.NewScanner(f)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}

if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
```

### Create and Write

```go
// Create (truncates if exists)
f, err := os.Create("file.txt")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

f.WriteString("Hello\n")
f.Write([]byte("World\n"))

// Buffered writing (better for many writes)
w := bufio.NewWriter(f)
w.WriteString("Buffered write\n")
w.Flush()  // Don't forget to flush!
```

### Append to File

```go
f, err := os.OpenFile("file.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
if err != nil {
    log.Fatal(err)
}
defer f.Close()

f.WriteString("Appended line\n")
```

### File Flags

| Flag | Description |
|------|-------------|
| `os.O_RDONLY` | Read only |
| `os.O_WRONLY` | Write only |
| `os.O_RDWR` | Read and write |
| `os.O_CREATE` | Create if not exists |
| `os.O_APPEND` | Append to file |
| `os.O_TRUNC` | Truncate file |

### Check if File Exists

```go
if _, err := os.Stat("file.txt"); os.IsNotExist(err) {
    fmt.Println("File does not exist")
}
```

### Directory Operations

```go
// Create directory
os.Mkdir("mydir", 0755)
os.MkdirAll("path/to/dir", 0755)  // Creates parents too

// Remove
os.Remove("file.txt")             // File or empty dir
os.RemoveAll("mydir")             // Dir and contents

// List directory
entries, _ := os.ReadDir(".")
for _, e := range entries {
    fmt.Println(e.Name(), e.IsDir())
}

// Get working directory
wd, _ := os.Getwd()

// Change directory
os.Chdir("/path/to/dir")
```

### Copy File

```go
func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    return err
}
```

### Path Operations

```go
import "path/filepath"

filepath.Join("dir", "subdir", "file.txt")  // "dir/subdir/file.txt" (OS-aware)
filepath.Dir("/a/b/c.txt")                   // "/a/b"
filepath.Base("/a/b/c.txt")                  // "c.txt"
filepath.Ext("/a/b/c.txt")                   // ".txt"
filepath.Abs("file.txt")                     // Full absolute path
filepath.Glob("*.go")                        // Match pattern
```

---

## HTTP Client & Server

### Import

```go
import "net/http"
```

### Simple GET Request

```go
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(body))
fmt.Println("Status:", resp.StatusCode)
```

### GET with Timeout (Using Context)

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
if err != nil {
    log.Fatal(err)
}

resp, err := http.DefaultClient.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

### POST Request with JSON

```go
data := map[string]string{"name": "John", "email": "john@example.com"}
jsonData, _ := json.Marshal(data)

resp, err := http.Post(
    "https://api.example.com/users",
    "application/json",
    bytes.NewBuffer(jsonData),
)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

fmt.Println("Status:", resp.StatusCode)
```

### Custom Request with Headers

```go
client := &http.Client{Timeout: 10 * time.Second}

req, err := http.NewRequest("POST", "https://api.example.com/data", bytes.NewBuffer(jsonData))
if err != nil {
    log.Fatal(err)
}

req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer token123")

resp, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

### Simple HTTP Server

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Home page")
    })

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Handler Methods

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Request info
    fmt.Println("Method:", r.Method)           // GET, POST, etc.
    fmt.Println("URL:", r.URL.Path)            // /users/123
    fmt.Println("Query:", r.URL.Query())       // map of query params

    // Headers
    fmt.Println("User-Agent:", r.Header.Get("User-Agent"))

    // Read body
    body, _ := io.ReadAll(r.Body)
    defer r.Body.Close()

    // Response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)  // 200
    w.Write([]byte(`{"status": "ok"}`))
}
```

### JSON API Handler

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case "GET":
        user := User{ID: 1, Name: "John"}
        json.NewEncoder(w).Encode(user)

    case "POST":
        var user User
        json.NewDecoder(r.Body).Decode(&user)
        fmt.Printf("Received: %+v\n", user)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}
```

### HTTP Quick Reference

| Operation | Code |
|-----------|------|
| GET request | `http.Get(url)` |
| POST request | `http.Post(url, contentType, body)` |
| Custom request | `http.NewRequest(method, url, body)` |
| Set header | `req.Header.Set("Key", "Value")` |
| Read response | `io.ReadAll(resp.Body)` |
| Start server | `http.ListenAndServe(":8080", nil)` |
| Register handler | `http.HandleFunc("/path", handler)` |
| Send JSON | `json.NewEncoder(w).Encode(data)` |
| Read JSON | `json.NewDecoder(r.Body).Decode(&data)` |

---

## Regular Expressions

### Import

```go
import "regexp"
```

### Basic Matching

```go
// Check if matches
matched, _ := regexp.MatchString(`\d+`, "abc123")  // true

// Compile for reuse (better performance)
re := regexp.MustCompile(`\d+`)  // Panics if invalid
re.MatchString("abc123")         // true
```

### Find Matches

```go
re := regexp.MustCompile(`\d+`)

// Find first match
re.FindString("abc123def456")           // "123"

// Find all matches
re.FindAllString("abc123def456", -1)    // ["123", "456"]
re.FindAllString("abc123def456", 1)     // ["123"] (limit to 1)

// Find with position
loc := re.FindStringIndex("abc123")     // [3, 6] (start, end)
```

### Capture Groups

```go
re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
matches := re.FindStringSubmatch("user@example.com")
// matches[0] = "user@example.com" (full match)
// matches[1] = "user"
// matches[2] = "example"
// matches[3] = "com"

// Named groups
re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
matches := re.FindStringSubmatch("user@example.com")
names := re.SubexpNames()  // ["", "user", "domain"]
```

### Replace

```go
re := regexp.MustCompile(`\d+`)

// Replace all
re.ReplaceAllString("abc123def456", "X")       // "abcXdefX"

// Replace with function
re.ReplaceAllStringFunc("abc123def456", func(s string) string {
    return "[" + s + "]"
})  // "abc[123]def[456]"
```

### Split

```go
re := regexp.MustCompile(`\s+`)  // Split on whitespace
re.Split("a  b   c", -1)         // ["a", "b", "c"]
```

### Common Regex Patterns

| Pattern | Matches |
|---------|---------|
| `\d` | Digit [0-9] |
| `\D` | Non-digit |
| `\w` | Word char [a-zA-Z0-9_] |
| `\W` | Non-word char |
| `\s` | Whitespace |
| `\S` | Non-whitespace |
| `.` | Any char (except newline) |
| `*` | 0 or more |
| `+` | 1 or more |
| `?` | 0 or 1 |
| `{n}` | Exactly n |
| `{n,m}` | Between n and m |
| `^` | Start of string |
| `$` | End of string |
| `[abc]` | a, b, or c |
| `[^abc]` | Not a, b, or c |
| `(...)` | Capture group |
| `(?:...)` | Non-capture group |

---

## Sorting

### Import

```go
import "sort"
```

### Sort Built-in Types

```go
// Sort ints
nums := []int{3, 1, 4, 1, 5, 9}
sort.Ints(nums)        // [1, 1, 3, 4, 5, 9]

// Sort strings
strs := []string{"banana", "apple", "cherry"}
sort.Strings(strs)     // ["apple", "banana", "cherry"]

// Sort floats
floats := []float64{3.14, 1.41, 2.72}
sort.Float64s(floats)  // [1.41, 2.72, 3.14]
```

### Check if Sorted

```go
sort.IntsAreSorted([]int{1, 2, 3})       // true
sort.StringsAreSorted([]string{"a","b"}) // true
```

### Reverse Sort

```go
nums := []int{3, 1, 4, 1, 5}
sort.Sort(sort.Reverse(sort.IntSlice(nums)))  // [5, 4, 3, 1, 1]
```

### Sort Custom Types

```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}

// Sort by Age
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
// [{Bob 25} {Alice 30} {Charlie 35}]

// Sort by Name
sort.Slice(people, func(i, j int) bool {
    return people[i].Name < people[j].Name
})
```

### Stable Sort (Preserves Order of Equal Elements)

```go
sort.SliceStable(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

### Binary Search

```go
nums := []int{1, 2, 3, 4, 5}  // Must be sorted!

// Search returns index where value would be inserted
i := sort.SearchInts(nums, 3)  // 2

// Check if found
if i < len(nums) && nums[i] == 3 {
    fmt.Println("Found at index", i)
}
```

---

## Random Numbers

### Import

```go
import (
    "math/rand"
    "time"
)
```

### Seeding (Required for Variety)

```go
// Go < 1.20: Must seed manually
rand.Seed(time.Now().UnixNano())

// Go 1.20+: Auto-seeded, but can reseed
rand.Seed(time.Now().UnixNano())
```

### Generate Random Numbers

```go
rand.Int()                    // Random int
rand.Intn(100)               // 0 to 99
rand.Float64()               // 0.0 to 1.0

// Random in range [min, max)
min, max := 10, 20
n := rand.Intn(max-min) + min  // 10 to 19
```

### Random Selection

```go
// Random element from slice
items := []string{"a", "b", "c", "d"}
random := items[rand.Intn(len(items))]

// Shuffle slice
rand.Shuffle(len(items), func(i, j int) {
    items[i], items[j] = items[j], items[i]
})
```

### Random String

```go
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

randomString(10)  // "aBcDeFgHiJ"
```

### Crypto Random (Secure)

```go
import "crypto/rand"

// For security-sensitive code (tokens, passwords, keys)
b := make([]byte, 16)
_, err := rand.Read(b)  // Cryptographically secure
```

---

## Logging

### Import

```go
import "log"
```

### Basic Logging

```go
log.Print("Info message")               // 2024/01/15 14:30:45 Info message
log.Println("With newline")             // Same but adds newline
log.Printf("User %s logged in", name)   // Formatted

log.Fatal("Error!")     // Print + os.Exit(1)
log.Panic("Error!")     // Print + panic()
```

### Configure Logger

```go
// Set flags (prefix format)
log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// 2024/01/15 14:30:45 main.go:10: message

// All flags:
log.Ldate         // 2024/01/15
log.Ltime         // 14:30:45
log.Lmicroseconds // 14:30:45.123456
log.Llongfile     // /full/path/main.go:10
log.Lshortfile    // main.go:10
log.LUTC          // Use UTC
log.Lmsgprefix    // Prefix before message (not before date)
log.LstdFlags     // Ldate | Ltime

// Set prefix
log.SetPrefix("[APP] ")
// [APP] 2024/01/15 14:30:45 message
```

### Log to File

```go
f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}
defer f.Close()

log.SetOutput(f)
log.Println("This goes to file")

// Or create new logger
logger := log.New(f, "[APP] ", log.LstdFlags|log.Lshortfile)
logger.Println("Custom logger")
```

### Multiple Log Levels (Manual)

```go
var (
    Info  = log.New(os.Stdout, "INFO: ", log.LstdFlags)
    Warn  = log.New(os.Stdout, "WARN: ", log.LstdFlags)
    Error = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
)

Info.Println("Server started")
Warn.Println("Memory high")
Error.Println("Database connection failed")
```

### log/slog (Go 1.21+) - Structured Logging

```go
import "log/slog"

// Default text output
slog.Info("User logged in", "user", "john", "ip", "192.168.1.1")
// time=2024-01-15T14:30:45 level=INFO msg="User logged in" user=john ip=192.168.1.1

// JSON output
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)
slog.Info("User logged in", "user", "john")
// {"time":"2024-01-15T14:30:45","level":"INFO","msg":"User logged in","user":"john"}

// Log levels
slog.Debug("debug message")
slog.Info("info message")
slog.Warn("warning message")
slog.Error("error message")
```

---

## Interfaces

### Interface Definition

```go
type Speaker interface {
    Speak() string
}

type Dog struct{ Name string }
func (d Dog) Speak() string { return "Woof!" }

type Cat struct{ Name string }
func (c Cat) Speak() string { return "Meow!" }

// Both satisfy Speaker interface
func MakeSpeak(s Speaker) {
    fmt.Println(s.Speak())
}

MakeSpeak(Dog{"Buddy"})  // Woof!
MakeSpeak(Cat{"Whiskers"})  // Meow!
```

### Empty Interface

```go
// interface{} or 'any' accepts any type
func Print(v interface{}) {
    fmt.Println(v)
}

// Same with 'any' (Go 1.18+)
func Print(v any) {
    fmt.Println(v)
}

// Type assertion
var i interface{} = "hello"
s := i.(string)        // panics if wrong type
s, ok := i.(string)    // safe - ok is false if wrong type
```

### Common Interfaces

```go
// Stringer (like toString)
type Stringer interface {
    String() string
}

type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

fmt.Println(Person{"Alice", 30})  // Alice (30 years)

// Error interface
type error interface {
    Error() string
}

// Reader/Writer
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

---

## Error Handling

### Basic Error Handling

```go
import "errors"

// Return error
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Handle error
result, err := divide(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

### Custom Errors

```go
// Using fmt.Errorf
err := fmt.Errorf("failed to process %s: %w", filename, originalErr)

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Usage
return ValidationError{Field: "email", Message: "invalid format"}
```

### Error Wrapping (Go 1.13+)

```go
import "errors"

// Wrap error
err := fmt.Errorf("failed to read config: %w", originalErr)

// Unwrap and check
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("File not found")
}

// Get underlying error type
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path:", pathErr.Path)
}
```

### Panic and Recover

```go
// Panic - stops normal execution
func mustParse(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        panic(err)  // use sparingly!
    }
    return i
}

// Recover - catch panic
func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()

    panic("something bad happened")
}
```

### Go vs Try/Catch (JavaScript/Python/Java)

Go doesn't have try/catch. Instead, it uses **explicit error returns**.

**JavaScript try/catch:**
```javascript
try {
    const data = readFile("config.json");
    const parsed = JSON.parse(data);
    console.log(parsed.name);
} catch (error) {
    console.log("Error:", error.message);
} finally {
    cleanup();
}
```

**Go equivalent:**
```go
data, err := os.ReadFile("config.json")
if err != nil {
    fmt.Println("Error:", err)
    return
}

var parsed Config
err = json.Unmarshal(data, &parsed)
if err != nil {
    fmt.Println("Error:", err)
    return
}

fmt.Println(parsed.Name)

// "finally" equivalent: use defer at the start
defer cleanup()
```

### Comparison Table

| Concept | JavaScript/Java | Go |
|---------|-----------------|-----|
| Throw error | `throw new Error("msg")` | `return errors.New("msg")` |
| Catch error | `catch (e) { ... }` | `if err != nil { ... }` |
| Finally | `finally { ... }` | `defer func() { ... }()` |
| Rethrow | `throw e` | `return err` |
| Custom error | `class MyError extends Error` | `type MyError struct` with `Error() string` |
| Crash program | `throw` (uncaught) | `panic("msg")` |
| Catch crash | N/A (process dies) | `recover()` in deferred function |

### Pattern: Handle Multiple Errors

**JavaScript:**
```javascript
try {
    step1();
    step2();
    step3();
} catch (e) {
    // One handler for all errors
}
```

**Go:**
```go
if err := step1(); err != nil {
    return fmt.Errorf("step1 failed: %w", err)
}
if err := step2(); err != nil {
    return fmt.Errorf("step2 failed: %w", err)
}
if err := step3(); err != nil {
    return fmt.Errorf("step3 failed: %w", err)
}
```

### When to Use panic/recover vs Error Returns

| Situation | Use |
|-----------|-----|
| Expected failures (file not found, bad input) | **Return error** |
| Programming bugs (nil pointer, out of bounds) | **Panic** (automatic) |
| Unrecoverable state | **Panic** |
| Library boundary (must not crash caller) | **Recover** |
| HTTP handler (one request shouldn't crash server) | **Recover** |

**Rule of thumb:** If the caller can handle it, return an error. If it's a bug, let it panic.

### Visual: Go Error Flow vs Try/Catch

```
JavaScript/Python/Java:
┌──────────────────────────────────────────────┐
│  try {                                       │
│      riskyOperation();  ──throws──►  catch   │
│      anotherOperation();             │       │
│  } catch (e) {         ◄─────────────┘       │
│      handleError(e);                         │
│  }                                           │
└──────────────────────────────────────────────┘
Errors "jump" to catch block

Go:
┌──────────────────────────────────────────────┐
│  result, err := riskyOperation()             │
│  if err != nil {  ◄── Check immediately      │
│      return err   ◄── Handle or propagate    │
│  }                                           │
│                                              │
│  result2, err := anotherOperation()          │
│  if err != nil {  ◄── Check again            │
│      return err                              │
│  }                                           │
└──────────────────────────────────────────────┘
Errors are values, checked at each step
```

---

## Pointers

### Basics

```go
x := 42
p := &x      // p is pointer to x
fmt.Println(*p)  // 42 (dereference)
*p = 21      // modify through pointer
fmt.Println(x)   // 21

// Pointer to new value
p := new(int)    // *int pointing to 0
*p = 100
```

### Pointers with Functions

```go
// Pass by value (copy)
func double(x int) {
    x *= 2  // doesn't affect original
}

// Pass by pointer (modify original)
func doublePtr(x *int) {
    *x *= 2
}

n := 5
double(n)
fmt.Println(n)   // 5 (unchanged)

doublePtr(&n)
fmt.Println(n)   // 10 (changed!)
```

### When to Use Pointers

```go
// 1. When you need to modify the original
func (p *Person) SetAge(age int) {
    p.Age = age
}

// 2. Large structs (avoid copying)
func ProcessLargeData(data *LargeStruct) { }

// 3. When value might be nil
func findUser(id int) *User {
    // return nil if not found
}
```

---

## make vs new

### new() - Allocates Memory, Returns Pointer

`new(T)` allocates memory for type T, initializes to **zero value**, returns `*T`.

```go
// new returns a pointer to zero-valued memory
p := new(int)       // *int pointing to 0
fmt.Println(*p)     // 0

s := new(string)    // *string pointing to ""
fmt.Println(*s)     // "" (empty string)

// Equivalent to:
var i int
p := &i             // Same result as new(int)
```

### make() - Creates Slices, Maps, Channels

`make()` creates and initializes slices, maps, and channels. Returns the type itself (not a pointer).

```go
// make for slices
s := make([]int, 5)        // length 5, capacity 5
s := make([]int, 5, 10)    // length 5, capacity 10

// make for maps
m := make(map[string]int)  // empty map, ready to use

// make for channels
ch := make(chan int)       // unbuffered channel
ch := make(chan int, 5)    // buffered channel (capacity 5)
```

### Comparison Table

| | `new(T)` | `make(T)` |
|--|----------|-----------|
| Works with | Any type | Slices, Maps, Channels only |
| Returns | Pointer (`*T`) | Value (`T`) |
| Initializes to | Zero value | Usable data structure |
| Memory | Allocates | Allocates + initializes internal data |

### Visual: new vs make

```
new(int):
┌─────────────────────────┐
│  new(int)               │
│     │                   │
│     ▼                   │
│  ┌─────┐                │
│  │  0  │  ← zero value  │
│  └──┬──┘                │
│     │                   │
│  Returns: *int (pointer)│
└─────────────────────────┘

make([]int, 3):
┌─────────────────────────────────┐
│  make([]int, 3)                 │
│     │                           │
│     ▼                           │
│  ┌─────┬─────┬─────┐            │
│  │  0  │  0  │  0  │  ← backing │
│  └─────┴─────┴─────┘    array   │
│     │                           │
│  Returns: []int (slice header)  │
│  - pointer to array             │
│  - length: 3                    │
│  - capacity: 3                  │
└─────────────────────────────────┘
```

### When to Use Which

```go
// Use new() when you need a pointer to a simple type
p := new(MyStruct)   // Same as: p := &MyStruct{}

// Use make() for slices, maps, channels
s := make([]int, 0, 100)      // Pre-allocate capacity
m := make(map[string]int)     // Must make() before using map!
ch := make(chan int, 10)      // Buffered channel
```

### Common Mistake: nil map

```go
var m map[string]int  // m is nil!
m["key"] = 1          // PANIC: assignment to entry in nil map

m := make(map[string]int)  // CORRECT: now it's usable
m["key"] = 1               // Works!
```

---

## Goroutines

### Starting a Goroutine

```go
// Basic goroutine
go func() {
    fmt.Println("Running in goroutine")
}()

// Named function
func sayHello() {
    fmt.Println("Hello!")
}
go sayHello()

// With arguments
go func(msg string) {
    fmt.Println(msg)
}("Hello from goroutine")
```

### Goroutine Basics

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Start goroutine
    go func() {
        for i := 0; i < 3; i++ {
            fmt.Println("Goroutine:", i)
            time.Sleep(100 * time.Millisecond)
        }
    }()

    // Main continues immediately
    for i := 0; i < 3; i++ {
        fmt.Println("Main:", i)
        time.Sleep(100 * time.Millisecond)
    }

    // Wait a bit for goroutine to finish (bad practice - use sync.WaitGroup!)
    time.Sleep(500 * time.Millisecond)
}
```

### Goroutine Closure Gotcha

```go
// WRONG - all goroutines share same variable
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // Might print: 3, 3, 3
    }()
}

// CORRECT - pass value as argument
for i := 0; i < 3; i++ {
    go func(n int) {
        fmt.Println(n)  // Prints: 0, 1, 2 (any order)
    }(i)
}

// CORRECT (Go 1.22+) - loop variable is per-iteration
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // Works correctly in Go 1.22+
    }()
}
```

### Visual: Goroutine Execution

```
main() starts
    │
    ├── go func1()  ──────────────────────►  [runs concurrently]
    │                                              │
    ├── go func2()  ──────────────────────►  [runs concurrently]
    │                                              │
    ├── main continues...                          │
    │                                              │
    ▼                                              ▼
main() ends ─── Program exits! All goroutines killed!
```

**Important:** When `main()` exits, all goroutines are terminated immediately!

### Waiting for Goroutines (Preview)

```go
// Use sync.WaitGroup (see sync Package section)
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    // do work
}()

wg.Wait()  // Blocks until goroutine finishes
```

---

## Channels

### Creation

```go
// Unbuffered (synchronous)
ch := make(chan int)

// Buffered (can hold N values)
ch := make(chan int, 5)

// Directional (for function parameters)
func send(ch chan<- int) { }    // send-only
func recv(ch <-chan int) { }    // receive-only
```

### Operations

```go
ch <- value     // send (blocks if full/no receiver)
value := <-ch   // receive (blocks if empty)
close(ch)       // close channel (sender only!)

// Check if channel is closed
value, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}

// Range over channel (until closed)
for value := range ch {
    fmt.Println(value)
}
```

### Select

```go
select {
case msg := <-ch1:
    fmt.Println("From ch1:", msg)
case ch2 <- value:
    fmt.Println("Sent to ch2")
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("Nothing ready")
}
```

---

## sync Package

### Import

```go
import "sync"
```

### sync.WaitGroup - Wait for Goroutines

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)  // Add BEFORE starting goroutine

        go func(id int) {
            defer wg.Done()  // Mark done when goroutine exits
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Second)
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }

    wg.Wait()  // Blocks until counter reaches 0
    fmt.Println("All workers finished")
}
```

### WaitGroup Methods

| Method | Description |
|--------|-------------|
| `wg.Add(n)` | Add n to counter (call before goroutine) |
| `wg.Done()` | Decrement counter by 1 (same as Add(-1)) |
| `wg.Wait()` | Block until counter is 0 |

### sync.Mutex - Protect Shared Data

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()         // Acquire lock
    c.count++           // Safe to modify
    c.mu.Unlock()       // Release lock
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()  // defer ensures unlock even if panic
    return c.count
}

func main() {
    counter := &SafeCounter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Println("Count:", counter.Value())  // Always 1000
}
```

### sync.RWMutex - Multiple Readers, Single Writer

```go
type SafeCache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func (c *SafeCache) Get(key string) string {
    c.mu.RLock()          // Read lock - multiple readers OK
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()           // Write lock - exclusive access
    defer c.mu.Unlock()
    c.data[key] = value
}
```

| Lock Type | Multiple Concurrent? | Blocks? |
|-----------|---------------------|---------|
| `RLock()` | Yes (readers) | Only blocks writers |
| `Lock()` | No | Blocks all |

### sync.Once - Run Code Exactly Once

```go
var once sync.Once
var config *Config

func GetConfig() *Config {
    once.Do(func() {
        // This runs only once, even from multiple goroutines
        config = loadConfig()
        fmt.Println("Config loaded")
    })
    return config
}

// Safe to call from multiple goroutines
go GetConfig()  // Loads config
go GetConfig()  // Returns cached config
go GetConfig()  // Returns cached config
```

### sync.Map - Concurrent Map

```go
var m sync.Map

// Store
m.Store("key", "value")

// Load
value, ok := m.Load("key")
if ok {
    fmt.Println(value.(string))
}

// LoadOrStore - returns existing or stores new
actual, loaded := m.LoadOrStore("key", "default")

// Delete
m.Delete("key")

// Range - iterate
m.Range(func(key, value any) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true  // continue iteration
})
```

### sync/atomic - Atomic Operations

```go
import "sync/atomic"

var counter int64

// Atomic add
atomic.AddInt64(&counter, 1)

// Atomic load
value := atomic.LoadInt64(&counter)

// Atomic store
atomic.StoreInt64(&counter, 100)

// Compare and swap
swapped := atomic.CompareAndSwapInt64(&counter, 100, 200)  // true if was 100
```

### sync Quick Reference

| Type | Use Case |
|------|----------|
| `WaitGroup` | Wait for goroutines to complete |
| `Mutex` | Protect shared data (one accessor) |
| `RWMutex` | Multiple readers, single writer |
| `Once` | Initialize exactly once (singletons) |
| `Map` | Concurrent map access |
| `atomic` | Simple counters, flags |

---

## Context

Context is used for **cancellation**, **timeouts**, and **passing request-scoped values** across function calls and goroutines.

### Import

```go
import "context"
```

### Creating Contexts

```go
// Background - root context, never cancelled (starting point)
ctx := context.Background()

// TODO - placeholder when unsure which context to use
ctx := context.TODO()
```

### context.WithTimeout - Auto-Cancel After Duration

```go
// Signature:
// func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // ALWAYS call cancel to release resources!

// ctx.Done() will receive after 5 seconds
// ctx.Err() will return context.DeadlineExceeded
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `parent` | `context.Context` | Parent context (usually `context.Background()`) |
| `timeout` | `time.Duration` | How long until auto-cancel |

**Returns:**
| Return | Type | Description |
|--------|------|-------------|
| `ctx` | `context.Context` | New context with timeout |
| `cancel` | `context.CancelFunc` | Function to manually cancel early |

### context.WithDeadline - Cancel at Specific Time

```go
// Signature:
// func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)

deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

### context.WithCancel - Manual Cancellation

```go
// Signature:
// func WithCancel(parent Context) (Context, CancelFunc)

ctx, cancel := context.WithCancel(context.Background())

// Later, when you want to stop everything:
cancel()  // Signals all goroutines using this ctx to stop
```

### Using Context in Functions

```go
// Pass ctx as FIRST parameter (Go convention)
func doWork(ctx context.Context, data string) error {
    // Check if already cancelled
    select {
    case <-ctx.Done():
        return ctx.Err()  // Return early
    default:
        // Continue working
    }

    // ... do work ...
    return nil
}
```

### Checking for Cancellation

```go
// ctx.Done() - channel that closes when cancelled
<-ctx.Done()  // Blocks until cancelled

// ctx.Err() - why it was cancelled
ctx.Err()  // Returns:
           // - nil                       (not cancelled yet)
           // - context.Canceled          (cancel() was called)
           // - context.DeadlineExceeded  (timeout/deadline reached)

// Common pattern: check in select
select {
case <-ctx.Done():
    fmt.Println("Cancelled:", ctx.Err())
    return
case result := <-workChannel:
    fmt.Println("Got result:", result)
}
```

### Context with Values (Use Sparingly)

```go
// Signature:
// func WithValue(parent Context, key, val any) Context

type contextKey string
const userIDKey contextKey = "userID"

// Set value
ctx := context.WithValue(context.Background(), userIDKey, "user-123")

// Get value
userID := ctx.Value(userIDKey).(string)  // "user-123"
```

### HTTP Request Context

```go
// HTTP handlers get context from request
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Context that cancels if client disconnects

    // Pass to database, API calls, etc.
    result, err := queryDatabase(ctx, "SELECT ...")
}
```

### Context Quick Reference

| Function | Purpose | Auto-cancels? |
|----------|---------|---------------|
| `context.Background()` | Root context | Never |
| `context.TODO()` | Placeholder | Never |
| `context.WithCancel(parent)` | Manual cancel | When `cancel()` called |
| `context.WithTimeout(parent, duration)` | Timeout | After duration |
| `context.WithDeadline(parent, time)` | Deadline | At specific time |
| `context.WithValue(parent, key, val)` | Pass values | Inherits from parent |

| Method | Returns | Description |
|--------|---------|-------------|
| `ctx.Done()` | `<-chan struct{}` | Channel that closes on cancel |
| `ctx.Err()` | `error` | `nil`, `Canceled`, or `DeadlineExceeded` |
| `ctx.Deadline()` | `time.Time, bool` | When it will timeout (if set) |
| `ctx.Value(key)` | `any` | Get value by key |

### Example: HTTP Call with Timeout

```go
func fetchAPI(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err  // Returns error if ctx times out
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

// Usage
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := fetchAPI(ctx, "https://api.example.com/data")
if err != nil {
    if err == context.DeadlineExceeded {
        fmt.Println("Request timed out!")
    }
}
```

---

## Testing

### Test File Naming

```
myfile.go       → myfile_test.go    (same package)
math.go         → math_test.go
```

### Basic Test

```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

### Run Tests

```bash
go test                    # Run tests in current package
go test ./...              # Run all tests recursively
go test -v                 # Verbose output
go test -run TestAdd       # Run specific test
go test -cover             # Show coverage
go test -coverprofile=c.out && go tool cover -html=c.out  # Coverage report
```

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Test Helpers

```go
// Error methods
t.Error("message")           // Log error, continue
t.Errorf("got %d", val)      // Formatted error, continue
t.Fatal("message")           // Log error, stop test
t.Fatalf("got %d", val)      // Formatted error, stop test

// Other
t.Log("debug info")          // Only shown with -v
t.Skip("skipping this test") // Skip test
t.Helper()                   // Mark as helper (better error locations)
```

### Setup and Teardown

```go
func TestMain(m *testing.M) {
    // Setup before all tests
    setup()

    code := m.Run()  // Run all tests

    // Teardown after all tests
    teardown()

    os.Exit(code)
}

// Per-test setup
func TestSomething(t *testing.T) {
    // Setup
    cleanup := setupTest()
    defer cleanup()  // Teardown

    // Test code
}
```

### Benchmarks

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

// Run: go test -bench=.
// Run: go test -bench=BenchmarkAdd -benchmem
```

### Example Tests (Documentation)

```go
func ExampleAdd() {
    fmt.Println(Add(2, 3))
    // Output: 5
}
```

---

## Go Modules & Commands

### Initialize Module

```bash
go mod init github.com/user/myproject
```

Creates `go.mod`:
```
module github.com/user/myproject

go 1.21
```

### Common Commands

```bash
# Run
go run main.go              # Run single file
go run .                    # Run package

# Build
go build                    # Build current package
go build -o myapp           # Build with output name
go install                  # Build and install to $GOPATH/bin

# Dependencies
go get github.com/pkg/errors         # Add dependency
go get github.com/pkg/errors@v0.9.1  # Specific version
go get -u ./...                      # Update all dependencies
go mod tidy                          # Clean up go.mod/go.sum
go mod download                      # Download dependencies
go mod vendor                        # Copy deps to vendor/

# Tools
go fmt ./...               # Format all code
go vet ./...               # Static analysis
go doc fmt.Println         # View documentation
go generate ./...          # Run go:generate directives
```

### go.mod File

```
module github.com/user/myproject

go 1.21

require (
    github.com/pkg/errors v0.9.1
    github.com/gin-gonic/gin v1.9.0
)

require (
    // indirect dependencies (managed automatically)
    golang.org/x/sys v0.5.0 // indirect
)
```

### Replace Directive (Local Development)

```
// In go.mod - use local version instead of remote
replace github.com/user/mylib => ../mylib
```

### Build Tags

```go
// +build linux,amd64

package main
// This file only compiles on linux/amd64
```

Or Go 1.17+:
```go
//go:build linux && amd64
```

### Cross-Compilation

```bash
# Build for different OS/architecture
GOOS=linux GOARCH=amd64 go build -o myapp-linux
GOOS=windows GOARCH=amd64 go build -o myapp.exe
GOOS=darwin GOARCH=arm64 go build -o myapp-mac
```

| GOOS | Platform |
|------|----------|
| `linux` | Linux |
| `darwin` | macOS |
| `windows` | Windows |

---

## Generics

### Generic Function (Go 1.18+)

```go
// Type parameter in square brackets
func Min[T int | float64](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Usage
Min[int](3, 5)        // 3
Min(3.14, 2.71)       // 2.71 (type inferred)
```

### Type Constraints

```go
// Using built-in constraints
import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Custom constraint
type Number interface {
    int | int64 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

### Generic Struct

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
    n := len(s.items) - 1
    item := s.items[n]
    s.items = s.items[:n]
    return item
}

// Usage
s := Stack[int]{}
s.Push(1)
s.Push(2)
s.Pop()  // 2
```

### Common Constraints

```go
// any - any type (same as interface{})
func Print[T any](v T) { fmt.Println(v) }

// comparable - types that support == and !=
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

// constraints.Ordered - types that support < > <= >=
// Includes: int, int8, int16, int32, int64,
//           uint, uint8, uint16, uint32, uint64,
//           float32, float64, string
```

### Type Approximation (~)

```go
// ~ means "underlying type"
type MyInt int

type Integer interface {
    ~int | ~int64  // Includes types with int/int64 as underlying type
}

func Double[T Integer](n T) T {
    return n * 2
}

var x MyInt = 5
Double(x)  // Works! MyInt's underlying type is int
```

---

## Command Line Args

### os.Args - Basic Arguments

```go
import "os"

func main() {
    // os.Args[0] is the program name
    // os.Args[1:] are the arguments

    fmt.Println("Program:", os.Args[0])
    fmt.Println("Args:", os.Args[1:])

    if len(os.Args) > 1 {
        fmt.Println("First arg:", os.Args[1])
    }
}
```

```bash
./myapp hello world
# Program: ./myapp
# Args: [hello world]
# First arg: hello
```

### flag Package - Parsed Arguments

```go
import "flag"

func main() {
    // Define flags
    name := flag.String("name", "World", "Name to greet")
    age := flag.Int("age", 0, "Your age")
    verbose := flag.Bool("v", false, "Verbose output")

    // Parse command line
    flag.Parse()

    // Use values (flags are pointers!)
    fmt.Printf("Hello, %s!\n", *name)
    fmt.Printf("Age: %d\n", *age)

    if *verbose {
        fmt.Println("Verbose mode enabled")
    }

    // Non-flag arguments
    fmt.Println("Other args:", flag.Args())
}
```

```bash
./myapp -name=John -age=30 -v extra args
# Hello, John!
# Age: 30
# Verbose mode enabled
# Other args: [extra args]

./myapp -h  # Shows help automatically
```

### flag with Variables

```go
var (
    port int
    host string
)

func init() {
    flag.IntVar(&port, "port", 8080, "Server port")
    flag.StringVar(&host, "host", "localhost", "Server host")
}

func main() {
    flag.Parse()
    fmt.Printf("Server: %s:%d\n", host, port)
}
```

### Subcommands

```go
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: myapp <command> [args]")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "get":
        getCmd := flag.NewFlagSet("get", flag.ExitOnError)
        url := getCmd.String("url", "", "URL to fetch")
        getCmd.Parse(os.Args[2:])
        fmt.Println("Getting:", *url)

    case "post":
        postCmd := flag.NewFlagSet("post", flag.ExitOnError)
        data := postCmd.String("data", "", "Data to post")
        postCmd.Parse(os.Args[2:])
        fmt.Println("Posting:", *data)

    default:
        fmt.Println("Unknown command:", os.Args[1])
    }
}
```

---

## Environment Variables

### Import

```go
import "os"
```

### Read Environment Variables

```go
// Get value (empty string if not set)
value := os.Getenv("HOME")

// Get with existence check
value, exists := os.LookupEnv("API_KEY")
if !exists {
    log.Fatal("API_KEY not set")
}

// Get with default
func getEnvOrDefault(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}

port := getEnvOrDefault("PORT", "8080")
```

### Set Environment Variables

```go
os.Setenv("MY_VAR", "my_value")

// Unset
os.Unsetenv("MY_VAR")
```

### List All Environment Variables

```go
for _, env := range os.Environ() {
    pair := strings.SplitN(env, "=", 2)
    fmt.Printf("%s = %s\n", pair[0], pair[1])
}
```

### Common Pattern: Config from Environment

```go
type Config struct {
    Port     string
    DBHost   string
    APIKey   string
    Debug    bool
}

func LoadConfig() *Config {
    return &Config{
        Port:   getEnvOrDefault("PORT", "8080"),
        DBHost: getEnvOrDefault("DB_HOST", "localhost"),
        APIKey: os.Getenv("API_KEY"),  // Required, check separately
        Debug:  os.Getenv("DEBUG") == "true",
    }
}

func main() {
    config := LoadConfig()

    if config.APIKey == "" {
        log.Fatal("API_KEY environment variable required")
    }

    fmt.Printf("Starting server on port %s\n", config.Port)
}
```

### .env Files (Manual Loading)

```go
// Simple .env loader
func loadEnvFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    for _, line := range strings.Split(string(data), "\n") {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            os.Setenv(parts[0], parts[1])
        }
    }
    return nil
}

// Usage
loadEnvFile(".env")
```

---

## Common Patterns

### Check Error and Return Early

```go
func doSomething() error {
    if err := step1(); err != nil {
        return fmt.Errorf("step1 failed: %w", err)
    }

    if err := step2(); err != nil {
        return fmt.Errorf("step2 failed: %w", err)
    }

    return nil
}
```

### Defer for Cleanup

```go
func readFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // always runs

    // work with file...
    return nil
}
```

### Comma-Ok Idiom

```go
// Map lookup
value, ok := myMap[key]

// Type assertion
str, ok := value.(string)

// Channel receive
msg, ok := <-ch
```

### Constructor Pattern

```go
type Server struct {
    host string
    port int
}

func NewServer(host string, port int) *Server {
    return &Server{
        host: host,
        port: port,
    }
}
```

### Functional Options Pattern

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func NewServer(host string, opts ...Option) *Server {
    s := &Server{host: host, port: 8080, timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer("localhost", WithPort(9000), WithTimeout(1*time.Minute))
```

---

## Quick Reference Table

| Operation | Syntax |
|-----------|--------|
| Declare variable | `var x int` or `x := 0` |
| Create slice | `s := make([]int, 5)` |
| Create map | `m := make(map[string]int)` |
| Create channel | `ch := make(chan int)` |
| Create buffered channel | `ch := make(chan int, 10)` |
| Get current time | `time.Now()` |
| Format time | `time.Now().Format("15:04:05")` |
| Sleep | `time.Sleep(1 * time.Second)` |
| Measure duration | `elapsed := time.Since(start)` |
| Start goroutine | `go func() { }()` |
| Send to channel | `ch <- value` |
| Receive from channel | `value := <-ch` |
| Close channel | `close(ch)` |
| Create error | `errors.New("message")` |
| Wrap error | `fmt.Errorf("context: %w", err)` |
| Convert int to string | `strconv.Itoa(42)` |
| Convert string to int | `strconv.Atoi("42")` |

---

## Debugging Concurrent Code

```go
import (
    "fmt"
    "time"
)

// Helper function for timestamped logging
func log(format string, args ...interface{}) {
    timestamp := time.Now().Format("15:04:05.000")
    fmt.Printf("[%s] "+format+"\n", append([]interface{}{timestamp}, args...)...)
}

// Usage
func main() {
    go func() {
        log("Goroutine 1: starting")
        time.Sleep(100 * time.Millisecond)
        log("Goroutine 1: done")
    }()

    go func() {
        log("Goroutine 2: starting")
        time.Sleep(50 * time.Millisecond)
        log("Goroutine 2: done")
    }()

    time.Sleep(200 * time.Millisecond)
}

// Output:
// [14:30:45.001] Goroutine 1: starting
// [14:30:45.001] Goroutine 2: starting
// [14:30:45.051] Goroutine 2: done
// [14:30:45.101] Goroutine 1: done
```
