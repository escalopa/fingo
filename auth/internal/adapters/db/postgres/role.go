package mypostgres

import (
	"context"
	"database/sql"
	db "github.com/escalopa/fingo/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

// RolesRepository is a repository for roles that implements the RolesRepository interface
type RolesRepository struct {
	q db.Querier
}

// NewRolesRepository returns a new RolesRepository
func NewRolesRepository(conn *sql.DB, opts ...func(*RolesRepository)) *RolesRepository {
	rr := &RolesRepository{q: db.New(conn)}
	for _, opt := range opts {
		opt(rr)
	}
	return rr
}

// CreateRole creates a new role with a given name
func (rr *RolesRepository) CreateRole(ctx context.Context, name string) error {
	err := rr.q.CreateRole(ctx, name)
	if err != nil {
		if isUniqueViolationError(err) {
			return errs.B(err).Code(errs.AlreadyExists).Msg("role already exists").Err()
		}
		return errs.B(err).Code(errs.Internal).Msg("failed to create role").Err()
	}
	return nil
}

// GetRoleByName returns a role by its name
func (rr *RolesRepository) GetRoleByName(ctx context.Context, name string) (core.Role, error) {
	role, err := rr.q.GetRoleByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.Role{}, errs.B(err).Code(errs.NotFound).Msg("role not found").Err()
		}
		return core.Role{}, errs.B(err).Code(errs.Internal).Msg("failed to get role by name").Err()
	}
	return core.Role{ID: role.ID, Name: role.Name}, nil
}

// GetUserRoles returns the roles of a user
func (rr *RolesRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]core.Role, error) {
	roles, err := rr.q.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, errs.B(err).Code(errs.Internal).Msg("failed to get user roles").Err()
	}
	// Convert the roles to core.Role
	coreRoles := make([]core.Role, len(roles))
	for i := 0; i < len(roles); i++ {
		coreRoles[i] = core.Role{Name: roles[i]}
	}
	return coreRoles, nil
}

// GrantRole grants a role to a user
func (rr *RolesRepository) GrantRole(ctx context.Context, params core.GrantRoleToUserParams) error {
	rows, err := rr.q.GrantRoleToUser(ctx, db.GrantRoleToUserParams{
		UserID: params.UserID,
		RoleID: params.RoleID,
	})
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to grant role to user").Err()
	}
	if rows == 0 {
		return errs.B(err).Code(errs.NotFound).Msg("failed to grant role to user, no rows affected").Err()
	}
	return nil
}

// RevokeRole revokes a role from a user
func (rr *RolesRepository) RevokeRole(ctx context.Context, params core.RevokeRoleFromUserParams) error {
	rows, err := rr.q.RevokeRoleFromUser(ctx, db.RevokeRoleFromUserParams{
		UserID: params.UserID,
		RoleID: params.RoleID,
	})
	if err != nil {
		return errs.B(err).Code(errs.NotFound).Msg("failed to grant role to user, no rows affected").Err()
	}
	if rows == 0 {
		return errs.B(err).Code(errs.Internal).Msg("failed to revoke role from user, no rows affected").Err()
	}
	return nil
}
