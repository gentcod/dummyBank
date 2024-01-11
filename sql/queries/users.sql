-- name: CreateUser :one
INSERT INTO users (id, username, full_name, email, harshed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET harshed_password = $2, password_changed_at = $3
WHERE id = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;