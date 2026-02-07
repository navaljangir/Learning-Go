# Makefile Explained: Complete Guide

## What is a Makefile?

A **Makefile** is a special file used by the `make` command (a build automation tool) to run tasks. Think of it like `package.json` scripts in Node.js, but much more powerful.

### Node.js vs Make Comparison

**Node.js (package.json):**
```json
{
  "scripts": {
    "start": "node server.js",
    "build": "webpack build",
    "test": "jest"
  }
}
```
Run with: `npm run start`

**Make (Makefile):**
```makefile
start:
	node server.js

build:
	webpack build

test:
	jest
```
Run with: `make start`

### Why Use Makefiles in Go?

1. **Language-agnostic** - Works on any system with `make` installed
2. **No npm required** - Go developers may not have Node.js installed
3. **Dependency management** - Targets can depend on other targets
4. **Conditional logic** - Built-in shell scripting support
5. **Standard tool** - Been around since 1976, well-tested and reliable

---

## Makefile Syntax: Every Keyword Explained

### 1. Basic Structure

```makefile
target: dependencies
	command
	another-command
```

- **target** - Name of the task (like `build`, `run`, `test`)
- **dependencies** - Other targets that must run first (optional)
- **command** - Shell command to execute (MUST start with a TAB, not spaces!)

### 2. `.PHONY` Keyword

```makefile
.PHONY: clean build
```

**What it does:**
- Tells `make` these are **not actual files**
- Normally, `make` checks if a file named `target` exists and is up-to-date
- `.PHONY` says "always run this, even if a file with this name exists"

**Why it matters:**

```bash
# Without .PHONY
$ touch clean          # Creates a file named "clean"
$ make clean           # Make thinks "clean" is up-to-date, does nothing!

# With .PHONY
$ touch clean          # Creates a file named "clean"
$ make clean           # Runs anyway because .PHONY tells it to ignore files
```

**In our Makefile:**
```makefile
.PHONY: help build run dev migrate-up migrate-down ...
```
All our targets are commands, not files, so we declare them all as `.PHONY`.

---

### 3. `@` Symbol (Silence Prefix)

```makefile
build:
	@echo "Building..."    # Hides the command itself
	go build               # Shows this command
```

**Without `@`:**
```bash
$ make build
echo "Building..."       # ← Shows the command
Building...              # ← Shows the output
go build                 # ← Shows the command
```

**With `@`:**
```bash
$ make build
Building...              # ← Only shows the output
```

**Why use it?**
- Makes output cleaner and more user-friendly
- Hides implementation details
- Shows only what matters to the user

---

### 4. Comments

```makefile
# This is a comment (starts with #)
build: ## This is a help comment (after ##)
	@echo "Building..."  # This is also a comment
```

- `#` - Regular comment (ignored by `make`)
- `##` - Special marker for help text (our Makefile extracts these)

---

### 5. Variables

```makefile
# Define variables
APP_NAME=todo-app
BUILD_DIR=bin

# Use variables with $(VAR_NAME)
build:
	go build -o $(BUILD_DIR)/$(APP_NAME)
```

**In our Makefile, we use command-line variables:**
```makefile
migrate-create:
	@if [ -z "$(NAME)" ]; then
		echo "Error: NAME is required"
		exit 1
	fi
	bash scripts/migrate.sh create $(NAME)
```

Usage: `make migrate-create NAME=add_users_table`

---

### 6. Special Makefile Variables

- `$(MAKEFILE_LIST)` - List of all Makefiles being processed
- `$$1`, `$$2` - Arguments in shell commands (double `$` to escape)
- `$@` - The target name (not used in our Makefile)
- `$<` - First dependency (not used in our Makefile)

---

## Line-by-Line Breakdown: Our Makefile

### Line 1: `.PHONY` Declaration

```makefile
.PHONY: help build run dev migrate-up migrate-down migrate-create migrate-version setup clean deps lint fmt vet sqlc-generate sqlc-verify install-tools all
```

**Breakdown:**
- `.PHONY:` - Directive telling `make` these are command names, not files
- Lists all targets in the Makefile
- Ensures `make` always runs these commands, never tries to check file timestamps

---

### Lines 3-7: `help` Target

```makefile
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
```

**Breakdown:**

1. `help:` - Target name
2. `## Show this help message` - Description for help output
3. `@echo 'Usage: make [target]'` - Print usage instructions
4. `@awk '...' $(MAKEFILE_LIST)` - **Magic line!**

**How the AWK command works:**

