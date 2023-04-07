-- name: CreateCard :exec
INSERT INTO cards (number, account_id)
VALUES ($1, $2);

-- name: GetCard :one
SELECT *
FROM cards
WHERE number = $1;

-- name: GetAccountCards :many
SELECT *
FROM cards
WHERE account_id = $1;

-- name: GetCardAccount :one
SELECT a.id as id, a.user_id as owner_id, a.name, a.balance, cc.name as currency
FROM cards c
       JOIN accounts a on a.id = c.account_id
       JOIN currency cc on a.currency_id = cc.id
WHERE c.number = $1
LIMIT 1;

-- name: GetUserCards :many
SELECT cards.*, accounts.currency_id
FROM cards
       INNER JOIN accounts ON accounts.id = cards.account_id
WHERE accounts.user_id = $1;

-- name: GetCardBalance :one
SELECT balance
FROM accounts
WHERE id = (SELECT account_id
            FROM cards
            WHERE number = $1);

-- name: DeleteCard :exec
DELETE
FROM cards
WHERE number = $1;

-- name: DeleteAccountCards :exec
DELETE
FROM cards
WHERE account_id = $1;
