-- name: CreateTodo :exec
INSERT INTO todos (id, user_id, list_id, title, description, completed, priority, due_date, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTodoByID :one
SELECT
    t.id,
    t.user_id,
    t.list_id,
    t.title,
    t.description,
    t.completed,
    t.priority,
    t.due_date,
    t.created_at,
    t.updated_at,
    t.completed_at,
    t.deleted_at,
    l.name as list_name
FROM todos t
LEFT JOIN todo_lists l ON t.list_id = l.id AND l.deleted_at IS NULL
WHERE t.id = ? AND t.deleted_at IS NULL;

-- name: GetTodosByUserID :many
SELECT
    t.id,
    t.user_id,
    t.list_id,
    t.title,
    t.description,
    t.completed,
    t.priority,
    t.due_date,
    t.created_at,
    t.updated_at,
    t.completed_at,
    t.deleted_at,
    l.name as list_name
FROM todos t
LEFT JOIN todo_lists l ON t.list_id = l.id AND l.deleted_at IS NULL
WHERE t.user_id = ? AND t.deleted_at IS NULL
ORDER BY t.created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateTodo :exec
UPDATE todos
SET title = ?, description = ?, completed = ?, priority = ?, due_date = ?, updated_at = ?, completed_at = ?, list_id = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteTodo :exec
UPDATE todos
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at IS NULL;

-- name: CountTodosByUser :one
SELECT COUNT(*)
FROM todos
WHERE user_id = ? AND deleted_at IS NULL;

-- name: CountTodos :one
SELECT COUNT(*)
FROM todos
WHERE deleted_at IS NULL;

-- name: GetTodosFiltered :many
SELECT id, user_id, list_id, title, description, completed, priority, due_date, created_at, updated_at, completed_at, deleted_at
FROM todos
WHERE deleted_at IS NULL
  AND (sqlc.narg('user_id') IS NULL OR user_id = sqlc.narg('user_id'))
  AND (sqlc.narg('completed') IS NULL OR completed = sqlc.narg('completed'))
  AND (sqlc.narg('priority') IS NULL OR priority = sqlc.narg('priority'))
  AND (sqlc.narg('due_date_from') IS NULL OR due_date >= sqlc.narg('due_date_from'))
  AND (sqlc.narg('due_date_to') IS NULL OR due_date <= sqlc.narg('due_date_to'))
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountTodosFiltered :one
SELECT COUNT(*)
FROM todos
WHERE deleted_at IS NULL
  AND (sqlc.narg('user_id') IS NULL OR user_id = sqlc.narg('user_id'))
  AND (sqlc.narg('completed') IS NULL OR completed = sqlc.narg('completed'))
  AND (sqlc.narg('priority') IS NULL OR priority = sqlc.narg('priority'))
  AND (sqlc.narg('due_date_from') IS NULL OR due_date >= sqlc.narg('due_date_from'))
  AND (sqlc.narg('due_date_to') IS NULL OR due_date <= sqlc.narg('due_date_to'));
