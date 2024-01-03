-- name: CreateUser :one
INSERT INTO users (id, harshed_password, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY full_name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users SET harshed_password = $2, password_changed_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;