-- name: CreateAccount :exec
INSERT INTO accounts (user_id, currency_id, balance, name)
VALUES ($1, $2, 0, $3)
RETURNING id;

-- name: GetAccount :one
SELECT a.id, a.user_id, a.name, a.balance, a.currency_id, c.name as currency_name
FROM accounts a
       JOIN currency c on a.currency_id = c.id
WHERE a.id = $1
LIMIT 1;

-- name: GetAccounts :many
SELECT a.id, a.user_id, a.name, a.balance, a.currency_id, c.name as currency_name
FROM accounts a
       JOIN currency c on a.currency_id = c.id
WHERE a.user_id = $1;

-- name: AddAccountBalance :exec
UPDATE accounts
SET balance = balance + $2
WHERE id = $1;

-- name: SubAccountBalance :exec
UPDATE accounts
SET balance = balance - $2
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = $1;
