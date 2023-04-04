-- name: CreateAccount :execresult
INSERT INTO accounts (user_id, currency_id, name)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccounts :many
SELECT *
FROM accounts
WHERE user_id = $1;

-- name: AddAccountBalance :exec
UPDATE accounts
SET balance = balance + $2
WHERE id = $1;

-- name: SubAccountBalance :exec
UPDATE accounts
SET balance = balance + $2
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = $1;
