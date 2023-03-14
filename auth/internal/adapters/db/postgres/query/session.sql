-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, access_token, refresh_token, user_agent, client_ip, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetSessionByID :one
SELECT *
FROM sessions
WHERE id = $1
LIMIT 1;

-- name: GetUserSessions :many
SELECT *
FROM sessions
WHERE user_id = $1;

-- name: UpdateSessionTokens :execrows
UPDATE sessions
SET access_token  = $2,
    refresh_token = $3,
    expires_at    = $4
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
