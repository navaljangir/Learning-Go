# TODO App with Clean Architecture

A production-ready TODO application built with Go, following clean architecture principles.

## Features

- **User Authentication** - Register, login with JWT tokens
- **User Profile** - View and update user profile
- **TODO Management** - Create, read, update, delete todos
- **Priority System** - Set priorities (low, medium, high, urgent)
- **Due Dates** - Set and track due dates
- **Completion Tracking** - Mark todos as complete/incomplete
- **Pagination** - Efficient pagination for todo lists
- **Soft Deletes** - Recover accidentally deleted data

## Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin Web Framework
- **Database:** MySQL 8.0 (runs in Docker)
- **Query Builder:** sqlc (type-safe SQL)
- **Driver:** /mysql (database/sql)
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt

## Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
├── domain/              # Domain layer (business entities and repository interfaces)
├── internal/            # Application layer (services, DTOs, repository implementations)
├── api/                 # Presentation layer (handlers, middleware, routes)
├── pkg/                 # Shared utilities (JWT, hashing, responses)
├── config/              # Configuration management
├── migrations/          # Database migrations
└── cmd/api/             # Application entry point
```

### Dependency Rule

- **Domain layer** has no dependencies (pure business logic)
- **Application layer** depends only on domain
- **Infrastructure** implements domain interfaces
- **Presentation** depends on application and domain

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker (for running MySQL)
- golang-migrate CLI (for migrations)
- sqlc (for generating type-safe SQL code)

### Installation

1. **Start MySQL container (ONE COMMAND):**
   ```bash
   docker run -d --name todo_mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=todo_db -p 3306:3306 mysql:8.0
   ```

2. **Wait for MySQL to start (10-15 seconds):**
   ```bash
   sleep 15
   ```

3. **Clone/navigate to the directory:**
   ```bash
   cd learn/07_todo_clean_architecture
   ```

4. **Setup environment variables:**
   ```bash
   cp .env.example .env
   # Already configured for Docker MySQL at localhost:3306
   ```

5. **Install dependencies:**
   ```bash
   go mod download
   ```

6. **Run migrations:**
   ```bash
   make migrate-up
   ```

7. **Run the application:**
   ```bash
   make run
   ```

The server will start on `http://localhost:8080`

> See `QUICKSTART.md` for quick copy-paste commands!

## API Endpoints

### Public Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `GET /health` - Health check

### Protected Endpoints (Require JWT Token)

- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users/profile` - Update user profile
- `GET /api/v1/todos` - List user's todos (with pagination)
- `POST /api/v1/todos` - Create a new todo
- `GET /api/v1/todos/:id` - Get a specific todo
- `PUT /api/v1/todos/:id` - Update a todo
- `PATCH /api/v1/todos/:id/toggle` - Toggle todo completion
- `DELETE /api/v1/todos/:id` - Delete a todo (soft delete)

## Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

### Create a TODO

```bash
TOKEN="your-jwt-token"
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "priority": "medium",
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

### List TODOs

```bash
curl -X GET "http://localhost:8080/api/v1/todos?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

## Makefile Commands

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make sqlc-generate` - Generate type-safe SQL code
- `make sqlc-verify` - Verify SQL queries
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create name=xxx` - Create a new migration
- `make setup` - Setup database directory
- `make clean` - Clean build artifacts
- `make deps` - Download dependencies
- `make install-tools` - Install all dev tools (sqlc, migrate, etc.)

## Project Structure

```
.
├── api/
│   ├── handler/           # HTTP request handlers
│   ├── middleware/        # HTTP middleware (auth, logging, etc.)
│   └── router/            # Route definitions
├── cmd/
│   └── api/
│       └── main.go        # Application entry point
├── config/
│   └── config.go          # Configuration management
├── domain/
│   ├── entity/            # Business entities
│   └── repository/        # Repository interfaces
├── internal/
│   ├── dto/               # Data Transfer Objects
│   ├── repository/        # Repository layer
│   │   ├── queries/       # SQL queries for sqlc
│   │   ├── sqlc/          # Generated Go code (DO NOT EDIT)
│   │   └── sqlc_impl/     # Repository implementations using sqlc
│   └── service/           # Business logic services
├── migrations/            # Database migrations
├── pkg/
│   ├── constants/         # Application constants
│   └── utils/             # Utility functions
├── scripts/
│   ├── migrate.sh         # Migration script
│   └── setup.sh           # Setup script
├── docs/                  # Documentation
│   ├── SQLC_MIGRATION.md        # sqlc migration guide
│   └── SQLC_QUICK_REFERENCE.md  # sqlc quick reference
├── .env.example           # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── sqlc.yaml              # sqlc configuration
├── Makefile
├── MIGRATION_SUMMARY.md   # Database layer migration summary
└── README.md
```

## Security Features

- **Password Hashing** - bcrypt with cost factor 10
- **JWT Authentication** - Stateless authentication with token expiry
- **SQL Injection Prevention** - Parameterized queries only
- **Input Validation** - Gin binding validators
- **Authorization** - User-specific resource access control
- **Soft Deletes** - Data recovery capability

## Why MySQL in Docker?

- **No Local Installation** - Don't need MySQL installed on your machine
- **Isolated** - Separate from any existing databases
- **Easy to Reset** - Just delete and recreate the container
- **Production-like** - Similar to real deployment scenarios
- **Simple Command** - One command to start everything

```bash
docker run -d --name todo_mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=todo_db -p 3306:3306 mysql:8.0
```

## Docker Commands

```bash
# Start MySQL
docker start todo_mysql

# Stop MySQL
docker stop todo_mysql

# View logs
docker logs todo_mysql

# Connect to MySQL shell
docker exec -it todo_mysql mysql -uroot -prootpassword todo_db

# Remove container (deletes data)
docker rm -f todo_mysql
```

See `MYSQL_DOCKER_SETUP.md` for detailed documentation.

## Contributing

This is a learning project demonstrating clean architecture principles in Go.

## License

MIT License
