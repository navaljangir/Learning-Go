-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, full_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUserByID :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET full_name = $2, updated_at = $3
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, username, email, password_hash, full_name, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CheckUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL) AS exists;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL) AS exists;
