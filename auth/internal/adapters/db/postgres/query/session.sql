-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, refresh_token, user_agent, client_ip, expires_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetSessionByID :one
SELECT *
FROM sessions
WHERE id = $1
LIMIT 1;

-- name: GetUserSessions :many
SELECT *
FROM sessions
WHERE user_id = $1;

-- name: UpdateSessionRefreshToken :execrows
UPDATE sessions
SET refresh_token = $2,
    expires_at    = $3
WHERE id = $1;

-- name: SetSessionIsBlocked :execrows
UPDATE sessions
SET is_blocked = $2
WHERE id = $1;

-- name: GetUserDevices :many
SELECT user_agent, client_ip
FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteSessionByID :execrows
DELETE
FROM sessions
WHERE id = $1;
