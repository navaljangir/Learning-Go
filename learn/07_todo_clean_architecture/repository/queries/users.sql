-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, full_name, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE username = ? AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE email = ? AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET full_name = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CheckUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND deleted_at IS NULL) AS username_exists;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL) AS email_exists;
