-- name: CreateUser :exec
INSERT INTO users(external_id)
VALUES ($1)
RETURNING id;

-- name: GetUserByExternalID :one
SELECT id
FROM users
WHERE external_id = $1;
