-- name: CreateCard :execresult
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
