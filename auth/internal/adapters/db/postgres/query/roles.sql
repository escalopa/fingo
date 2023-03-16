-- name: CreateRole :exec
INSERT INTO roles (name)
VALUES ($1);

-- name: GetRoleByName :one
SELECT *
FROM roles
WHERE name = $1
LIMIT 1;

-- name: GetUserRoles :many
SELECT r.name
FROM user_roles ur
       JOIN roles r on r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: GrantRoleToUser :execrows
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2);

-- name: RevokeRoleFromUser :execrows
DELETE
FROM user_roles
WHERE user_id = $1
  AND role_id = $2;

-- name: HasPrivillage :one
SELECT COUNT(*)
FROM user_roles ur
       JOIN roles r on r.id = ur.id
WHERE user_id = $1
  AND r.name = $2;
