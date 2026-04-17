-- name: CreateUser :one
INSERT INTO users (firm_id, email, display_name, role, auth_method, password_hash, notification_frequency)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsersByFirmID :many
SELECT * FROM users
WHERE firm_id = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateUser :one
UPDATE users
SET role = COALESCE(sqlc.narg('role'), role),
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    notification_frequency = COALESCE(sqlc.narg('notification_frequency'), notification_frequency),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users SET is_active = false WHERE id = $1;
