-- name: CreateTransfer :one
INSERT INTO transfers (id, sender_id, recipient_id, amount, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;