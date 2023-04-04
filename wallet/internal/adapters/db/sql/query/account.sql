-- name: CreateAccount :exec
INSERT INTO accounts (user_id, currency_id, name)
VALUES ($1, $2, $3);

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccounts :many
SELECT *
FROM accounts
WHERE user_id = $1;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = $1;
