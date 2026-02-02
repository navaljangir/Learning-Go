# Development Tools Explained

Understanding Makefile, Air, and other dev tools in Go.

---

## Quick Answer

| Tool | Purpose | Auto-reload? |
|------|---------|-------------|
| **Makefile** | Task runner (shortcuts for commands) | ❌ No |
| **Air** | Hot reload tool (watches file changes) | ✅ Yes |
| **go run** | Compile & run directly | ❌ No |
| **go build** | Create binary executable | ❌ No |

---

## 1. Makefile - Task Runner

### What it does
Makefile is just a **shortcut creator**. It doesn't watch files or auto-reload anything.

### Example
```makefile
run: ## Run the application
	@go run cmd/api/main.go

dev: ## Run with hot reload
	@air
```

When you type `make run`, it just executes `go run cmd/api/main.go`.

### Workflow WITHOUT hot reload
```
1. make run          # Server starts
2. Edit code
3. Ctrl+C            # Stop server manually
4. make run          # Start again manually
5. Repeat...         # Tedious!
```

---

## 2. Air - Hot Reload Tool

### What it does
Air **watches your Go files** and automatically:
1. Detects file changes
2. Stops the running server
3. Recompiles your code
4. Restarts the server

### Workflow WITH hot reload (air)
```
1. make dev          # Server starts with air
2. Edit code
3. Save file         # Air automatically detects change
4. Air restarts      # Server reloads automatically!
5. Repeat...         # No manual restart needed!
```

### Installation
```bash
# Install air globally
go install github.com/air-verse/air@latest

# Verify installation
air -v
```

### Configuration
Create `.air.toml` in project root:
```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]
```

This tells air:
- Which command to build your app
- Where to put the binary
- Which files to watch (.go files)
- Which directories to ignore

---

## 3. Comparison: make run vs make dev

### `make run` (No hot reload)

```bash
make run
```

**What happens:**
```
1. Runs: go run cmd/api/main.go
2. Server starts
3. If you edit code... nothing happens
4. You must Ctrl+C and run again
```

**Use when:**
- Quick one-time test
- Production deployment
- You don't plan to edit code

---

### `make dev` (With air hot reload)

```bash
make dev
```

**What happens:**
```
1. Runs: air (with .air.toml config)
2. Air compiles and starts server
3. Air watches all .go files
4. When you save a file:
   ├─ Air detects change
   ├─ Air stops old server
   ├─ Air recompiles code
   └─ Air starts new server
5. All automatic! 🎉
```

**Use when:**
- Active development
- Testing frequently
- Learning and experimenting

---

## 4. Other Common Tools

### nodemon (Node.js equivalent)
If you're from Node.js background:
```bash
# Node.js
nodemon server.js

# Go equivalent
air
```

Both watch files and auto-restart!

### entr (Unix tool)
Alternative to air:
```bash
# Watch files and run command on change
ls **/*.go | entr -r go run cmd/api/main.go
```

### watchexec (Rust tool)
Another alternative:
```bash
watchexec -e go -r go run cmd/api/main.go
```

---

## 5. Your Project Setup

### Current Makefile Commands

```bash
make help         # Show all available commands
make run          # Run once (no hot reload)
make dev          # Run with hot reload (air)
make build        # Create binary in bin/
make race         # Run with race detector
make clean        # Remove build artifacts
make deps         # Download dependencies

# Testing commands
make test-create  # Test create endpoint
make test-batch   # Test batch endpoint
make test-notify  # Test notification endpoint
make test-stats   # View statistics
```

### Recommended Workflow

**For learning/development:**
```bash
# Terminal 1: Run server with hot reload
make dev

# Terminal 2: Test your changes
make test-create
curl http://localhost:8080/api/v1/todos
```

**For testing concurrency bugs:**
```bash
# Run with race detector
make race

# Make requests from multiple terminals
while true; do curl http://localhost:8080/api/v1/stats; done
```

---

## 6. Why Both Makefile AND Air?

### Makefile provides:
- Consistent commands across projects
- Documentation of available tasks
- Complex multi-step commands
- Project-specific shortcuts

### Air provides:
- Automatic hot reload
- Faster development cycle
- Less context switching
- Immediate feedback

### Together:
```makefile
dev: ## Run with hot reload
	@air

build: ## Build production binary
	@go build -o bin/app cmd/api/main.go

test: ## Run tests
	@go test ./...
```

One Makefile, multiple tools orchestrated!

---

## 7. Real-World Development Cycle

### Without Air (Painful)
```
1. make run
2. Test endpoint... bug found
3. Ctrl+C
4. Fix code
5. make run
6. Wait for compile...
7. Test endpoint... another bug
8. Ctrl+C
9. Fix code
10. make run
11. Repeat 50 times... 😫
```

### With Air (Smooth)
```
1. make dev
2. Test endpoint... bug found
3. Fix code, save file
4. [Air auto-restarts]
5. Test endpoint immediately
6. Repeat as needed... 😊
```

**Time saved:** Seconds per iteration × 50 iterations = Minutes per session!

---

## 8. Advanced: Custom Air Config

### Basic config (.air.toml)
```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  include_ext = ["go", "html", "tmpl"]
  exclude_dir = ["tmp", "vendor", "bin"]
  delay = 1000  # Wait 1s before rebuild
```

### With environment variables
```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  full_bin = "APP_ENV=dev ./tmp/main"  # Pass env vars
```

### With pre/post commands
```toml
[build]
  pre_cmd = ["echo 'Building...'"]
  post_cmd = ["echo 'Build complete!'"]
```

---

## 9. Troubleshooting

### Air not found
```bash
# Install it
go install github.com/air-verse/air@latest

# Add to PATH (if needed)
export PATH=$PATH:$(go env GOPATH)/bin
```

### Air not detecting changes
```bash
# Check .air.toml configuration
# Ensure include_ext has "go"
# Ensure your file isn't in exclude_dir
```

### Port already in use
```bash
# Find process on port 8080
lsof -i :8080

# Kill it
kill -9 <PID>
```

### Air restarts too frequently
```bash
# Increase delay in .air.toml
[build]
  delay = 2000  # Wait 2 seconds
```

---

## 10. Quick Reference

### Install Tools
```bash
# Air (hot reload)
go install github.com/air-verse/air@latest

# golangci-lint (linting)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# dlv (debugger)
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Common Commands
```bash
# Development
make dev              # Hot reload (best for coding)
make run              # Single run
make race             # Detect races

# Building
make build            # Production binary
go build -race        # Binary with race detector

# Testing
make test             # Run tests
go test -race ./...   # Tests with race detector
go test -bench=.      # Run benchmarks
```

---

## Summary

| Scenario | Command | Auto-reload? |
|----------|---------|-------------|
| **Active Development** | `make dev` | ✅ Yes (air) |
| **Quick Test** | `make run` | ❌ No |
| **Concurrency Testing** | `make race` | ❌ No |
| **Production Build** | `make build` | ❌ No |

**Bottom Line:**
- **Makefile** = Task shortcuts
- **Air** = Auto-reload magic
- Together = Happy developer 😊

---

## Your Next Steps

1. Install air:
   ```bash
   go install github.com/air-verse/air@latest
   ```

2. Run with hot reload:
   ```bash
   make dev
   ```

3. Edit some code and save - watch it auto-reload!

4. Compare the experience:
   - Try `make run`, edit code, manually restart
   - Try `make dev`, edit code, auto-reloads
   - Feel the difference!

Now you understand why we use both! 🚀
