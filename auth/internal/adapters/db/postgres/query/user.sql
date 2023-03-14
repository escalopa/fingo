-- name: CreateUser :exec
INSERT INTO users (id,
                   first_name,
                   last_name,
                   username,
                   gender,
                   birthday,
                   email,
                   phone_number,
                   hashed_password)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: DeleteUserByID :execrows
DELETE
FROM users
WHERE id = $1;
