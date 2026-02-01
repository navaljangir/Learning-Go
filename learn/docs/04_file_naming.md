# Go File Naming (main.go isn't special!)

## Key Concept

In Go, **file names don't matter**. What matters is:
1. `package main` - marks this as an executable
2. `func main()` - the entry point

## Valid Examples

All of these work exactly the same:

```
# Option 1: Traditional
my_app/
└── main.go          ← package main + func main()

# Option 2: Custom name
my_app/
└── app.go           ← package main + func main()

# Option 3: Multiple files
my_app/
├── app.go           ← package main + func main()
├── helpers.go       ← package main (no func main - only one allowed)
└── config.go        ← package main
```

## Running Your Code

```bash
# Run the entire package (recommended)
go run .

# Run specific file(s)
go run app.go
go run app.go helpers.go config.go

# Build the entire package
go build .
go build -o myapp.exe .
```

## Node.js Comparison

| Node.js | Go |
|---------|-----|
| `package.json` → `"main": "app.js"` | Any file with `func main()` |
| File name matters | File name doesn't matter |
| `node app.js` | `go run .` |

## The Rule: One main() Per Package

```go
// app.go
package main

func main() {           // ✅ Only ONE main function
    startServer()
}

// helpers.go
package main

func startServer() {    // ✅ Helper functions are fine
    // ...
}

func main() {           // ❌ ERROR: main redeclared
}
```

## Multiple Entry Points (Multiple Mains)

If you need multiple executables, use `cmd/` folder:

```
my_project/
├── cmd/
│   ├── api/
│   │   └── main.go      ← go run ./cmd/api
│   ├── worker/
│   │   └── main.go      ← go run ./cmd/worker
│   └── cli/
│       └── main.go      ← go run ./cmd/cli
├── internal/
│   └── shared.go        ← shared code
└── go.mod
```

```bash
# Run different entry points
go run ./cmd/api
go run ./cmd/worker
go run ./cmd/cli

# Build all
go build -o bin/api ./cmd/api
go build -o bin/worker ./cmd/worker
```

## Air with Multiple Entry Points

```toml
# .air.toml for running ./cmd/api
[build]
cmd = "go build -o ./tmp/main.exe ./cmd/api"
bin = "./tmp/main.exe"
```

## Summary

| Question | Answer |
|----------|--------|
| Must file be named `main.go`? | No, any name works |
| What makes it executable? | `package main` + `func main()` |
| Can I have multiple `.go` files? | Yes, all with `package main` |
| Can I have multiple `main()` functions? | No, only one per package |
| How to have multiple executables? | Use `cmd/` subfolders |
