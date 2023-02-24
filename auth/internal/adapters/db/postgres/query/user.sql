-- name: CreateUser :exec
INSERT INTO users (id, name, username, email, hashed_password, is_verified)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: SetUserIsVerified :exec
UPDATE users
SET is_verified = $1
WHERE id = $1;

-- name: ChangeUserEmail :exec
UPDATE users
SET email       = $2,
    is_verified = false
WHERE id = $1;

-- name: ChangePassword :exec
UPDATE users
SET hashed_password = $2
WHERE id = $1;

-- name: ChangeNames :exec
UPDATE users
SET name     = coalesce(sqlc.narg('name'), name),
    username = coalesce(sqlc.narg('username'), name)
WHERE id = sqlc.arg('id');

-- name: DeleteUserByID :exec
DELETE
FROM users
WHERE id = $1;
