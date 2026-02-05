# SQLC Complete Guide

A comprehensive guide to understanding and using SQLC for type-safe database operations in Go.

## Table of Contents
1. [What is SQLC?](#what-is-sqlc)
2. [Project Structure](#project-structure)
3. [Understanding SQLC](#understanding-sqlc)
4. [Configuration (sqlc.yaml)](#configuration-sqlcyaml)
5. [Query Annotations](#query-annotations)
6. [SQL to Go Type Mappings](#sql-to-go-type-mappings)
7. [Understanding sql.Null Types](#understanding-sqlnull-types)
8. [How SQLC Works - Step by Step](#how-sqlc-works---step-by-step)
9. [Common Query Patterns](#common-query-patterns)
10. [Type Converters](#type-converters)
11. [Usage in Repository](#usage-in-repository)
12. [Development Workflow](#development-workflow)
13. [Commands Reference](#commands-reference)
14. [Best Practices](#best-practices)
15. [Troubleshooting](#troubleshooting)
16. [Resources](#resources)

---

## What is SQLC?

`sqlc` is a tool that generates type-safe Go code from SQL queries. Instead of writing manual SQL queries and scanning results, you write SQL queries in `.sql` files, and sqlc generates fully type-safe Go code.

### Benefits

- **Type Safety**: Catch SQL errors at compile time, not runtime
- **No ORM overhead**: Direct SQL queries with zero performance cost
- **Maintainability**: SQL queries in separate files, easier to review and test
- **Auto-generated boilerplate**: No more manual row scanning
- **Database native**: Works perfectly with `database/sql` and database drivers
- **IDE Support**: Full autocomplete for query parameters and return types

### What SQLC Does NOT Generate

SQLC generates **Go structs and methods**, not interfaces. Here's what it creates:

| Generated Item | Description | Example |
|----------------|-------------|---------|
| **Structs (Models)** | Data models matching your database tables | `type User struct {...}` |
| **Methods (Queries)** | Functions to execute your SQL queries | `func (q *Queries) GetUserByID(...)` |
| **Queries Container** | Struct that holds all query methods | `type Queries struct { db DBTX }` |
| **Querier Interface** | Interface for mocking (optional) | `type Querier interface {...}` |

---

## Project Structure

```
learn/07_todo_clean_architecture/
├── sqlc.yaml                              # sqlc configuration
├── migrations/                            # Database schema migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000002_create_todos_table.up.sql
├── internal/
│   └── repository/
│       ├── queries/                       # SQL query files (YOU WRITE THESE)
│       │   ├── users.sql
│       │   ├── todos.sql
│       │   └── todo_lists.sql
│       ├── sqlc/                          # Generated Go code (DO NOT EDIT)
│       │   ├── db.go                      # Queries struct + New() function
│       │   ├── models.go                  # Table structs (User, Todo, etc.)
│       │   ├── querier.go                 # Querier interface (if enabled)
│       │   ├── users.sql.go               # Generated from queries/users.sql
│       │   ├── todos.sql.go               # Generated from queries/todos.sql
│       │   └── todo_lists.sql.go          # Generated from queries/todo_lists.sql
│       └── sqlc_impl/                     # Repository implementations using sqlc
│           ├── user_repository.go
│           ├── todo_repository.go
│           └── todo_list_repository.go
```

**Key Directories:**
- `migrations/` - SQLC reads these to understand your database schema
- `queries/` - You write SQL queries here
- `sqlc/` - SQLC generates Go code here (never edit manually)
- `sqlc_impl/` - Your repository implementations that use SQLC

---

## Understanding SQLC

### 1. Structs - Data Models

SQLC reads your database schema and creates Go structs to match:

```go
// Generated from your 'users' table
type User struct {
    ID        string
    Username  string
    Email     sql.NullString  // Nullable column
    CreatedAt time.Time
}
```

**What we call them:** These are called **"models"** or **"generated structs"**

### 2. Methods - Query Functions

SQLC reads your `.sql` files and creates Go functions:

```go
// Generated from: -- name: GetUserByID :one
func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
    // ... generated query code
}
```

**What we call them:** These are called **"generated queries"** or **"query methods"**

### 3. Queries Struct - The Container

```go
type Queries struct {
    db DBTX  // Can be *sql.DB or *sql.Tx
}

func New(db DBTX) *Queries {
    return &Queries{db: db}
}
```

**What this is:** The `Queries` struct holds all your generated query methods

---

## Configuration (sqlc.yaml)

Your project's SQLC configuration:

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "./internal/repository/queries"
    schema: "./migrations"
    gen:
      go:
        package: "sqlc"
        out: "./internal/repository/sqlc"
        sql_package: "database/sql"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_pointers_for_null_types: true
```

### Key Configuration Options

| Key | Value in Your Config | What It Does |
|-----|---------------------|--------------|
| `version` | `"2"` | SQLC config file version |
| `engine` | `"mysql"` | Database type (mysql, postgresql, sqlite) |
| `queries` | `"./internal/repository/queries"` | **Where SQLC reads your .sql files** |
| `schema` | `"./migrations"` | **Where SQLC reads your migrations to understand table structure** |
| `package` | `"sqlc"` | **Name of the generated Go package** |
| `out` | `"./internal/repository/sqlc"` | **Where SQLC writes generated Go files** |

### Generation Options (The `gen.go` section)

| Option | Your Value | What It Does | Example |
|--------|-----------|--------------|---------|
| `sql_package` | `"database/sql"` | Use standard library `database/sql` types | `sql.NullString`, `sql.NullTime` |
| `emit_json_tags` | `true` | Add JSON tags to structs | `json:"id"` on struct fields |
| `emit_prepared_queries` | `false` | Don't generate prepared statement versions | Simpler code, slightly slower |
| `emit_interface` | `true` | **Generate Querier interface** | Allows mocking for tests |
| `emit_exact_table_names` | `false` | Struct names don't exactly match table names | `User` instead of `users` |
| `emit_empty_slices` | `true` | Return `[]Type{}` instead of `nil` for empty results | Easier JSON marshaling |
| `emit_pointers_for_null_types` | `true` | Use `*string` for nullable columns | Clear distinction: `nil` vs `""` |

---

## Query Annotations

Every SQL query must have an annotation comment:

```sql
-- name: MethodName :type
```

### Annotation Types

| Type | Returns | Use For | Example |
|------|---------|---------|---------|
| `:one` | Single row | `SELECT ... WHERE id = ?` | `GetUserByID(ctx, id) (User, error)` |
| `:many` | Multiple rows | `SELECT ... WHERE user_id = ?` | `ListUsers(ctx) ([]User, error)` |
| `:exec` | No return value | `INSERT`, `UPDATE`, `DELETE` | `CreateUser(ctx, params) error` |
| `:execresult` | `sql.Result` | Get rows affected / last insert ID | `DeleteUser(ctx, id) (sql.Result, error)` |

### Examples

```sql
-- :one - Returns a single row
-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? AND deleted_at IS NULL;

-- :many - Returns multiple rows
-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC;

-- :exec - Executes without return
-- name: CreateUser :exec
INSERT INTO users (id, username, email) VALUES (?, ?, ?);

-- :execresult - Returns sql.Result
-- name: DeleteUser :execresult
DELETE FROM users WHERE id = ?;
```

---

## SQL to Go Type Mappings

### MySQL Types

| MySQL Type | Go Type (Non-Nullable) | Go Type (Nullable) |
|------------|------------------------|-------------------|
| `CHAR(36)` (UUID) | `string` | `sql.NullString` |
| `VARCHAR` | `string` | `sql.NullString` |
| `TEXT` | `string` | `sql.NullString` |
| `INT` | `int32` | `sql.NullInt32` |
| `BIGINT` | `int64` | `sql.NullInt64` |
| `BOOLEAN` | `bool` | `sql.NullBool` |
| `TIMESTAMP` | `time.Time` | `sql.NullTime` |
| `DATETIME` | `time.Time` | `sql.NullTime` |
| `FLOAT` | `float32` | - |
| `DOUBLE` | `float64` | `sql.NullFloat64` |

### PostgreSQL Types

| PostgreSQL Type | Go Type (Non-Nullable) | Go Type (Nullable) |
|----------------|------------------------|-------------------|
| `UUID` | `uuid.UUID` | `uuid.NullUUID` |
| `VARCHAR` | `string` | `sql.NullString` |
| `TEXT` | `string` | `sql.NullString` |
| `INTEGER` | `int32` | `sql.NullInt32` |
| `BIGINT` | `int64` | `sql.NullInt64` |
| `BOOLEAN` | `bool` | `sql.NullBool` |
| `TIMESTAMP` | `time.Time` | `sql.NullTime` |

---

## Understanding sql.Null Types

### The Problem: Go Can't Represent Database NULL

In SQL:
- `NULL` = "no value"
- `""` (empty string) = "value that is empty"

In Go:
- Empty string `""` is a valid value
- How do we represent "no value"?

**Solution:** `sql.NullString`, `sql.NullTime`, etc.

### sql.NullString Structure

```go
type NullString struct {
    String string  // The actual value
    Valid  bool    // Is this NULL or not?
}
```

### How It Works

| Database Value | Go Representation |
|----------------|-------------------|
| `NULL` | `sql.NullString{String: "", Valid: false}` |
| `""` (empty string) | `sql.NullString{String: "", Valid: true}` |
| `"hello"` | `sql.NullString{String: "hello", Valid: true}` |

### Usage Pattern

```go
// Reading from database
var email sql.NullString

if email.Valid {
    // Database had a value (even if empty string)
    fmt.Println("Email:", email.String)
} else {
    // Database had NULL
    fmt.Println("No email")
}

// Writing a value to database
email := sql.NullString{
    String: "user@example.com",
    Valid:  true,  // Set Valid to true when you have a value
}

// Writing NULL to database
email := sql.NullString{
    String: "",     // Doesn't matter
    Valid:  false,  // Valid = false means NULL
}
```

### Common Null Types

| Type | Go Type | Used For |
|------|---------|----------|
| `sql.NullString` | `string` | VARCHAR, TEXT columns |
| `sql.NullInt32` | `int32` | INT columns |
| `sql.NullInt64` | `int64` | BIGINT columns |
| `sql.NullFloat64` | `float64` | FLOAT, DOUBLE columns |
| `sql.NullBool` | `bool` | BOOLEAN columns |
| `sql.NullTime` | `time.Time` | TIMESTAMP, DATETIME columns |

### Quick Check Pattern

```go
// Has value?
if nullString.Valid {
    value := nullString.String
}

// Create with value
ns := sql.NullString{String: "hello", Valid: true}

// Create NULL
ns := sql.NullString{Valid: false}
```

---

## How SQLC Works - Step by Step

### Step 1: You Write SQL

```sql
-- internal/repository/queries/users.sql

-- name: GetUserByID :one
SELECT id, username, email, created_at
FROM users
WHERE id = ?;
```

### Step 2: You Run `sqlc generate`

```bash
sqlc generate
```

### Step 3: SQLC Reads Your Migrations

It looks in `./migrations` to understand your database schema:

```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NULL,  -- ⭐ NULL = optional column
    created_at TIMESTAMP NOT NULL
);
```

### Step 4: SQLC Generates Go Code

**File:** `internal/repository/sqlc/users.sql.go`

```go
type GetUserByIDRow struct {
    ID        string
    Username  string
    Email     sql.NullString  // ⭐ NULL column → sql.NullString
    CreatedAt time.Time
}

func (q *Queries) GetUserByID(ctx context.Context, id string) (GetUserByIDRow, error) {
    row := q.db.QueryRowContext(ctx, getUserByID, id)
    var i GetUserByIDRow
    err := row.Scan(&i.ID, &i.Username, &i.Email, &i.CreatedAt)
    return i, err
}
```

### Step 5: You Use It In Your Code

```go
queries := sqlc.New(db)
user, err := queries.GetUserByID(ctx, "some-uuid")
if err != nil {
    return err
}

// Check if email is NULL
if user.Email.Valid {
    fmt.Println("Email:", user.Email.String)
} else {
    fmt.Println("No email set")
}
```

---

## Common Query Patterns

### SELECT One

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? AND deleted_at IS NULL;
```

**Generated:**
```go
func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error)
```

### SELECT Many

```sql
-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
```

**Generated:**
```go
type ListUsersParams struct {
    Limit  int32
    Offset int32
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
```

### INSERT

```sql
-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);
```

**Generated:**
```go
type CreateUserParams struct {
    ID           string
    Username     string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error
```

### UPDATE

```sql
-- name: UpdateUser :exec
UPDATE users
SET full_name = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL;
```

**Generated:**
```go
type UpdateUserParams struct {
    FullName  string
    UpdatedAt time.Time
    ID        string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error
```

### DELETE (Soft Delete)

```sql
-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;
```

**Generated:**
```go
func (q *Queries) SoftDeleteUser(ctx context.Context, id string) error
```

### COUNT

```sql
-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;
```

**Generated:**
```go
func (q *Queries) CountUsers(ctx context.Context) (int64, error)
```

### Filtering with Optional Parameters (MySQL with sqlc.arg)

```sql
-- name: SearchTodos :many
SELECT * FROM todos
WHERE deleted_at IS NULL
  AND (? IS NULL OR user_id = ?)
  AND (? IS NULL OR completed = ?)
ORDER BY created_at DESC;
```

**Generated:**
```go
type SearchTodosParams struct {
    UserID      sql.NullString
    UserID_2    string
    Completed   sql.NullBool
    Completed_2 bool
}
```

### Inserting with Generated UUID

```sql
-- name: CreateTodo :exec
INSERT INTO todos (id, user_id, title, description, ...)
VALUES (?, ?, ?, ?, ...);
```

**Usage:**
```go
todo := entity.NewTodo(...)  // Generates UUID internally
params := sqlc.CreateTodoParams{
    ID:     todo.ID.String(),  // Use the generated UUID
    UserID: todo.UserID.String(),
    Title:  todo.Title,
    // ...
}
r.queries.CreateTodo(ctx, params)
```

---

## Type Converters

Since SQLC uses `sql.Null*` types and domain entities use regular types, we need converter functions:

### String Converters

```go
// Non-nullable string to sql.NullString
func toNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}

// sql.NullString to non-nullable string
func fromNullString(ns sql.NullString) string {
    if !ns.Valid {
        return ""
    }
    return ns.String
}
```

### Time Converters

```go
// time.Time to sql.NullTime
func toNullTime(t time.Time) sql.NullTime {
    if t.IsZero() {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{Time: t, Valid: true}
}

// *time.Time to sql.NullTime
func toNullTimePtr(t *time.Time) sql.NullTime {
    if t == nil {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{Time: *t, Valid: true}
}

// sql.NullTime to *time.Time
func fromNullTimePtr(nt sql.NullTime) *time.Time {
    if !nt.Valid {
        return nil
    }
    return &nt.Time
}
```

### UUID Converters

```go
// *uuid.UUID to sql.NullString (for MySQL)
func uuidToNullString(id *uuid.UUID) sql.NullString {
    if id == nil {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: id.String(), Valid: true}
}

// sql.NullString to *uuid.UUID (for MySQL)
func nullStringToUUID(ns sql.NullString) *uuid.UUID {
    if !ns.Valid {
        return nil
    }
    id, err := uuid.Parse(ns.String)
    if err != nil {
        return nil
    }
    return &id
}
```

### Converting Domain Entities

```go
// Pattern 1: Converting Go Pointer to sql.Null

// Entity has *uuid.UUID (pointer = nullable)
type Todo struct {
    ListID *uuid.UUID  // nil = no list
}

// SQLC expects sql.NullString
func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
    var listID sql.NullString

    if todo.ListID != nil {
        // Entity has a value → Valid = true
        listID = sql.NullString{
            String: todo.ListID.String(),
            Valid:  true,
        }
    } else {
        // Entity has nil → Valid = false (NULL in DB)
        listID = sql.NullString{Valid: false}
    }

    params := sqlc.CreateTodoParams{
        ListID: listID,
        // ...
    }
    return r.queries.CreateTodo(ctx, params)
}
```

```go
// Pattern 2: Converting sql.Null to Go Pointer

// SQLC returns sql.NullString
type GetTodoByIDRow struct {
    ListID sql.NullString
}

// Convert to Entity with *uuid.UUID
func sqlcTodoToEntity(t GetTodoByIDRow) *entity.Todo {
    var listID *uuid.UUID

    if t.ListID.Valid {
        // Database had a value → convert to pointer
        parsed, _ := uuid.Parse(t.ListID.String)
        listID = &parsed
    }
    // If not Valid, listID stays nil

    return &entity.Todo{
        ListID: listID,  // nil or *uuid.UUID
        // ...
    }
}
```

---

## Usage in Repository

### Basic Query

```go
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    user, err := r.queries.GetUserByID(ctx, id.String())
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return sqlcUserToEntity(user), nil
}
```

### Query with Parameters

```go
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
    params := sqlc.ListUsersParams{
        Limit:  int32(limit),
        Offset: int32(offset),
    }

    users, err := r.queries.ListUsers(ctx, params)
    if err != nil {
        return nil, err
    }

    result := make([]*entity.User, len(users))
    for i, u := range users {
        result[i] = sqlcUserToEntity(u)
    }

    return result, nil
}
```

### Insert

```go
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    params := sqlc.CreateUserParams{
        ID:           user.ID.String(),
        Username:     user.Username,
        Email:        user.Email,
        PasswordHash: user.PasswordHash,
        FullName:     toNullString(user.FullName),
        CreatedAt:    user.CreatedAt,
        UpdatedAt:    user.UpdatedAt,
    }

    return r.queries.CreateUser(ctx, params)
}
```

### Optional Filters

```go
func (r *todoRepository) FindWithFilters(ctx context.Context, filter repository.TodoFilter, limit, offset int) ([]*entity.Todo, error) {
    params := sqlc.GetTodosFilteredParams{
        Limit:  int32(limit),
        Offset: int32(offset),
    }

    if filter.UserID != nil {
        params.UserID = sql.NullString{
            String: filter.UserID.String(),
            Valid:  true,
        }
    }

    if filter.Completed != nil {
        params.Completed = sql.NullBool{
            Bool:  *filter.Completed,
            Valid: true,
        }
    }

    todos, err := r.queries.GetTodosFiltered(ctx, params)
    if err != nil {
        return nil, err
    }

    result := make([]*entity.Todo, len(todos))
    for i, t := range todos {
        result[i] = sqlcTodoToEntity(t)
    }

    return result, nil
}
```

### Handling Not Found Errors

```go
user, err := r.queries.GetUserByID(ctx, id)
if err == sql.ErrNoRows {
    return nil, errors.New("user not found")
}
if err != nil {
    return nil, err
}
```

---

## Development Workflow

### Adding a New Query

1. **Write SQL query** in `internal/repository/queries/*.sql`:
   ```sql
   -- name: GetUsersByRole :many
   SELECT * FROM users WHERE role = ? AND deleted_at IS NULL;
   ```

2. **Generate Go code**:
   ```bash
   make sqlc-generate
   # or
   sqlc generate
   ```

3. **Use in repository**:
   ```go
   func (r *userRepository) FindByRole(ctx context.Context, role string) ([]*entity.User, error) {
       users, err := r.queries.GetUsersByRole(ctx, role)
       if err != nil {
           return nil, err
       }

       result := make([]*entity.User, len(users))
       for i, u := range users {
           result[i] = sqlcUserToEntity(u)
       }

       return result, nil
   }
   ```

### Modifying Database Schema

1. **Create migration files** in `migrations/`:
   ```bash
   migrate create -ext sql -dir migrations -seq add_role_to_users
   ```

2. **Write up migration**:
   ```sql
   -- migrations/000003_add_role_to_users.up.sql
   ALTER TABLE users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';
   ```

3. **Write down migration**:
   ```sql
   -- migrations/000003_add_role_to_users.down.sql
   ALTER TABLE users DROP COLUMN role;
   ```

4. **Run migrations**:
   ```bash
   make migrate-up
   ```

5. **Update SQL queries** in `internal/repository/queries/`

6. **Regenerate SQLC code**:
   ```bash
   make sqlc-generate
   ```

### Complete Flow: New Feature

```bash
# 1. Create migrations
migrate create -ext sql -dir migrations -seq create_todo_lists

# 2. Write migration SQL (up and down)
# Edit migrations/000003_create_todo_lists.up.sql
# Edit migrations/000003_create_todo_lists.down.sql

# 3. Run migrations
make migrate-up

# 4. Write queries
# Edit internal/repository/queries/todo_lists.sql

# 5. Generate SQLC code
make sqlc-generate

# 6. Implement repository
# Create internal/repository/sqlc_impl/todo_list_repository.go

# 7. Wire up in main.go
# Create domain/repository/todo_list_repository.go interface
# Register in dependency injection

# 8. Test
go test ./...
```

---

## Commands Reference

### Common Commands

```bash
# Generate Go code from SQL queries
make sqlc-generate

# Or directly:
sqlc generate

# Verify queries without generating
make sqlc-verify

# Run migrations
make migrate-up
make migrate-down
```

### Makefile Commands

```makefile
# Generate sqlc code
sqlc-generate:
	sqlc generate

# Verify sqlc queries
sqlc-verify:
	sqlc verify

# Run migrations
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

# Install tools
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

## Best Practices

### 1. Write SQL in `.sql` Files
Don't write SQL in Go code. Keep SQL queries in separate `.sql` files for better maintainability.

**Good:**
```sql
-- internal/repository/queries/users.sql
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? AND deleted_at IS NULL;
```

**Bad:**
```go
// Don't do this
query := "SELECT * FROM users WHERE email = ? AND deleted_at IS NULL"
```

### 2. Use Meaningful Query Names
Use descriptive names that clearly indicate what the query does.

**Good:**
```sql
-- name: GetUserByEmail :one
-- name: ListActiveUsers :many
-- name: CreateUser :exec
```

**Bad:**
```sql
-- name: Query1 :one
-- name: GetUser :one  (ambiguous - by ID? by email?)
-- name: DoStuff :exec
```

### 3. Handle NULL Properly
Use nullable types for optional fields and convert appropriately at boundaries.

```go
// Domain entity with pointer (nullable)
type Todo struct {
    ListID *uuid.UUID
}

// Convert to SQLC type
listID := sql.NullString{Valid: false}
if todo.ListID != nil {
    listID = sql.NullString{String: todo.ListID.String(), Valid: true}
}
```

### 4. Convert at Boundaries
Keep SQLC types in repository layer, convert to domain entities at boundaries.

```go
// Repository layer - uses SQLC types
func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
    params := sqlc.CreateTodoParams{
        // Convert entity → SQLC params
    }
    return r.queries.CreateTodo(ctx, params)
}

// Service layer - uses domain entities
func (s *todoService) CreateTodo(ctx context.Context, todo *entity.Todo) error {
    return s.repo.Create(ctx, todo)
}
```

### 5. Keep Queries Simple
Put complex logic in Go, keep SQL queries straightforward.

**Good:**
```sql
-- name: GetUserTodos :many
SELECT * FROM todos WHERE user_id = ? AND deleted_at IS NULL;
```

```go
// Complex filtering in Go
todos, err := r.queries.GetUserTodos(ctx, userID)
if err != nil {
    return nil, err
}

// Filter in Go
result := filterByPriority(todos, filter.Priority)
result = sortByDueDate(result)
return result, nil
```

**Bad:**
```sql
-- Too complex - hard to maintain
SELECT t.*,
       CASE WHEN priority = 'high' THEN 1 WHEN priority = 'medium' THEN 2 ELSE 3 END as priority_order
FROM todos t
WHERE user_id = ?
  AND deleted_at IS NULL
  AND (? IS NULL OR priority = ?)
ORDER BY priority_order, due_date;
```

### 6. Test Your Repositories
Write tests for repository methods.

```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := sqlc_impl.NewUserRepository(db)
    user := &entity.User{
        ID:       uuid.New(),
        Username: "testuser",
        Email:    "test@example.com",
    }

    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)

    // Verify
    found, err := repo.FindByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Username, found.Username)
}
```

### 7. Regenerate After Schema Changes
Always run `sqlc generate` after modifying migrations or queries.

```bash
# After editing migrations or queries
sqlc generate

# Verify everything compiles
go build ./...

# Run tests
go test ./...
```

### 8. Version Control
- **Commit:** `sqlc.yaml`, `migrations/`, `queries/`
- **Commit:** Generated code in `sqlc/` (helps with code review)
- **Don't edit:** Never manually edit generated files

### 9. Use Soft Deletes
Prefer soft deletes (setting `deleted_at`) over hard deletes.

```sql
-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL;

-- Always filter out soft-deleted records
-- name: GetActiveUsers :many
SELECT * FROM users WHERE deleted_at IS NULL;
```

### 10. Handle Errors Consistently
Always check for `sql.ErrNoRows` separately from other errors.

```go
user, err := r.queries.GetUserByID(ctx, id)
if err == sql.ErrNoRows {
    return nil, ErrUserNotFound  // Custom domain error
}
if err != nil {
    return nil, fmt.Errorf("database error: %w", err)
}
```

---

## Troubleshooting

### "sqlc: command not found"

**Solution:**
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Or add to Makefile
make install-tools
```

Ensure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### "relation does not exist" Error

**Problem:**
```
internal/repository/queries/todo_lists.sql:1:1: relation "todo_lists" does not exist
```

**Why:**
SQLC reads migrations to understand database structure, but your database doesn't have the table yet.

**Solution (in order):**
```bash
# Step 1: Run migrations (creates tables)
make migrate-up

# Step 2: Now SQLC can validate queries
sqlc generate

# Step 3: Build and run
go build ./cmd/api
```

### Generated Code Has Wrong Types

**Problem:** SQLC generates `interface{}` instead of specific types.

**Solution:** Check your `sqlc.yaml` configuration, especially the `sql_package` setting:

```yaml
gen:
  go:
    sql_package: "database/sql"  # Use standard library types
```

### Query Not Found After Generation

**Problem:** Query method doesn't exist after running `sqlc generate`.

**Solution:** Check your query annotation format:

```sql
-- Correct format (note the space after --)
-- name: QueryName :one

-- Wrong (missing space)
--name: QueryName :one

-- Wrong (missing colon)
-- name: QueryName one
```

### NULL Handling Errors

**Problem:** Runtime errors with NULL values.

**Solution:** Use appropriate nullable types:

```go
// For nullable VARCHAR
listID sql.NullString

// For nullable BOOLEAN
completed sql.NullBool

// For nullable TIMESTAMP
deletedAt sql.NullTime

// For nullable UUID (MySQL with CHAR(36))
userID sql.NullString

// For nullable UUID (PostgreSQL)
userID uuid.NullUUID
```

### Type Mismatch Errors

**Problem:**
```
cannot use params (type CreateUserParams) as type in argument
```

**Solution:** Make sure you're passing the correct params struct:

```go
// SQLC generates a params struct for multi-parameter queries
params := sqlc.CreateUserParams{
    ID:       id.String(),
    Username: username,
    // ... all required fields
}

r.queries.CreateUser(ctx, params)
```

### Syntax Errors in Generated Code

**Problem:** Generated code doesn't compile.

**Solution:**
1. Check your SQL syntax in `.sql` files
2. Run `sqlc verify` to check for issues
3. Delete `sqlc/` folder and regenerate:
   ```bash
   rm -rf internal/repository/sqlc/*
   sqlc generate
   ```

### Database Connection Issues

**Problem:** Can't connect to database for query validation.

**Solution:** SQLC can work without database connection. Use `emit_db_tags: false` in config:

```yaml
gen:
  go:
    emit_db_tags: false  # Don't require DB connection for generation
```

### Parameter Placeholder Issues

**MySQL uses `?`:**
```sql
SELECT * FROM users WHERE id = ? AND email = ?;
```

**PostgreSQL uses `$1, $2, ...`:**
```sql
SELECT * FROM users WHERE id = $1 AND email = $2;
```

Make sure your placeholders match your database engine setting in `sqlc.yaml`.

---

## Resources

### Official Documentation
- [SQLC Documentation](https://docs.sqlc.dev/)
- [SQLC GitHub](https://github.com/sqlc-dev/sqlc)
- [Go database/sql Tutorial](https://go.dev/doc/database/querying)

### Related Tools
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [dbmate](https://github.com/amacneil/dbmate) - Alternative migration tool
- [Atlas](https://atlasgo.io/) - Modern schema management

### Learning Resources
- [SQLC Playground](https://play.sqlc.dev/) - Try SQLC online
- [Go Database Best Practices](https://www.alexedwards.net/blog/practical-persistence-sql)
- [database/sql Tutorial](https://go.dev/doc/tutorial/database-access)

### Example Projects
- [SQLC Examples](https://github.com/sqlc-dev/sqlc/tree/main/examples)
- [This Project](../) - Real-world clean architecture with SQLC

---

## Real Example from Your Code

### Your SQL Query

```sql
-- name: CreateTodo :exec
INSERT INTO todos (id, user_id, list_id, title, description, ...)
VALUES (?, ?, ?, ?, ?, ...);
```

### Your Table (from migration)

```sql
CREATE TABLE todos (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    list_id CHAR(36) NULL,          -- ⭐ NULLABLE column
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,          -- ⭐ NULLABLE column
    ...
);
```

### SQLC Generated Struct

```go
type CreateTodoParams struct {
    ID          string
    UserID      string
    ListID      sql.NullString  // ⭐ Because list_id is NULL in SQL
    Title       string
    Description sql.NullString  // ⭐ Because description is NULL in SQL
}
```

### How You Use It

```go
// Todo WITHOUT a list (global todo)
params := sqlc.CreateTodoParams{
    ID:     uuid.New().String(),
    UserID: userID.String(),
    ListID: sql.NullString{Valid: false},  // ⭐ NULL in database
    Title:  "Buy groceries",
}

// Todo WITH a list
params := sqlc.CreateTodoParams{
    ID:     uuid.New().String(),
    UserID: userID.String(),
    ListID: sql.NullString{
        String: "some-list-uuid",
        Valid:  true,  // ⭐ NOT NULL, has a value
    },
    Title:  "Finish project",
}

queries.CreateTodo(ctx, params)
```

---

## Summary

**What SQLC Generates:** Structs (models) + Methods (queries), not interfaces

**sqlc.yaml:** Tells SQLC where to read SQL files, where migrations are, and where to output Go code

**sql.Null Types:** Go's way of representing database NULL values using a struct with `{Value, Valid bool}` pattern

**Why Valid Field:** Distinguishes between NULL (`Valid: false`) and actual empty value (`Valid: true, String: ""`)

**Development Flow:** Write SQL → Run migrations → Generate code → Implement repository → Test

**Key Benefit:** Type-safe database operations with zero runtime overhead