```awk
BEGIN {FS = ":.*?## "}
```
- `BEGIN` - Run before processing lines
- `FS = ":.*?## "` - Field Separator: Split on `:` followed by `##`
- Example: `build: ## Build the app` becomes:
  - Field 1 (`$$1`): `build`
  - Field 2 (`$$2`): `Build the app`

```awk
/^[a-zA-Z_-]+:.*?## /
```
- Regex pattern: Match lines like `build: ## Description`
- `^[a-zA-Z_-]+:` - Starts with target name and colon
- `.*?## ` - Has `##` somewhere after

```awk
{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}
```
- `printf` - Formatted print
- `\033[36m` - ANSI color code (cyan)
- `%-18s` - Left-align string, pad to 18 characters
- `\033[0m` - Reset color
- `$$1` - First field (target name) - **double `$` to escape in Makefile**
- `$$2` - Second field (description)

**Output:**
```bash
$ make help
Usage: make [target]

Available targets:
  build              Build the application
  run                Run the application
  dev                Run with hot reload (requires air)
```

---

### Lines 9-12: `build` Target

```makefile
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/todo-app ./cmd/api
	@echo "✓ Build complete! Binary: bin/todo-app"
```

**Breakdown:**

1. `build:` - Target name
2. `@echo "Building application..."` - Status message
3. `@go build -o bin/todo-app ./cmd/api` - **The actual build command**
   - `go build` - Compile Go code
   - `-o bin/todo-app` - Output binary to `bin/todo-app`
   - `./cmd/api` - Package to build (entry point)
4. `@echo "✓ Build complete! Binary: bin/todo-app"` - Success message

**How to use:**
```bash
make build
```

---

### Lines 14-16: `run` Target

```makefile
run: ## Run the application
	@echo "Starting application..."
	@go run ./cmd/api/main.go
```

**Breakdown:**

1. `@go run ./cmd/api/main.go` - Run Go application without building binary
   - `go run` - Compile and run in one step
   - `./cmd/api/main.go` - Entry point file

**Difference from `build`:**
- `build` - Creates a binary file (`bin/todo-app`)
- `run` - Compiles and runs immediately, no binary saved

**How to use:**
```bash
make run
```

---

### Lines 18-26: `dev` Target

```makefile
dev: ## Run with hot reload (requires air)
	@echo "Running with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Error: 'air' not found. Install it with:"; \
		echo "  go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi
```

**Breakdown:**

1. `@if command -v air > /dev/null; then \` - **Check if `air` is installed**
   - `command -v air` - Check if command exists
   - `> /dev/null` - Discard output (only check exit code)
   - `\` - Line continuation (command continues on next line)

2. `air;` - Run the `air` hot reload tool

3. `else` - If `air` not found:
   - Print error message
   - Show installation instructions
   - `exit 1` - Exit with error code

**What is `air`?**
- Hot reload tool for Go (like `nodemon` in Node.js)
- Watches files and restarts server on changes
- Install: `go install github.com/air-verse/air@latest`

**Why use backslash `\`?**
```makefile
# Without backslash (won't work)
if command -v air > /dev/null; then
	air

# With backslash (works)
if command -v air > /dev/null; then \
	air; \
fi
```
- Makefiles need entire shell script on one logical line
- `\` tells Make "command continues on next line"

---

### Lines 28-30: `migrate-up` Target

```makefile
migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@bash scripts/migrate.sh up
```

**Breakdown:**

1. `@bash scripts/migrate.sh up` - Execute external script
   - `bash` - Run script with bash shell
   - `scripts/migrate.sh` - Path to script file
   - `up` - Argument passed to script (runs "up" migrations)

**How to use:**
```bash
make migrate-up    # Apply new migrations to database
```

---

### Lines 32-34: `migrate-down` Target

```makefile
migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@bash scripts/migrate.sh down
```

**Breakdown:**

- Same as `migrate-up` but passes `down` argument
- Rolls back the most recent migration

**How to use:**
```bash
make migrate-down    # Undo the last migration
```

---

### Lines 36-43: `migrate-force` Target

```makefile
migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required"; \
		echo "Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@bash scripts/migrate.sh force $(VERSION)
