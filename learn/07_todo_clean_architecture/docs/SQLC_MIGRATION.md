# SQLC Migration Guide

This document explains the migration from manual SQL queries to `sqlc` for type-safe database operations.

## What is sqlc?

`sqlc` is a tool that generates type-safe Go code from SQL queries. Instead of writing manual SQL queries and scanning results, you write SQL queries in `.sql` files, and sqlc generates fully type-safe Go code.

## Benefits

- **Type Safety**: Catch SQL errors at compile time, not runtime
- **No ORM overhead**: Direct SQL queries with zero performance cost
- **Maintainability**: SQL queries in separate files, easier to review and test
- **Auto-generated boilerplate**: No more manual row scanning
- **PostgreSQL native**: Works perfectly with `database/sql` and `lib/pq`

## Project Structure

```
learn/07_todo_clean_architecture/
├── sqlc.yaml                              # sqlc configuration
├── internal/
│   └── repository/
│       ├── queries/                       # SQL query files (write your queries here)
│       │   ├── users.sql
│       │   └── todos.sql
│       ├── sqlc/                          # Generated Go code (DO NOT EDIT)
│       │   ├── db.go
│       │   ├── models.go
│       │   ├── querier.go
│       │   ├── users.sql.go
│       │   └── todos.sql.go
│       └── sqlc_impl/                     # Repository implementations using sqlc
│           ├── user_repository.go
│           └── todo_repository.go
```

## How It Works

### 1. Write SQL Queries

Create `.sql` files in `internal/repository/queries/`:

```sql
-- name: GetUserByID :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, full_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);
```

### 2. Generate Go Code

Run sqlc to generate type-safe Go code:

```bash
sqlc generate
```

This creates methods like:
- `GetUserByID(ctx, id)` - returns `(User, error)`
- `CreateUser(ctx, CreateUserParams)` - executes the insert

### 3. Use in Repository

```go
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    user, err := r.queries.GetUserByID(ctx, id)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return sqlcUserToEntity(user), nil
}
```

## Query Annotations

- `:one` - Returns a single row (`QueryRow`)
- `:many` - Returns multiple rows (`Query`)
- `:exec` - Executes query without returning rows (`Exec`)

## Nullable Parameters

For optional filters, use `sqlc.narg()`:

```sql
-- name: GetTodosFiltered :many
SELECT * FROM todos
WHERE deleted_at IS NULL
  AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
  AND (sqlc.narg('completed')::boolean IS NULL OR completed = sqlc.narg('completed'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
```

This generates:

```go
type GetTodosFilteredParams struct {
    UserID    uuid.NullUUID  `json:"user_id"`
    Completed sql.NullBool   `json:"completed"`
    Limit     int32          `json:"limit"`
    Offset    int32          `json:"offset"`
}
```

## Type Conversions

Since sqlc uses `sql.Null*` types and domain entities use regular types, we have converter functions:

```go
// toNullString converts string to sql.NullString
func toNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}

// fromNullString converts sql.NullString to string
func fromNullString(ns sql.NullString) string {
    if !ns.Valid {
        return ""
    }
    return ns.String
}
```

## Development Workflow

### Adding a New Query

1. **Write SQL query** in `internal/repository/queries/*.sql`:
   ```sql
   -- name: GetUsersByRole :many
   SELECT * FROM users WHERE role = $1 AND deleted_at IS NULL;
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
       // ... convert and return
   }
   ```

### Modifying Database Schema

1. Create migration files in `migrations/`
2. Run migrations: `make migrate-up`
3. Update SQL queries in `internal/repository/queries/`
4. Regenerate: `make sqlc-generate`

## Configuration (sqlc.yaml)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./internal/repository/queries"
    schema: "./migrations"
    gen:
      go:
        package: "sqlc"
        out: "./internal/repository/sqlc"
        sql_package: "database/sql"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
        emit_pointers_for_null_types: true
```

## Makefile Commands

```makefile
# Generate sqlc code
sqlc-generate:
	sqlc generate

# Verify sqlc queries
sqlc-verify:
	sqlc verify
```

## Migration from Old Implementation

The old `internal/repository/postgres/` implementations have been replaced with:
- `internal/repository/sqlc/` - Generated code
- `internal/repository/sqlc_impl/` - Repository implementations using sqlc

**Benefits over old approach:**
- No manual `Scan()` calls
- Compile-time type checking
- Automatic NULL handling
- Query validation during generation
- IDE autocomplete for query parameters

## Common Patterns

### Inserting with Generated UUID

```sql
-- name: CreateTodo :exec
INSERT INTO todos (id, user_id, title, ...)
VALUES ($1, $2, $3, ...);
```

```go
todo := entity.NewTodo(...)
params := sqlc.CreateTodoParams{
    ID:     todo.ID,  // generated in entity.NewTodo()
    UserID: todo.UserID,
    // ...
}
r.queries.CreateTodo(ctx, params)
```

### Filtering with Optional Parameters

```go
params := sqlc.GetTodosFilteredParams{
    Limit:  int32(limit),
    Offset: int32(offset),
}

if filter.UserID != nil {
    params.UserID = uuid.NullUUID{UUID: *filter.UserID, Valid: true}
}

todos, err := r.queries.GetTodosFiltered(ctx, params)
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

## Resources

- [sqlc Documentation](https://docs.sqlc.dev/)
- [sqlc GitHub](https://github.com/sqlc-dev/sqlc)
- [Go database/sql Tutorial](https://go.dev/doc/database/querying)
