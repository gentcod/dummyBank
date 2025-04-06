-- name: CreateUser :one
INSERT INTO users (id, username, full_name, email, harshed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET 
   harshed_password = COALESCE(sqlc.narg(harshed_password), harshed_password), 
   full_name = COALESCE(sqlc.narg(full_name), full_name), 
   email = COALESCE(sqlc.narg(email), email),
   password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
   updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: GetUser :one
SELECT id, username, full_name, email, is_email_verified, created_at FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserWithPassword :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: VerifyUserEmail :one
UPDATE users 
SET 
   is_email_verified = true,
   updated_at = now()
WHERE email = $1
RETURNING is_email_verified, updated_at;