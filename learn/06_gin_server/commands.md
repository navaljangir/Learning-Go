# Gin Server Commands Reference

## Running the Server

```bash
# Run directly
go run .

# Build and run
go build -o server.exe . && ./server.exe

# Run with specific port (modify constants/constants.go)
```

## Development with Hot Reload (Air - Go's nodemon)

```bash
# Install Air globally
go install github.com/air-verse/air@latest

# Initialize Air config (creates .air.toml)
air init

# Run with hot reload
air
```

## API Testing Commands

### Health Check
```bash
curl http://localhost:8080/health
```

### Register a New User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@test.com","password":"secret123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}'
```

### Get Profile (Protected - requires token)
```bash
curl http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Get All Users (Protected)
```bash
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Build Commands

```bash
# Build for current OS
go build -o server.exe .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server .

# Build for Mac
GOOS=darwin GOARCH=amd64 go build -o server .

# Build with optimizations (smaller binary)
go build -ldflags="-s -w" -o server .
```

## Dependency Management

```bash
# Add a new dependency
go get github.com/package/name

# Update dependencies
go get -u ./...

# Clean up unused dependencies
go mod tidy

# View dependency graph
go mod graph
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./...
```
