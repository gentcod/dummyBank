-- name: CreateAccount :one
INSERT INTO accounts (id, owner, balance, currency, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts
ORDER BY owner
LIMIT $1
OFFSET $2;