```

**Breakdown:**

1. `@if [ -z "$(VERSION)" ]; then \` - **Check if VERSION variable is set**
   - `[ -z "$(VERSION)" ]` - Test if string is empty
   - `-z` - "Zero length" test
   - `$(VERSION)` - Variable from command line

2. If VERSION is empty:
   - Print error and usage instructions
   - `exit 1` - Exit with error

3. If VERSION is set:
   - `bash scripts/migrate.sh force $(VERSION)` - Force migration to specific version

**How to use:**
```bash
make migrate-force VERSION=3    # Force database to migration version 3
make migrate-force              # Error: VERSION is required
```

**Why force migrations?**
- When migrations are out of sync
- When you need to skip a broken migration
- When resetting migration state

---

### Lines 45-46: `migrate-version` Target

```makefile
migrate-version: ## Show current migration version
	@bash scripts/migrate.sh version
```

**Breakdown:**

- Calls script with `version` argument
- Shows current migration version number

**How to use:**
```bash
make migrate-version    # Output: 20240115123456 (current version)
```

---

### Lines 48-55: `migrate-create` Target

```makefile
migrate-create: ## Create a new migration (usage: make migrate-create NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required"; \
		echo "Usage: make migrate-create NAME=add_users_table"; \
		exit 1; \
	fi
	@bash scripts/migrate.sh create $(NAME)
```

**Breakdown:**

- Same pattern as `migrate-force`
- Requires `NAME` variable
- Creates new migration files with given name

**How to use:**
```bash
make migrate-create NAME=add_users_table
# Creates: 20240115123456_add_users_table.up.sql
#          20240115123456_add_users_table.down.sql
```

---

### Lines 57-59: `setup` Target

```makefile
setup: ## Setup database
	@echo "Setting up database..."
	@bash scripts/setup.sh
```

**Breakdown:**

- Runs initial database setup script
- Creates database, runs migrations, seeds data, etc.

**How to use:**
```bash
make setup    # First-time database setup
```

---

### Lines 61-66: `clean` Target

```makefile
clean: ## Clean build artifacts and cache
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "✓ Clean complete!"
```

**Breakdown:**

1. `@rm -rf bin/` - Remove build directory
   - `rm` - Remove command
   - `-r` - Recursive (delete directory and contents)
   - `-f` - Force (don't prompt, ignore errors if doesn't exist)

2. `@rm -f coverage.out coverage.html` - Remove coverage files
   - `-f` - Force (don't error if files don't exist)

3. `@go clean` - Go's built-in clean command
   - Removes object files and cached test results

**How to use:**
```bash
make clean    # Clean up all generated files
```

---

### Lines 68-72: `deps` Target

```makefile
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated!"
```

**Breakdown:**

1. `@go mod download` - Download all dependencies
   - Fetches packages to local cache
   - Doesn't modify `go.mod` or `go.sum`

2. `@go mod tidy` - Clean up dependencies
   - Adds missing dependencies
   - Removes unused dependencies
   - Updates `go.mod` and `go.sum`

**Node.js equivalent:**
```bash
npm install    # Similar to: go mod download && go mod tidy
```

**How to use:**
```bash
make deps    # Update dependencies
```

---

### Lines 74-77: `fmt` Target

```makefile
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted!"
```

**Breakdown:**

1. `@go fmt ./...` - Format all Go files
   - `go fmt` - Go's built-in formatter
   - `./...` - Recursive pattern (all packages in current directory and subdirectories)

**Node.js equivalent:**
```bash
npx prettier --write .    # Similar to: go fmt ./...
```

**How to use:**
```bash
make fmt    # Format all code
```

---

### Lines 79-87: `lint` Target

```makefile
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Error: 'golangci-lint' not found. Install it with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
```

**Breakdown:**

- Same pattern as `dev` target
- Checks if `golangci-lint` is installed
- Runs linter if available, shows install instructions if not

**What is `golangci-lint`?**
- Comprehensive Go linter (runs multiple linters at once)
- Like ESLint for Node.js

**How to use:**
```bash
make lint    # Run all linters
```

---

### Lines 89-92: `vet` Target

```makefile
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ No issues found!"
```

**Breakdown:**

1. `@go vet ./...` - Static analysis tool
   - Built into Go (no installation needed)
   - Finds suspicious code that compiles but may have bugs
   - Example catches: unreachable code, wrong printf formats, unused variables

**Difference from linter:**
- `go vet` - Built-in, catches probable bugs
- `golangci-lint` - Third-party, catches style issues, potential bugs, complexity

**How to use:**
```bash
make vet    # Run static analysis
```

---

### Lines 94-97: `sqlc-generate` Target

```makefile
sqlc-generate: ## Generate sqlc code from SQL queries
	@echo "Generating sqlc code..."
	@sqlc generate
	@echo "✓ sqlc code generated!"
