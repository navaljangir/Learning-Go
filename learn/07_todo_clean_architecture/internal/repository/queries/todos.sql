-- name: CreateTodo :exec
INSERT INTO todos (id, user_id, title, description, completed, priority, due_date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetTodoByID :one
SELECT id, user_id, title, description, completed, priority, due_date, created_at, updated_at, completed_at, deleted_at
FROM todos
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTodosByUserID :many
SELECT id, user_id, title, description, completed, priority, due_date, created_at, updated_at, completed_at, deleted_at
FROM todos
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTodo :exec
UPDATE todos
SET title = $2, description = $3, completed = $4, priority = $5, due_date = $6, updated_at = $7, completed_at = $8
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteTodo :exec
UPDATE todos
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CountTodosByUser :one
SELECT COUNT(*)
FROM todos
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CountTodos :one
SELECT COUNT(*)
FROM todos
WHERE deleted_at IS NULL;

-- name: GetTodosFiltered :many
SELECT id, user_id, title, description, completed, priority, due_date, created_at, updated_at, completed_at, deleted_at
FROM todos
WHERE deleted_at IS NULL
  AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
  AND (sqlc.narg('completed')::boolean IS NULL OR completed = sqlc.narg('completed'))
  AND (sqlc.narg('priority')::varchar IS NULL OR priority = sqlc.narg('priority'))
  AND (sqlc.narg('from_date')::timestamp IS NULL OR due_date >= sqlc.narg('from_date'))
  AND (sqlc.narg('to_date')::timestamp IS NULL OR due_date <= sqlc.narg('to_date'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountTodosFiltered :one
SELECT COUNT(*)
FROM todos
WHERE deleted_at IS NULL
  AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
  AND (sqlc.narg('completed')::boolean IS NULL OR completed = sqlc.narg('completed'))
  AND (sqlc.narg('priority')::varchar IS NULL OR priority = sqlc.narg('priority'))
  AND (sqlc.narg('from_date')::timestamp IS NULL OR due_date >= sqlc.narg('from_date'))
  AND (sqlc.narg('to_date')::timestamp IS NULL OR due_date <= sqlc.narg('to_date'));
