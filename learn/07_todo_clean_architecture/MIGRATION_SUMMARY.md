# Database Layer Migration to sqlc

## Summary

Successfully migrated from manual SQL query writing to **sqlc** for type-safe database operations.

## What Changed

### Before (Manual SQL)
```go
// Manual query writing with scanning
query := `
    SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
    FROM users WHERE id = $1 AND deleted_at IS NULL
`
user := &entity.User{}
err := r.db.QueryRowContext(ctx, query, id).Scan(
    &user.ID, &user.Username, &user.Email, &user.PasswordHash,
    &user.FullName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
)
```

### After (sqlc)
```go
// Type-safe generated code
user, err := r.queries.GetUserByID(ctx, id)
```

## New Structure

```
internal/repository/
├── queries/              # SQL query files (you write these)
│   ├── users.sql
│   └── todos.sql
├── sqlc/                 # Generated Go code (DO NOT EDIT)
│   ├── db.go
│   ├── models.go
│   ├── querier.go
│   ├── users.sql.go
│   └── todos.sql.go
└── sqlc_impl/            # Repository implementations using sqlc
    ├── user_repository.go
    └── todo_repository.go
```

## Benefits

✅ **Type Safety** - Compile-time checking of SQL queries
✅ **No More Manual Scanning** - Automatically generated scanning code
✅ **Better Maintainability** - SQL in separate files, easier to review
✅ **Zero Performance Cost** - Direct SQL, no ORM overhead
✅ **IDE Support** - Full autocomplete for queries and parameters

## Migration Steps Completed

1. ✅ Installed sqlc
2. ✅ Created `sqlc.yaml` configuration
3. ✅ Created SQL query files in `internal/repository/queries/`
4. ✅ Generated type-safe Go code with `sqlc generate`
5. ✅ Implemented new repositories in `internal/repository/sqlc_impl/`
6. ✅ Updated `cmd/api/main.go` to use new repositories
7. ✅ Added Makefile targets for sqlc
8. ✅ Created documentation

## How to Use

### Writing New Queries

1. Add SQL query to `internal/repository/queries/*.sql`:
   ```sql
   -- name: GetActiveUsers :many
   SELECT * FROM users 
   WHERE deleted_at IS NULL AND active = true
   ORDER BY created_at DESC;
   ```

2. Generate Go code:
   ```bash
   make sqlc-generate
   ```

3. Use in repository:
   ```go
   users, err := r.queries.GetActiveUsers(ctx)
   ```

### Common Commands

```bash
# Generate code
make sqlc-generate

# Build application
make build

# Run migrations
make migrate-up

# Install all tools (including sqlc)
make install-tools
```

## Files Modified

- `cmd/api/main.go` - Updated to use new sqlc_impl repositories
- `Makefile` - Added sqlc-generate and sqlc-verify targets
- `go.mod` - Dependencies already compatible (using database/sql + lib/pq)

## Files Created

- `sqlc.yaml` - sqlc configuration
- `internal/repository/queries/users.sql` - User SQL queries
- `internal/repository/queries/todos.sql` - Todo SQL queries
- `internal/repository/sqlc/` - Generated code (auto-created)
- `internal/repository/sqlc_impl/user_repository.go` - New user repo implementation
- `internal/repository/sqlc_impl/todo_repository.go` - New todo repo implementation
- `docs/SQLC_MIGRATION.md` - Detailed migration guide
- `docs/SQLC_QUICK_REFERENCE.md` - Quick reference for sqlc

## Old Implementation

The previous manual implementation in `internal/repository/postgres/` is still present but no longer used. You can:
- **Keep it** as reference for comparison
- **Delete it** to clean up the codebase

To delete old implementation:
```bash
rm -rf internal/repository/postgres/
```

## Testing

The application compiles successfully and uses the new sqlc-based repositories. To test:

```bash
# Build
make build

# Run (make sure PostgreSQL is running)
make run

# Or use the binary directly
./bin/api
```

## Next Steps

1. **Test the application** with a running PostgreSQL database
2. **Write unit tests** for the new repository implementations
3. **Consider removing** old `internal/repository/postgres/` directory
4. **Add more queries** as needed using the sqlc workflow

## Documentation

- Full guide: [`docs/SQLC_MIGRATION.md`](docs/SQLC_MIGRATION.md)
- Quick reference: [`docs/SQLC_QUICK_REFERENCE.md`](docs/SQLC_QUICK_REFERENCE.md)
- sqlc official docs: https://docs.sqlc.dev/

## Dependencies

No new dependencies added! Using existing:
- `database/sql` (standard library)
- `github.com/lib/pq` (PostgreSQL driver)
- `github.com/google/uuid`

sqlc is a **code generation tool**, not a runtime dependency.