```

**Breakdown:**

1. `@sqlc generate` - Generate Go code from SQL
   - Reads `sqlc.yaml` configuration
   - Processes `.sql` files in specified directories
   - Generates type-safe Go functions

**What is sqlc?**
- Tool that generates Go code from SQL
- Like Prisma for Node.js, but uses raw SQL
- Generates: structs, query functions, and type-safe interfaces

**How to use:**
```bash
make sqlc-generate    # Regenerate database code
```

---

### Lines 99-102: `sqlc-verify` Target

```makefile
sqlc-verify: ## Verify sqlc queries without generating code
	@echo "Verifying sqlc queries..."
	@sqlc verify
	@echo "✓ sqlc queries verified!"
```

**Breakdown:**

1. `@sqlc verify` - Check SQL queries without generating code
   - Validates SQL syntax
   - Checks against database schema
   - Faster than full generation

**How to use:**
```bash
make sqlc-verify    # Check SQL queries for errors
```

---

### Lines 104-110: `install-tools` Target

```makefile
install-tools: ## Install development tools (MySQL version)
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "✓ Tools installed!"
```

**Breakdown:**

1. `@go install github.com/air-verse/air@latest` - Install hot reload tool
   - `go install` - Install Go binary to `$GOPATH/bin`
   - `@latest` - Get latest version

2. `@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` - Install linter

3. `@go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
   - `-tags 'mysql'` - Build with MySQL support
   - Migration tool for database version control

4. `@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` - Install sqlc

**What is `-tags`?**
- Build tags - conditional compilation
- `-tags 'mysql'` - Include MySQL driver in the build
- Without it, only PostgreSQL might be supported

**How to use:**
```bash
make install-tools    # Install all development tools
```

---

### Lines 112-113: `all` Target

```makefile
all: deps fmt vet build ## Run all checks and build
	@echo "✓ All checks passed!"
```

**Breakdown:**

1. `all: deps fmt vet build` - **Target with dependencies**
   - `all` - Target name
   - `deps fmt vet build` - Other targets to run first
   - Make runs them in order: `deps` → `fmt` → `vet` → `build`

**How dependencies work:**
```makefile
all: deps build
```
When you run `make all`:
1. Make checks if `deps` is up-to-date (runs it)
2. Make checks if `build` is up-to-date (runs it)
3. Then runs commands in `all` target

**How to use:**
```bash
make all    # Run complete build pipeline
```

---

## Common Make Commands

### Basic Usage

```bash
make                # Runs first target (help in our case)
make build          # Run specific target
make help           # Show all available targets
```

### Multiple Targets

```bash
make clean build    # Run multiple targets in sequence
```

### Variables

```bash
make migrate-create NAME=add_users_table    # Pass variables
make migrate-force VERSION=3                # Multiple variables possible
```

### Debugging

```bash
make -n build       # Dry run (show commands without executing)
make -d build       # Debug mode (verbose output)
```

---

## Shell Commands Reference

### Conditionals in Makefile

```makefile
target:
	@if [ -z "$(VAR)" ]; then \
		echo "Variable is empty"; \
	else \
		echo "Variable is: $(VAR)"; \
	fi
```

**Test operators:**
- `-z STRING` - True if string is empty
- `-n STRING` - True if string is not empty
- `-f FILE` - True if file exists
- `-d DIR` - True if directory exists
- `-x FILE` - True if file is executable

### Command Checks

```makefile
@if command -v tool > /dev/null; then \
	tool; \
else \
	echo "Tool not found"; \
fi
```

- `command -v tool` - Check if command exists
- `> /dev/null` - Discard output (check exit code only)

---

## Common Patterns

### Pattern 1: Check and Install

```makefile
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi
```

### Pattern 2: Variable Validation

```makefile
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME required"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
```

### Pattern 3: Cleanup Pattern

```makefile
clean:
	@rm -rf bin/
	@rm -f *.out *.html
	@go clean
```

### Pattern 4: Multi-step Build

```makefile
build: deps fmt vet
	@go build -o bin/app ./cmd/api
```

---

## Tips and Best Practices

### 1. Always Use Tabs (Not Spaces!)

```makefile
# WRONG (spaces)
build:
    go build ./cmd/api

# CORRECT (tab)
build:
	go build ./cmd/api
```

**Why?** Make requires tabs for commands. Spaces will cause:
```
Makefile:2: *** missing separator. Stop.
```

### 2. Use `.PHONY` for All Command Targets

