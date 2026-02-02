.PHONY: help build run test test-coverage migrate-up migrate-down migrate-create migrate-version setup clean deps lint fmt

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/todo-app ./cmd/api
	@echo "✓ Build complete! Binary: bin/todo-app"

run: ## Run the application
	@echo "Starting application..."
	@go run ./cmd/api/main.go

dev: ## Run with hot reload (requires air)
	@echo "Running with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Error: 'air' not found. Install it with:"; \
		echo "  go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi

test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test ./... -race -v

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@bash scripts/migrate.sh up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@bash scripts/migrate.sh down

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required"; \
		echo "Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@bash scripts/migrate.sh force $(VERSION)

migrate-version: ## Show current migration version
	@bash scripts/migrate.sh version

migrate-create: ## Create a new migration (usage: make migrate-create NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required"; \
		echo "Usage: make migrate-create NAME=add_users_table"; \
		exit 1; \
	fi
	@bash scripts/migrate.sh create $(NAME)

setup: ## Setup database
	@echo "Setting up database..."
	@bash scripts/setup.sh

clean: ## Clean build artifacts and cache
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "✓ Clean complete!"

deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated!"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted!"

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Error: 'golangci-lint' not found. Install it with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ No issues found!"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✓ Tools installed!"

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	@docker-compose up -d
	@echo "✓ Services started!"

docker-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	@docker-compose down
	@echo "✓ Services stopped!"

docker-logs: ## View docker-compose logs
	@docker-compose logs -f

all: deps fmt vet test build ## Run all checks and build
	@echo "✓ All checks passed!"
