# go get vs go install

## Quick Answer

| Command | Purpose | Example |
|---------|---------|---------|
| `go get` | Add **library** to your project | `go get github.com/gin-gonic/gin` |
| `go install` | Install **executable tool** globally | `go install github.com/air-verse/air@latest` |

## go get - Add Dependencies

**Use for:** Adding packages/libraries to import in your code

```bash
# Add a package to your project
go get github.com/gin-gonic/gin

# Add specific version
go get github.com/gin-gonic/gin@v1.9.1

# Update to latest
go get -u github.com/gin-gonic/gin
```

**What happens:**
1. Downloads the package
2. Adds it to `go.mod` (like adding to `package.json`)
3. Creates/updates `go.sum` (like `package-lock.json`)

**Then you can import it:**
```go
import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
}
```

### Node.js Equivalent
```bash
# go get = npm install --save
npm install express
```

## go install - Install Tools

**Use for:** Installing CLI tools/executables globally

```bash
# Install Air (hot reload tool)
go install github.com/air-verse/air@latest

# Install golangci-lint (linter)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install swag (Swagger docs generator)
go install github.com/swaggo/swag/cmd/swag@latest
```

**What happens:**
1. Downloads the package
2. Builds it
3. Puts the executable in `$GOPATH/bin` (usually `~/go/bin`)
4. Does NOT modify your project's `go.mod`

**Then you can run it from anywhere:**
```bash
air
golangci-lint run
swag init
```

### Node.js Equivalent
```bash
# go install = npm install -g
npm install -g nodemon
```

## Visual Comparison

```
go get github.com/gin-gonic/gin
├── Downloads package
├── Adds to go.mod        ← Project-specific
├── Adds to go.sum
└── Use: import in code

go install github.com/air-verse/air@latest
├── Downloads package
├── Builds executable
├── Puts in ~/go/bin      ← Global tool
└── Use: run as command
```

## When to Use Which?

| Scenario | Command |
|----------|---------|
| Need to `import` it in code | `go get` |
| Need to run it as CLI command | `go install` |
| Adding Gin, GORM, JWT library | `go get` |
| Installing Air, linters, code generators | `go install` |

## Common Examples

### Libraries (go get)
```bash
go get github.com/gin-gonic/gin          # Web framework
go get github.com/golang-jwt/jwt/v5      # JWT library
go get gorm.io/gorm                       # ORM
go get github.com/go-redis/redis/v8      # Redis client
```

### Tools (go install)
```bash
go install github.com/air-verse/air@latest                           # Hot reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest # Linter
go install github.com/swaggo/swag/cmd/swag@latest                    # Swagger
go install golang.org/x/tools/gopls@latest                           # Language server
```

## Important Notes

### @latest is required for go install
```bash
# This works
go install github.com/air-verse/air@latest

# This fails (no version specified)
go install github.com/air-verse/air
```

### go get inside a module
```bash
# Must be inside a Go module (folder with go.mod)
cd my_project/
go get github.com/gin-gonic/gin   # Works

# Outside a module - won't work properly
cd ~/random/
go get github.com/gin-gonic/gin   # Error or unexpected behavior
```

## Quick Reference Table

| Aspect | go get | go install |
|--------|--------|------------|
| **Purpose** | Add library to project | Install CLI tool globally |
| **Modifies go.mod** | Yes | No |
| **Where it goes** | Project's module cache | `~/go/bin` |
| **How to use result** | `import` in code | Run as command |
| **Version required** | Optional | Required (`@latest` or `@v1.2.3`) |
| **Node.js equivalent** | `npm install` | `npm install -g` |
