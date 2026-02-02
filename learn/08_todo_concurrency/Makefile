.PHONY: help run dev build clean test deps race

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run: ## Run the application (no hot reload)
	@go run cmd/api/main.go

dev: ## Run with hot reload (watches for file changes)
	@echo "Starting with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ 'air' not found. Install it with:"; \
		echo "  go install github.com/air-verse/air@latest"; \
		echo ""; \
		echo "After installing, run 'make dev' again"; \
		exit 1; \
	fi

build: ## Build the application
	@echo "Building..."
	@go build -o bin/todo-api cmd/api/main.go
	@echo "✓ Build complete: bin/todo-api"

clean: ## Remove build artifacts
	@rm -rf bin/
	@echo "✓ Cleaned"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

test: ## Run tests
	@go test -v ./...

race: ## Run with race detector (detects concurrency bugs!)
	@echo "Running with race detector..."
	@echo "This will catch data races and concurrency bugs!"
	@go run -race cmd/api/main.go

fmt: ## Format code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

# Example curl commands for testing

.PHONY: test-create test-batch test-notify test-stats

test-create: ## Test creating a todo
	@echo "Creating a todo..."
	@curl -X POST http://localhost:8080/api/v1/todos \
		-H "Content-Type: application/json" \
		-d '{"title":"Learn Goroutines","description":"Master concurrent programming","priority":3}'
	@echo ""

test-batch: ## Test batch creation
	@echo "Creating batch of todos..."
	@curl -X POST http://localhost:8080/api/v1/todos/batch \
		-H "Content-Type: application/json" \
		-d '{"todos":[{"title":"Task 1","description":"First task","priority":2},{"title":"Task 2","description":"Second task","priority":1},{"title":"Task 3","description":"Third task","priority":3}]}'
	@echo ""

test-notify: ## Test notification (replace :id with actual todo ID)
	@echo "Sending notification..."
	@curl -X POST http://localhost:8080/api/v1/todos/1/notify \
		-H "Content-Type: application/json" \
		-d '{"message":"Don't forget this task!","delay_seconds":3}'
	@echo ""

test-stats: ## View statistics
	@echo "Fetching statistics..."
	@curl http://localhost:8080/api/v1/stats
	@echo ""
