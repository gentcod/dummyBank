-- name: CreateEntry :one
INSERT INTO entries (id, account_id, amount, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;
