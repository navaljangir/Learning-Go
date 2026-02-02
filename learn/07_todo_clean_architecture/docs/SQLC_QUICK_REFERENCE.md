# sqlc Quick Reference

Quick reference guide for working with sqlc in this project.

## Common Commands

```bash
# Generate Go code from SQL queries
make sqlc-generate

# Or directly:
sqlc generate

# Verify queries without generating
make sqlc-verify
```

## Query Annotations

| Annotation | Description | Returns |
|-----------|-------------|---------|
| `:one` | Single row query | `(Type, error)` |
| `:many` | Multiple rows query | `([]Type, error)` |
| `:exec` | Execute without return | `error` |
| `:execrows` | Execute and return rows affected | `(int64, error)` |

## Parameter Types

### Named Parameters
```sql
-- Use sqlc.arg() for required parameters
WHERE user_id = sqlc.arg('user_id')
LIMIT sqlc.arg('limit')
```

### Nullable Parameters
```sql
-- Use sqlc.narg() for optional parameters
WHERE (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
```

## SQL Types to Go Types

| PostgreSQL | Go Type (sql_package: "database/sql") |
|------------|----------------------------------------|
| `UUID` | `uuid.UUID` |
| `UUID` (nullable) | `uuid.NullUUID` |
| `VARCHAR` | `string` |
| `VARCHAR` (nullable) | `sql.NullString` |
| `INTEGER` | `int32` |
| `BIGINT` | `int64` |
| `BOOLEAN` | `bool` |
| `BOOLEAN` (nullable) | `sql.NullBool` |
| `TIMESTAMP` | `time.Time` |
| `TIMESTAMP` (nullable) | `sql.NullTime` |

## Common Patterns

### SELECT One
```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;
```

Generated:
```go
func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
```

### SELECT Many
```sql
-- name: ListUsers :many
SELECT * FROM users 
WHERE deleted_at IS NULL 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
```

Generated:
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
VALUES ($1, $2, $3, $4, $5, $6);
```

Generated:
```go
type CreateUserParams struct {
    ID           uuid.UUID
    Username     string
    Email        string
    PasswordHash string
    CreatedAt    sql.NullTime
    UpdatedAt    sql.NullTime
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error
```

### UPDATE
```sql
-- name: UpdateUser :exec
UPDATE users
SET full_name = $2, updated_at = $3
WHERE id = $1 AND deleted_at IS NULL;
```

### DELETE (Soft Delete)
```sql
-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
```

### COUNT
```sql
-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;
```

Generated:
```go
func (q *Queries) CountUsers(ctx context.Context) (int64, error)
```

### Filtering with Optional Parameters
```sql
-- name: SearchUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
  AND (sqlc.narg('username')::varchar IS NULL OR username ILIKE sqlc.narg('username'))
  AND (sqlc.narg('email')::varchar IS NULL OR email ILIKE sqlc.narg('email'))
ORDER BY created_at DESC;
```

Generated:
```go
type SearchUsersParams struct {
    Username sql.NullString
    Email    sql.NullString
}

func (q *Queries) SearchUsers(ctx context.Context, arg SearchUsersParams) ([]User, error)
```

## Usage in Repository

### Basic Query
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
        ID:           user.ID,
        Username:     user.Username,
        Email:        user.Email,
        PasswordHash: user.PasswordHash,
        FullName:     toNullString(user.FullName),
        CreatedAt:    toNullTime(user.CreatedAt),
        UpdatedAt:    toNullTime(user.UpdatedAt),
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
        params.UserID = uuid.NullUUID{UUID: *filter.UserID, Valid: true}
    }
    
    if filter.Completed != nil {
        params.Completed = sql.NullBool{Bool: *filter.Completed, Valid: true}
    }
    
    todos, err := r.queries.GetTodosFiltered(ctx, params)
    if err != nil {
        return nil, err
    }
    
    // Convert to domain entities...
}
```

## Type Converters

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
// *uuid.UUID to uuid.NullUUID
if userID != nil {
    params.UserID = uuid.NullUUID{UUID: *userID, Valid: true}
}
```

## Best Practices

1. **Write SQL in `.sql` files** - Don't write SQL in Go code
2. **Use meaningful query names** - `GetUserByEmail`, not `Query1`
3. **Handle NULL properly** - Use nullable types for optional fields
4. **Convert at boundaries** - Convert between sqlc types and domain entities
5. **Keep queries simple** - Complex logic in Go, simple queries in SQL
6. **Test queries** - Write tests for your repository methods
7. **Regenerate after schema changes** - Always run `make sqlc-generate`

## Troubleshooting

### "sqlc: command not found"
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# or
make install-tools
```

### Generated code has wrong types
Check your sqlc.yaml configuration, especially `sql_package` setting.

### Query not found
Make sure your query has the proper annotation format:
```sql
-- name: QueryName :one
```

### NULL handling errors
Use appropriate nullable types:
- `sql.NullString` for nullable VARCHAR
- `sql.NullBool` for nullable BOOLEAN
- `sql.NullTime` for nullable TIMESTAMP
- `uuid.NullUUID` for nullable UUID
