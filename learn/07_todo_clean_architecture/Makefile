.PHONY: help build run dev migrate-up migrate-down migrate-create migrate-version setup clean deps lint fmt vet sqlc-generate sqlc-verify install-tools all

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

sqlc-generate: ## Generate sqlc code from SQL queries
	@echo "Generating sqlc code..."
	@sqlc generate
	@echo "✓ sqlc code generated!"

sqlc-verify: ## Verify sqlc queries without generating code
	@echo "Verifying sqlc queries..."
	@sqlc verify
	@echo "✓ sqlc queries verified!"

install-tools: ## Install development tools (MySQL version)
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "✓ Tools installed!"

all: deps fmt vet build ## Run all checks and build
	@echo "✓ All checks passed!"
