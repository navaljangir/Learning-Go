# Codebase Context Map

> **For subagents**: Read this file FIRST before exploring. This is the single source of truth for project structure.
> **Last updated**: 2025-02-10

---

## Project: 07_todo_clean_architecture (Reference — Complete)

**Entry point**: `learn/07_todo_clean_architecture/cmd/api/main.go`
**Framework**: Gin
**DB**: PostgreSQL + SQLC (generated queries)
**Status**: Production-ready, all layers implemented

### Layer Map

| Layer | Interfaces | Implementations |
|-------|-----------|-----------------|
| **Entity** | — | `domain/entity/user.go`, `todo.go`, `todo_list.go` |
| **Repository** | `domain/repository/user_repository.go`, `todo_repository.go`, `todo_list_repository.go` | `internal/repository/sqlc_impl/user_repository.go`, `todo_repository.go`, `todo_list_repository.go` |
| **Service** | `domain/service/user_service.go`, `todo_service.go`, `todo_list_service.go` | `internal/service/user_service_impl.go`, `todo_service_impl.go`, `todo_list_service_impl.go` |
| **Handler** | `api/handler/interfaces.go` | `api/handler/auth_handler.go`, `user_handler.go`, `todo_handler.go`, `todo_list_handler.go` |
| **Router** | — | `api/router/router.go` |
| **Middleware** | — | `api/middleware/auth.go`, `logger.go`, `error_handler.go`, `cors.go` |

### Supporting Files

| Category | Files |
|----------|-------|
| **DTOs** | `internal/dto/auth_dto.go`, `user_dto.go`, `todo_dto.go`, `list_dto.go` |
| **Config** | `config/config.go` |
| **Utils** | `pkg/utils/jwt.go`, `response.go`, `error.go`, `hash.go` |
| **Constants** | `pkg/constants/constants.go` |
| **Validator** | `pkg/validator/validator.go` |
| **Migrations** | `migrations/000001_create_users_table`, `000002_create_todos_table`, `000003_add_todo_lists` |
| **SQLC** | `sqlc.yaml`, `internal/repository/queries/*.sql`, `internal/repository/sqlc/` (generated) |
| **Mocks** | `api/handler/mocks/` (auto-generated) |

### Dependency Flow

```
main.go → router.go → handlers → services (via interface) → repositories (via interface) → sqlc/DB
                     → middleware (auth, logger, cors, error)
```

---

## Project: 09_self_todo (Learning — In Progress)

**Entry point**: `learn/09_self_todo/cmd/api/main.go`
**Framework**: Gin
**DB**: Not configured yet
**Status**: Early skeleton, interfaces defined, implementations missing

### Layer Map

| Layer | Interfaces | Implementations |
|-------|-----------|-----------------|
| **Entity** | — | `domain/entity/user.go`, `todo.go` |
| **Repository** | `domain/repository/user_repository.go`, `todo_repository.go` | **NOT IMPLEMENTED** |
| **Service** | `domain/service/auth_service.go`, `todo_service.go` | **NOT IMPLEMENTED** |
| **Handler** | — | `api/handler/authHandler.go` (partial) |
| **Router** | — | Inline in `cmd/api/main.go` (setupRouter function) |
| **Middleware** | — | **NOT IMPLEMENTED** |

### Supporting Files

| Category | Files |
|----------|-------|
| **DTOs** | `internal/dto/auth_dto.go`, `todo_dto.go`, `todo_test_dto.go`, `user.go` |
| **Config** | `config/config.go` |

### What's Missing (TODO)

- [ ] Service implementations (`internal/service/`)
- [ ] Repository implementations (`internal/repository/`)
- [ ] Database setup (migrations, sqlc or GORM)
- [ ] Middleware (auth, logging, cors)
- [ ] Separate router file
- [ ] Utils (jwt, hashing, response helpers)
- [ ] Tests

---

## Conventions (Both Projects)

- **Naming**: Interfaces in `domain/`, implementations in `internal/`
- **Interface files**: Named `*_repository.go`, `*_service.go`
- **Implementation files**: Named `*_impl.go` (project 07)
- **Tests**: `*_test.go` in same package
- **DTOs**: Separate from entities, in `internal/dto/`
- **Dependency injection**: Constructor functions return interfaces, main.go wires everything
- **Error handling**: Return `error`, handlers convert to HTTP status codes
