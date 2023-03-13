-- name: CreateUser :exec
INSERT INTO users (id,
                   first_name,
                   last_name,
                   username,
                   gender,
                   birthday,
                   email,
                   phone_number,
                   hashed_password,
                   created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now());

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

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone_number = $1
LIMIT 1;

-- name: SetUserEmailIsVerified :exec
UPDATE users
SET is_verified_email = $2
WHERE id = $1;

-- name: SetUserPhoneIsVerified :exec
UPDATE users
SET is_verified_phone = $2
WHERE id = $1;

-- name: ChangeUserEmail :exec
UPDATE users
SET email = $2
WHERE id = $1;

-- name: ChangeUserPhone :exec
UPDATE users
SET phone_number = $2
WHERE id = $1;

-- name: ChangeUserPassword :exec
UPDATE users
SET hashed_password     = $2,
    password_changed_at = now()
WHERE id = $1;

-- name: ChangeNames :exec
UPDATE users
SET first_name = coalesce(sqlc.narg('first_name'), first_name),
    last_name  = coalesce(sqlc.narg('last_name'), last_name),
    username   = coalesce(sqlc.narg('username'), username),
    birthday   = coalesce(sqlc.narg('birthday'), birthday)
WHERE id = sqlc.arg('id');

-- name: DeleteUserByID :exec
DELETE
FROM users
WHERE id = $1;