```makefile
.PHONY: build clean test

build:
	go build

clean:
	rm -rf bin/
```

### 3. Use `@` for Cleaner Output

```makefile
# Without @
build:
	echo "Building..."
	go build

# Output:
# echo "Building..."
# Building...
# go build

# With @
build:
	@echo "Building..."
	@go build

# Output:
# Building...
```

### 4. Provide Help Comments

```makefile
build: ## Build the application
run: ## Run the application
```

### 5. Use Line Continuation for Long Commands

```makefile
build:
	@go build \
		-o bin/app \
		-ldflags="-s -w" \
		./cmd/api
```

### 6. Validate Required Variables

```makefile
deploy:
	@if [ -z "$(ENV)" ]; then \
		echo "Error: ENV required (dev/prod)"; \
		exit 1; \
	fi
	@echo "Deploying to $(ENV)..."
```

---

## Troubleshooting

### Error: "missing separator"

**Cause:** Using spaces instead of tabs

**Fix:** Replace spaces with actual tab characters

### Error: "No rule to make target"

**Cause:** Typo in target name

**Fix:**
```bash
make biuld     # Wrong
make build     # Correct
```

### Error: Command not found

**Cause:** Tool not installed

**Fix:**
```bash
make install-tools    # Install required tools
```

### Variable not working

```makefile
# Wrong
make migrate-create name=test    # Lowercase

# Correct
make migrate-create NAME=test    # Must match $(NAME) in Makefile
```

---

## Summary: Our Makefile Quick Reference

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make build` | Build binary to `bin/todo-app` |
| `make run` | Run application without building |
| `make dev` | Run with hot reload (requires air) |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-version` | Show current migration version |
| `make migrate-create NAME=x` | Create new migration |
| `make migrate-force VERSION=x` | Force migration version |
| `make setup` | Initial database setup |
| `make clean` | Remove build artifacts |
| `make deps` | Update dependencies |
| `make fmt` | Format all Go code |
| `make lint` | Run golangci-lint |
| `make vet` | Run go vet static analysis |
| `make sqlc-generate` | Generate code from SQL |
| `make sqlc-verify` | Verify SQL queries |
| `make install-tools` | Install dev tools |
| `make all` | Run full build pipeline |

---

## Node.js Developers: Make vs NPM Scripts

### Feature Comparison

| Feature | NPM Scripts | Make |
|---------|-------------|------|
| **Setup** | Comes with Node.js | May need to install |
| **Syntax** | JSON | Makefile DSL |
| **Conditionals** | Need external tools | Built-in shell |
| **Dependencies** | Sequential with `&&` | Native target dependencies |
| **Variables** | Process.env | Command-line or internal |
| **Cross-platform** | Better (if using cross-env) | May need adjustments |

### When to Use Each

**Use NPM scripts when:**
- Building frontend projects
- Working in Node.js ecosystem
- Need cross-platform compatibility
- Team is familiar with package.json

**Use Make when:**
- Building backend Go projects
- Need complex build logic
- Want language-agnostic tooling
- Need dependency management between tasks

---

## Advanced: How Make Actually Works

### 1. Dependency Resolution

```makefile
all: deps build

deps:
	go mod download

build: fmt
	go build

fmt:
	go fmt ./...
```

**Execution order:**
```
make all
  └─> deps (runs first)
  └─> build
       └─> fmt (runs before build)
```

### 2. File Timestamps (Why .PHONY Matters)

```makefile
# Without .PHONY
build:
	go build -o build ./cmd/api

# If file named "build" exists:
$ make build
make: 'build' is up to date.    # Doesn't run!
```

Make checks:
1. Does file `build` exist?
2. Are source files newer than `build`?
3. If no, skip (thinks it's already built)

```makefile
# With .PHONY
.PHONY: build

build:
	go build -o build ./cmd/api

$ make build    # Always runs!
```

### 3. Variable Expansion

```makefile
NAME=test

create:
	echo $(NAME)           # Expanded by Make
	echo $$PATH            # Expanded by shell ($$=$ in shell)
```

---

## Conclusion

You now understand:
- ✓ What a Makefile is and why it's used
- ✓ Every keyword and symbol (`.PHONY`, `@`, `\`, etc.)
- ✓ Every target in our Makefile
- ✓ How to use shell commands in Make
- ✓ Common patterns and best practices
- ✓ How to debug Make issues

**Next steps:**
1. Try running `make help`
2. Experiment with `make build` and `make run`
3. Create your own custom targets
4. Read the man page: `man make`
