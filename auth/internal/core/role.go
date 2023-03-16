package core

import "github.com/google/uuid"

type Role struct {
	Name string
}

// ------------------------- Params -------------------------

type GrantRoleToUserParams struct {
	UserID   uuid.UUID
	RoleName string
}

type RevokeRoleFromUserParams struct {
	UserID   uuid.UUID
	RoleName string
}

type HasPrivillageParams struct {
	UserID   uuid.UUID
	RoleName string
}
