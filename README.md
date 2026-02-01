# Learning Go for Backend Development

A structured learning path for mastering Go (Golang) with a focus on backend development.

---

## 📚 Documentation

Comprehensive Go reference documentation is available in the [`learn/docs/`](./learn/docs) directory:

| File | Description |
|------|-------------|
| [00_go_syntax.md](./learn/docs/00_go_syntax.md) | **Complete Go syntax reference** - Variables, types, control flow, functions, data structures, concurrency, standard library, and common patterns |
| [01_packages.md](./learn/docs/01_packages.md) | Package management and organization |
| [02_hot_reload_air.md](./learn/docs/02_hot_reload_air.md) | Hot reload setup with Air |
| [03_go_get_vs_install.md](./learn/docs/03_go_get_vs_install.md) | Understanding `go get` vs `go install` |
| [04_file_naming.md](./learn/docs/04_file_naming.md) | Go file naming conventions |
| [05_go_deep_concepts.md](./learn/docs/05_go_deep_concepts.md) | Advanced Go concepts |
| [06_concurrency.md](./learn/docs/06_concurrency.md) | **Comprehensive concurrency guide** - Goroutines, channels, context, sync package, real API examples, and memory leak detection |

**Start with:**
- [00_go_syntax.md](./learn/docs/00_go_syntax.md) - Complete syntax reference with 31 sections
- [06_concurrency.md](./learn/docs/06_concurrency.md) - Deep dive into Go's concurrency model with Node.js comparisons

---

## 🚀 Quick Start

### Learning Files Structure

```
learn/
├── docs/                          # 📚 Documentation (READ THESE!)
│   ├── 00_go_syntax.md           # Complete Go syntax reference
│   ├── 06_concurrency.md         # Concurrency deep dive
│   └── ...                       # Other guides
├── 00_commands_reference.md       # Go commands cheat sheet
├── 01_basics/                     # Variables, types, printing
│   ├── main.go
│   ├── buffer/                    # Channel buffers
│   ├── context/                   # Context examples
│   ├── mutex/                     # Mutex examples
│   ├── pingpong/                  # Channel ping-pong
│   └── workerPool/                # Worker pool pattern
├── 02_functions/                  # Functions, parameters, returns
├── 03_format_specifiers/          # Printf, %s, %d, \n explained
├── 04_simple_server/              # First HTTP server
├── 05_go_concepts/                # Structs, pointers, errors
└── 06_gin_server/                 # Gin framework examples
```

**Run any example:**
```bash
cd learn/01_basics && go run main.go
cd learn/01_basics/workerPool && go run main.go
```

---

## 🎯 Project Goals

- ✅ Master Go fundamentals and idioms
- ✅ Build production-ready REST APIs
- ✅ Understand Go's concurrency model (goroutines, channels, context)
- 🔄 Deploy containerized Go applications

---

## 💡 Key Concepts (Coming from Node.js)

| Concept | In JavaScript | In Go |
|---------|---------------|-------|
| Objects | `{ name: "x" }` | Structs |
| Classes | `class User {}` | Structs + methods |
| Async/await | Automatic event loop | Goroutines + channels |
| null/undefined | Both exist | Only `nil` |
| Types | Dynamic (`let x`) | Static (`var x int`) |
| Errors | try/catch | Return `error` value |
| Pointers | Hidden | Explicit (`*`, `&`) |
| Promises | Single value | Channels (multiple values) |

---

## 🛠️ Installation

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

### Windows
Download from [go.dev/dl](https://go.dev/dl/) and run the installer.

### Verify Installation
```bash
go version
go env GOPATH
```

---

## 📖 Learning Path

### Phase 1: Go Basics ✅
**Topics:** Variables, types, control flow, functions, structs, pointers, packages, error handling

**Documentation:**
- [00_go_syntax.md](./learn/docs/00_go_syntax.md) - Complete syntax reference

**Practice:** ✅ Built CLI tools and basic programs

---

### Phase 2: Concurrency ✅
**Topics:** Goroutines, channels, select, sync package, context, worker pools

**Documentation:**
- [06_concurrency.md](./learn/docs/06_concurrency.md) - Complete concurrency guide
  - Goroutines vs Node.js event loop
  - Channels (buffered vs unbuffered)
  - Context for cancellation and timeouts
  - Memory leaks vs resource leaks
  - Real-world API examples
  - Worker pool patterns

**Practice:** ✅ Built concurrent examples (worker pools, rate limiting, ping-pong)

---

### Phase 3: Backend Fundamentals 🔄
**Topics:** net/http, routing, middleware, JSON handling, validation

**Practice:** Build a basic CRUD API without frameworks

---

### Phase 4: Database Integration
**Topics:** database/sql, PostgreSQL/MySQL, migrations, GORM, Redis

---

### Phase 5: Building REST APIs
**Topics:** Authentication, authorization, rate limiting, documentation, error handling

---

### Phase 6: Production & Deployment
**Topics:** Testing, Docker, CI/CD, logging, observability, gRPC

---

## 🔧 Essential Commands

### Module Management
```bash
go mod init github.com/username/project   # Initialize module
go mod tidy                               # Sync dependencies
go get github.com/pkg                     # Add package
```

### Building & Running
```bash
go run .                                  # Run current package
go build -o bin/app ./cmd/api             # Build binary
```

### Testing & Quality
```bash
go test ./...                             # Run all tests
go test -race ./...                       # Race detector
go fmt ./...                              # Format code
go vet ./...                              # Static analysis
```

---

## 📚 Resources

- [Go Documentation](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [This Project's Docs](./learn/docs/) - **Start here!**

---

## 🏗️ Project Structure

```
Learning-Go/
├── README.md                      # You are here
├── CLAUDE.md                      # Claude Code instructions
├── learn/
│   ├── docs/                      # 📚 Main documentation
│   │   ├── 00_go_syntax.md       # Complete syntax reference
│   │   ├── 06_concurrency.md     # Concurrency deep dive
│   │   └── ...                   # Other guides
│   ├── 01_basics/                 # Basic examples
│   ├── 02_functions/              # Function examples
│   └── ...
├── cmd/                           # Application entry points
├── internal/                      # Private application code
├── pkg/                           # Public libraries
└── go.mod                         # Module definition
```

---

## 📝 License

This is a personal learning project. Feel free to use any code or documentation for your own learning!
