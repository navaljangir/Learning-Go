# Go Commands & Syntax Reference

Quick reference for Go commands, syntax, and patterns.

---

## Module Commands (like npm)

| Command | What it does | npm equivalent |
|---------|--------------|----------------|
| `go mod init myapp` | Create new module | `npm init` |
| `go mod tidy` | Add missing, remove unused deps | `npm install` (cleanup) |
| `go get github.com/pkg` | Add a package | `npm install pkg` |
| `go get github.com/pkg@v1.2.3` | Add specific version | `npm install pkg@1.2.3` |
| `go get -u ./...` | Update all dependencies | `npm update` |
| `go mod download` | Download all dependencies | `npm ci` |
| `go mod vendor` | Copy deps to vendor folder | - |

---

## Build & Run Commands

| Command | What it does |
|---------|--------------|
| `go run main.go` | Compile and run (temp binary) |
| `go run .` | Run current package |
| `go run ./cmd/api` | Run specific package |
| `go build` | Compile to binary |
| `go build -o myapp` | Compile with custom name |
| `go install` | Compile and install to $GOPATH/bin |

### Cross-compilation
```bash
# Build for Linux from macOS
GOOS=linux GOARCH=amd64 go build -o myapp-linux

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# Common OS/ARCH combinations:
# darwin/amd64  - macOS Intel
# darwin/arm64  - macOS Apple Silicon
# linux/amd64   - Linux 64-bit
# windows/amd64 - Windows 64-bit
```

---

## Testing Commands

| Command | What it does |
|---------|--------------|
| `go test` | Run tests in current package |
| `go test ./...` | Run ALL tests in project |
| `go test -v` | Verbose output |
| `go test -run TestName` | Run specific test |
| `go test -cover` | Show coverage % |
| `go test -coverprofile=c.out` | Save coverage data |
| `go tool cover -html=c.out` | View coverage in browser |
| `go test -race` | Detect race conditions |
| `go test -bench=.` | Run benchmarks |

---

## Installing Global Tools (like npm install -g)

```bash
# Install a CLI tool globally
go install github.com/air-verse/air@latest           # Hot reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linter
go install github.com/swaggo/swag/cmd/swag@latest    # Swagger generator
```

**Where do tools get installed?**
```bash
go env GOPATH    # Shows: C:\Users\yourname\go
# Tools go to: C:\Users\yourname\go\bin\
```

**Add to PATH (one time setup):**

| OS | How to add Go bin to PATH |
|----|---------------------------|
| **Windows** | Add `%USERPROFILE%\go\bin` to System Environment Variables |
| **macOS/Linux** | Add `export PATH=$PATH:$(go env GOPATH)/bin` to `~/.zshrc` or `~/.bashrc` |

**Windows step-by-step:**
1. `Win + R` → type `sysdm.cpl` → Enter
2. Advanced tab → Environment Variables
3. Edit `Path` under User variables → Add `%USERPROFILE%\go\bin`
4. Restart your terminal

---

## Code Quality Commands

| Command | What it does |
|---------|--------------|
| `go fmt ./...` | Format all code |
| `go vet ./...` | Find suspicious code |
| `gofmt -s -w .` | Format + simplify |
| `goimports -w .` | Format + fix imports |
| `golangci-lint run` | Run many linters (install separately) |

---

## Info Commands

| Command | What it does |
|---------|--------------|
| `go version` | Show Go version |
| `go env` | Show all Go environment variables |
| `go env GOPATH` | Show specific variable |
| `go list -m all` | List all dependencies |
| `go doc fmt.Println` | Show documentation |
| `go doc -all fmt` | Full package docs |

---

## go.mod File Explained

```go
module github.com/tejas/learningGo  // Module path (unique identifier)

go 1.21                              // Minimum Go version

require (
    github.com/gin-gonic/gin v1.9.1  // Direct dependency
    github.com/lib/pq v1.10.9        // Another dependency
)

require (
    // indirect = used by your dependencies, not by you directly
    golang.org/x/sys v0.5.0 // indirect
)
```

### Module Path Naming
```bash
# For GitHub projects
go mod init github.com/username/projectname

# For local/private projects (anything works)
go mod init myapp
go mod init company.com/myapp
```

---

## go.sum File

- Auto-generated checksum file
- Ensures everyone gets exact same dependency versions
- Like `package-lock.json`
- **Don't edit manually, commit to git**

---

## Printf Format Specifiers

| Specifier | Type | Example |
|-----------|------|---------|
| `%s` | string | `"hello"` |
| `%d` | integer | `42` |
| `%f` | float | `3.14159` |
| `%.2f` | float (2 decimals) | `3.14` |
| `%t` | boolean | `true` |
| `%v` | any value | anything |
| `%+v` | struct with field names | `{Name:Tejas Age:25}` |
| `%#v` | Go syntax representation | `User{Name:"Tejas"}` |
| `%T` | type of value | `int`, `string` |
| `%p` | pointer address | `0xc0000...` |
| `%%` | literal % | `%` |

### Escape Characters
| Char | Meaning |
|------|---------|
| `\n` | newline |
| `\t` | tab |
| `\\` | backslash |
| `\"` | quote |

---

## Variable Declaration

```go
// Explicit type
var name string = "Tejas"
var age int = 25

// Type inference (Go figures it out)
var name = "Tejas"

// Short declaration (most common, only in functions)
name := "Tejas"
age := 25

// Multiple variables
var a, b, c int = 1, 2, 3
x, y := 10, 20

// Constants
const PI = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)
```

---

## Pointers Quick Reference

```go
x := 10

&x      // "address of" x → returns pointer
*p      // "value at" pointer p → returns value

p := &x   // p is pointer to x
*p = 20   // change value through pointer
fmt.Println(x)  // x is now 20
```

### When to Use Pointers
- Modify original value in function
- Avoid copying large structs
- Indicate "optional" (nil = no value)

---

## Error Handling Pattern

```go
result, err := someFunction()
if err != nil {
    // handle error
    return err
}
// use result
```

---

## Project Structure

```
myapp/
├── cmd/           # Entry points (main packages)
│   └── api/
│       └── main.go
├── internal/      # Private code (can't be imported)
│   ├── handlers/
│   ├── models/
│   └── services/
├── pkg/           # Public libraries (can be imported)
├── go.mod
└── go.sum
```

---

## Common Gotchas

1. **Unused variables = compile error**
   ```go
   x := 5  // ERROR if x is never used
   ```

2. **Unused imports = compile error**
   ```go
   import "fmt"  // ERROR if fmt is never used
   ```

3. **`:=` only works inside functions**
   ```go
   package main
   x := 5  // ERROR at package level
   var x = 5  // OK
   ```

4. **Public vs Private = Capital letter**
   ```go
   func DoSomething() {}  // Public (exported)
   func doSomething() {}  // Private (unexported)
   ```
