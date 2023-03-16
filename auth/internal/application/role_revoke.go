package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"time"
)

// RevokeRoleParams contains the parameters for the RevokeRoleCommand
type RevokeRoleParams struct {
	UserID   string `validate:"required,uuid"`
	RoleName string `validate:"required,min=3,max=30"`
}

// RevokeRoleCommand is the interface for the RevokeRoleCommandImpl
type RevokeRoleCommand interface {
	Execute(ctx context.Context, params RevokeRoleParams) error
}

// RevokeRoleCommandImpl is the implementation of the RevokeRoleCommand
type RevokeRoleCommandImpl struct {
	v  Validator
	rr RoleRepository
}

// Execute executes the RevokeRoleCommand with the given parameters
func (c *RevokeRoleCommandImpl) Execute(ctx context.Context, params RevokeRoleParams) error {
	return executeWithContextTimeout(ctx, 5*time.Second, func() error {
		if err := c.v.Validate(params); err != nil {
			return err
		}
		// Read user id from context
		callerID, err := parseUserIDFromContext(ctx)
		if err != nil {
			return err
		}
		// Parse userID in form of uuid
		var userID uuid.UUID
		userID, err = uuid.Parse(params.UserID)
		if err != nil {
			return err
		}
		//Check if the caller is an admin user
		isAdmin, err := c.rr.HasPrivillage(ctx, core.HasPrivillageParams{
			UserID:   callerID,
			RoleName: "admin",
		})
		if err != nil {
			return err
		}
		if !isAdmin {
			return errs.B().Code(errs.Forbidden).Msg("not enough privillage to revoke role").Err()
		}
		// Revoke user the given role
		err = c.rr.RevokeRole(ctx, core.RevokeRoleFromUserParams{UserID: userID, RoleName: params.RoleName})
		if err != nil {
			return err
		}
		return nil
	})
}

// NewRevokeRoleCommand returns a new RevokeRoleCommand with the passed dependencies
func NewRevokeRoleCommand(v Validator, rr RoleRepository) RevokeRoleCommand {
	return &RevokeRoleCommandImpl{v: v, rr: rr}
}
