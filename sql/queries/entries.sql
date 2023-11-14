-- name: CreateEntry :one
INSERT INTO entries (id, account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: GetEntries :many
SELECT * FROM entries
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;