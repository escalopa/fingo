-- name: CreateTransferTransaction :exec
INSERT INTO transactions (type, amount, source_account_id, destination_account_id)
VALUES ('transfer', $1, $2, $3);

-- name: CreateDepositTransaction :exec
INSERT INTO transactions(type, amount, destination_account_id)
VALUES ('deposit', $1, $2);

-- name: CreateWithdrawTransaction :exec
INSERT INTO transactions(type, amount, source_account_id)
VALUES ('withdrawal', $1, $2);

-- name: GetTransaction :one
SELECT t.id,
       t.type,
       t.amount,
       source.id        as from_account_id,
       source.name      as from_account_name,
       destination.id   as to_account_id,
       destination.name as to_account_name,
       t.created_at,
       t.is_rolled_back
FROM transactions t
       LEFT JOIN accounts destination on destination.id = t.destination_account_id
       LEFT JOIN accounts source on source.id = t.source_account_id
WHERE t.id = $1;

-- name: GetTransactions :many
SELECT t.id,
       t.amount,
       t.type,
       source.name      as from_account_name,
       destination.name as to_account_name,
       t.created_at,
       t.is_rolled_back
FROM transactions t
       LEFT JOIN accounts destination on destination.id = t.destination_account_id
       LEFT JOIN accounts source on source.id = t.source_account_id
WHERE (source.id = sqlc.arg('account_id')
  OR destination.id = sqlc.arg('account_id'))
--   AND coalesce(sqlc.narg('transaction_type') IS NULL, t.type) = t.type
  AND coalesce(sqlc.narg('min_amount'), t.amount) <= t.amount
  AND coalesce(sqlc.narg('max_amount'), t.amount) >= t.amount
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: SetTransactionRolledBack :exec
UPDATE transactions
SET is_rolled_back = true
WHERE id = $1;
