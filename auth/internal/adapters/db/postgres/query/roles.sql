-- name: CreateRole :exec
INSERT INTO roles (name)
VALUES ($1);

-- name: GetUserRoles :many
SELECT r.name
FROM user_roles ur
       JOIN roles r on r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: GetRoleUsers :many
SELECT u.id
FROM user_roles ur
       JOIN users u on u.id = ur.user_id
WHERE ur.role_id = $1;

-- name: UpdateRole :exec
UPDATE roles
SET name = $2
WHERE id = $1;

-- name: DeleteRole :exec
DELETE
FROM roles
WHERE id = $1;
