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

- **Language:** Go 1.21
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL
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
- PostgreSQL 12 or higher
- golang-migrate CLI (for migrations)

### Installation

1. **Clone the repository**
   ```bash
   cd learn/07_todo_clean_architecture
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Create database**
   ```bash
   make setup
   ```

5. **Run migrations**
   ```bash
   make migrate-up
   ```

6. **Run the application**
   ```bash
   make run
   ```

The server will start on `http://localhost:8080`

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
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create name=xxx` - Create a new migration
- `make setup` - Setup database
- `make clean` - Clean build artifacts
- `make deps` - Download dependencies

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
│   ├── repository/        # Repository implementations
│   │   └── postgres/      # PostgreSQL implementations
│   └── service/           # Business logic services
├── migrations/            # Database migrations
├── pkg/
│   ├── constants/         # Application constants
│   └── utils/             # Utility functions
├── scripts/
│   ├── migrate.sh         # Migration script
│   └── setup.sh           # Setup script
├── .env.example           # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Security Features

- **Password Hashing** - bcrypt with cost factor 10
- **JWT Authentication** - Stateless authentication with token expiry
- **SQL Injection Prevention** - Parameterized queries only
- **Input Validation** - Gin binding validators
- **Authorization** - User-specific resource access control
- **Soft Deletes** - Data recovery capability

## Contributing

This is a learning project demonstrating clean architecture principles in Go.

## License

MIT License
