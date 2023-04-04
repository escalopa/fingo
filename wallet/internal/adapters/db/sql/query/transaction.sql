-- name: CreateTransaction :exec
INSERT INTO transactions (type, amount, source_account_id, destination_account_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetTransaction :one
SELECT transactions.id,
       type,
       amount,
       source.id        as from_account_id,
       source.name      as from_account_name,
       destination.id   as to_account_id,
       destination.name as to_account_name,
       created_at
FROM transactions
       JOIN accounts destination on destination.id = transactions.destination_account_id
       JOIN accounts source on destination.id = transactions.source_account_id
WHERE transactions.id = $1;

-- name: GetTransactions :many
SELECT transactions.id,
       transactions.amount,
       type,
       source.name      as from_account_name,
       destination.name as to_account_name,
       created_at
FROM transactions
       JOIN accounts destination on destination.id = transactions.destination_account_id
       JOIN accounts source on destination.id = transactions.source_account_id
WHERE source.id = sqlc.arg('account_id')
   OR destination.id = sqlc.arg('account_id')
  AND coalesce(sqlc.narg('transaction_type'), type) = type
  AND coalesce(sqlc.narg('from_date'), created_at) = created_at
  AND coalesce(sqlc.narg('to_date'), created_at) = created_at
  AND coalesce(sqlc.narg('from_amount'), transactions.amount) <= transactions.amount
  AND coalesce(sqlc.narg('to_amount'), transactions.amount) >= transactions.amount
  AND coalesce(sqlc.narg('is_rolled_back'), transactions.is_rolled_back) = transactions.is_rolled_back
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: SetTransactionRolledBack :exec
UPDATE transactions
SET is_rolled_back = true
WHERE id = $1;
