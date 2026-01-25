# Learning Go for Backend Development

A structured learning path for mastering Go (Golang) with a focus on backend development.

---

## Project Overview

This repository tracks my journey learning Go for backend development. Go is ideal for backend services due to its simplicity, strong concurrency support, fast compilation, and excellent standard library.

**Goals:**
- Master Go fundamentals and idioms
- Build production-ready REST APIs
- Understand Go's concurrency model
- Deploy containerized Go applications

---

## Go Installation

### macOS
```bash
brew install go
```

### Linux
```bash
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Verify Installation
```bash
go version
go env GOPATH
```

### Environment Setup
Add to your shell profile (`~/.zshrc` or `~/.bashrc`):
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

---

## Learning Roadmap

### Phase 1: Go Basics
**Duration: 1-2 weeks**

| Topic | Resources |
|-------|-----------|
| Syntax & Variables | `var`, `:=`, constants, zero values |
| Data Types | `int`, `string`, `bool`, `float64`, arrays, slices, maps |
| Control Flow | `if`, `for`, `switch`, `defer` |
| Functions | Multiple returns, named returns, variadic functions |
| Structs & Methods | Custom types, receiver functions |
| Pointers | `&`, `*`, when to use pointers |
| Packages & Modules | `go mod init`, imports, visibility |
| Error Handling | `error` type, custom errors, `errors.Is/As` |

**Practice:** Build a CLI tool (todo list, file organizer)

---

### Phase 2: Concurrency
**Duration: 1-2 weeks**

| Topic | Description |
|-------|-------------|
| Goroutines | `go` keyword, lightweight threads |
| Channels | Unbuffered, buffered, directional |
| Select Statement | Multiplexing channel operations |
| Sync Package | `WaitGroup`, `Mutex`, `Once` |
| Context | Cancellation, timeouts, deadlines |
| Patterns | Worker pools, fan-in/fan-out, pipelines |

**Practice:** Build a concurrent web scraper or file processor

---

### Phase 3: Backend Fundamentals
**Duration: 2 weeks**

| Topic | Description |
|-------|-------------|
| net/http | `http.Handler`, `http.ServeMux`, request/response |
| Routing | Path parameters, query strings |
| Middleware | Logging, auth, CORS, recovery |
| JSON Handling | `encoding/json`, struct tags |
| Request Validation | Input sanitization, validation libraries |
| Templating | `html/template` for server-side rendering |

**Practice:** Build a basic CRUD API without frameworks

---

### Phase 4: Database Integration
**Duration: 2 weeks**

| Topic | Description |
|-------|-------------|
| database/sql | Connection pools, prepared statements |
| PostgreSQL/MySQL | `lib/pq`, `go-sql-driver/mysql` |
| Migrations | `golang-migrate`, `goose` |
| GORM | ORM basics, relationships, hooks |
| sqlx | Enhanced database/sql with struct scanning |
| Redis | Caching, sessions with `go-redis` |

**Practice:** Add persistence to your Phase 3 API

---

### Phase 5: Building REST APIs
**Duration: 2-3 weeks**

| Topic | Description |
|-------|-------------|
| API Design | RESTful conventions, versioning |
| Authentication | JWT, OAuth2, sessions |
| Authorization | RBAC, middleware-based permissions |
| Rate Limiting | Token bucket, sliding window |
| Documentation | Swagger/OpenAPI with `swag` |
| Pagination | Cursor-based, offset pagination |
| Error Responses | Consistent error formats |

**Practice:** Build a complete API (blog, e-commerce, task manager)

---

### Phase 6: Advanced Topics
**Duration: Ongoing**

| Topic | Description |
|-------|-------------|
| Testing | `testing` package, table-driven tests, mocks |
| Benchmarking | `go test -bench`, profiling with `pprof` |
| Docker | Multi-stage builds, minimal images |
| CI/CD | GitHub Actions, automated testing |
| Configuration | `viper`, environment variables, 12-factor app |
| Logging | `slog` (stdlib), `zerolog`, structured logging |
| Observability | Metrics, tracing, health checks |
| gRPC | Protocol buffers, service definitions |

---

## Project Structure

Standard layout for Go backend projects:

```
myapp/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loading
│   ├── handlers/
│   │   ├── user.go           # HTTP handlers
│   │   └── product.go
│   ├── middleware/
│   │   ├── auth.go           # Authentication middleware
│   │   └── logging.go
│   ├── models/
│   │   └── user.go           # Data models/entities
│   ├── repository/
│   │   └── user_repo.go      # Database operations
│   └── services/
│       └── user_service.go   # Business logic
├── pkg/
│   └── utils/                # Shared utilities (importable)
├── migrations/
│   └── 001_create_users.sql
├── scripts/
│   └── setup.sh
├── .env.example
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

**Key Principles:**
- `cmd/` - Entry points for executables
- `internal/` - Private application code (not importable by other projects)
- `pkg/` - Public libraries (importable by other projects)
- Keep `main.go` minimal - wire dependencies and start server

---

## Useful Commands

### Module Management
```bash
go mod init github.com/username/project   # Initialize module
go mod tidy                               # Clean up dependencies
go get package@version                    # Add/update dependency
go mod vendor                             # Vendor dependencies
```

### Building & Running
```bash
go run .                                  # Run current package
go run ./cmd/api                          # Run specific package
go build -o bin/app ./cmd/api             # Build binary
go install                                # Install to GOPATH/bin
```

### Testing
```bash
go test ./...                             # Run all tests
go test -v ./...                          # Verbose output
go test -cover ./...                      # With coverage
go test -race ./...                       # Race detector
go test -bench=. ./...                    # Run benchmarks
```

### Code Quality
```bash
go fmt ./...                              # Format code
go vet ./...                              # Static analysis
golangci-lint run                         # Comprehensive linting
```

### Profiling & Debugging
```bash
go tool pprof http://localhost:6060/debug/pprof/profile
go test -cpuprofile=cpu.out -memprofile=mem.out
```

---

## Popular Frameworks & Libraries

### Web Frameworks

| Framework | Description | Best For |
|-----------|-------------|----------|
| **Gin** | Fast, minimalist | Production APIs, performance-critical |
| **Echo** | High performance, extensible | REST APIs, middleware-heavy apps |
| **Fiber** | Express-inspired, fastest | Developers from Node.js background |
| **Chi** | Lightweight, stdlib compatible | Projects preferring stdlib patterns |

### Quick Comparison

```go
// Gin
r := gin.Default()
r.GET("/users/:id", getUser)
r.Run(":8080")

// Echo
e := echo.New()
e.GET("/users/:id", getUser)
e.Start(":8080")

// Fiber
app := fiber.New()
app.Get("/users/:id", getUser)
app.Listen(":8080")

// Chi
r := chi.NewRouter()
r.Get("/users/{id}", getUser)
http.ListenAndServe(":8080", r)
```

### Essential Libraries

| Category | Library | Purpose |
|----------|---------|---------|
| **ORM** | GORM | Full-featured ORM |
| **SQL** | sqlx | Enhanced database/sql |
| **Validation** | validator | Struct validation |
| **Config** | viper | Configuration management |
| **Logging** | zerolog, slog | Structured logging |
| **Auth** | jwt-go | JWT tokens |
| **Testing** | testify | Assertions and mocks |
| **HTTP Client** | resty | REST client |
| **CLI** | cobra | CLI applications |

---

## Quick Start Template

```go
// cmd/api/main.go
package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}
```

Run with: `go run ./cmd/api`

---

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Wiki](https://github.com/golang/go/wiki)
- [Awesome Go](https://awesome-go.com/)
