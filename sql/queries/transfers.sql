-- name: CreateTransfer :one
INSERT INTO transfers (id, sender_id, recipient_id, amount)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: GetTransfers :many
SELECT * FROM transfers
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;