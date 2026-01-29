# Go Packages Explained

## What is a Package?

A package is Go's way of organizing code into reusable modules (similar to Node.js modules).

## Two Types of Packages

### 1. `package main` - Entry Point

```go
// main.go
package main  // Tells Go: "This is an executable program"

func main() {  // Required! Program execution starts here
    fmt.Println("Hello!")
}
```

- **Purpose**: Entry point of your program
- **Required**: Must have `func main()`
- **Compiles to**: Executable binary (.exe)
- **Node.js equivalent**: Your `index.js` or entry file in `package.json`

### 2. `package <name>` - Library/Module

```go
// server/server.go
package server  // Tells Go: "This is a reusable module"

func NewServer() *Server {
    // ...
}
```

- **Purpose**: Reusable code that other packages import
- **Cannot run directly**: No main function
- **Node.js equivalent**: A module you `require()` or `import`

## Package Naming Rules

| Rule | Example |
|------|---------|
| Folder name = Package name | `server/` folder → `package server` |
| All files in folder share same package | `server/server.go` and `server/config.go` both use `package server` |
| Lowercase names | `package server` not `package Server` |
| Short, concise names | `utils`, `models`, `handlers` |

## Importing Packages

```go
import (
    // Standard library (comes with Go)
    "fmt"
    "net/http"

    // Your local packages (module_name/folder_path)
    "gin_server/server"
    "gin_server/utils"

    // Third-party packages
    "github.com/gin-gonic/gin"
)
```

## Exported vs Private

In Go, **capitalization** controls visibility:

```go
package server

func NewServer() {}   // Exported (Capital N) - accessible from other packages
func helper() {}      // Private (lowercase h) - only accessible within this package

var Port = 8080       // Exported
var maxConn = 100     // Private
```

**Node.js comparison:**
```javascript
// In Node.js you explicitly export
module.exports = { NewServer }  // exported
const helper = () => {}         // private (not exported)
```

**In Go, it's automatic based on first letter!**

## Project Structure Example

```
my_project/
├── main.go              → package main (entry point)
├── go.mod               → module my_project
├── server/
│   └── server.go        → package server
├── handlers/
│   ├── auth.go          → package handlers
│   └── user.go          → package handlers (same folder = same package)
├── models/
│   └── user.go          → package models
└── utils/
    ├── response.go      → package utils
    └── jwt.go           → package utils
```

## How to Use

```go
// main.go
package main

import (
    "my_project/server"   // Import by module_name/folder_path
    "my_project/utils"
)

func main() {
    s := server.NewServer()    // Use as: package_name.FunctionName()
    utils.HashPassword("123")
}
```

## Common Mistakes

### 1. Wrong package name
```go
// File: server/server.go
package srv  // WRONG - should match folder name "server"
```

### 2. Multiple packages in same folder
```go
// server/server.go
package server

// server/config.go
package config  // WRONG - all files in server/ must be "package server"
```

### 3. Trying to run a library package
```bash
$ go run server/server.go
# Error: package server is not a main package
```

## Quick Reference

| Concept | Go | Node.js |
|---------|-----|---------|
| Entry point | `package main` + `func main()` | `index.js` |
| Module/Library | `package mylib` | `module.exports` |
| Import | `import "path/to/pkg"` | `require()` / `import` |
| Export function | Start with Capital letter | Add to `module.exports` |
| Private function | Start with lowercase | Don't export it |
