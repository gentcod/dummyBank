-- name: CreateAccount :one
INSERT INTO accounts (id, owner, balance, currency, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccounts :many
SELECT * FROM accounts
ORDER BY owner
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;