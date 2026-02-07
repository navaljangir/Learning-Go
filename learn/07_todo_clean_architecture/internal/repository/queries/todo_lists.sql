-- name: CreateTodoList :exec
INSERT INTO todo_lists (id, user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetTodoListByID :one
SELECT id, user_id, name, created_at, updated_at, deleted_at
FROM todo_lists
WHERE id = ? AND deleted_at IS NULL;

-- name: GetTodoListsByUserID :many
SELECT id, user_id, name, created_at, updated_at, deleted_at
FROM todo_lists
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateTodoList :exec
UPDATE todo_lists
SET name = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteTodoList :exec
UPDATE todo_lists
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at IS NULL;

-- name: CountTodoListsByUser :one
SELECT COUNT(*)
FROM todo_lists
WHERE user_id = ? AND deleted_at IS NULL;

-- name: GetTodosByListID :many
SELECT id, user_id, list_id, title, description, completed, priority, due_date, created_at, updated_at, completed_at, deleted_at
FROM todos
WHERE list_id = sqlc.arg('list_id') AND deleted_at IS NULL
ORDER BY created_at DESC;
