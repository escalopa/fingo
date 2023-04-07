-- name: GetCurrencyByID :one
SELECT *
FROM currency
WHERE id = $1
LIMIT 1;

-- name: GetCurrencyByName :one
SELECT id
FROM currency
WHERE name = $1
LIMIT 1;
