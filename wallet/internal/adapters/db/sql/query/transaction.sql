-- name: CreateTransaction :exec
INSERT INTO transactions (type, amount, source_account_id, destination_account_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetTransaction :one
SELECT transactions.id, type, amount, source_account_id, destination_account_id
FROM transactions
WHERE id = $1;

-- name: GetTransactionsByAccount :many
SELECT transactions.id,
       type,
       transactions.amount,
       source.name            as from_account_name,
       source_account_id      as from_account_id,
       destination.name       as to_account_name,
       destination_account_id as to_account_id,
       created_at
FROM transactions
       JOIN accounts destination on destination.id = transactions.destination_account_id
       JOIN accounts source on destination.id = transactions.source_account_id
WHERE (source_account_id = $1 OR destination_account_id = $1);

-- name: GetTransactionsByAccountAndType :many
SELECT transactions.id,
       type,
       transactions.amount,
       source.name            as from_account_name,
       source_account_id      as from_account_id,
       destination.name       as to_account_name,
       destination_account_id as to_account_id,
       transactions.created_at
FROM transactions
       JOIN accounts destination on destination.id = transactions.destination_account_id
       JOIN accounts source on destination.id = transactions.source_account_id
WHERE (source_account_id = $1 OR destination_account_id = $1)
  AND type = $2;

-- name: GetTransactionsByAccountAndDate :many
SELECT transactions.id,
       type,
       transactions.amount,
       source.name            as from_account_name,
       source_account_id      as from_account_id,
       destination.name       as to_account_name,
       destination_account_id as to_account_id,
       created_at
FROM transactions
       JOIN accounts destination on destination.id = transactions.destination_account_id
       JOIN accounts source on destination.id = transactions.source_account_id
WHERE (source_account_id = $1 OR destination_account_id = $1)
  AND created_at >= $2
  AND created_at <= $3;
