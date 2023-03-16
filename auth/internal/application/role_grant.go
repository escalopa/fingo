package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"time"
)

// GrantRoleParams contains the parameters for the GrantRoleCommand
type GrantRoleParams struct {
	UserID   string `validate:"required,uuid"`
	RoleName string `validate:"required,min=3,max=30"`
}

// GrantRoleCommand is the interface for the GrantRoleCommandImpl
type GrantRoleCommand interface {
	Execute(ctx context.Context, params GrantRoleParams) error
}

// GrantRoleCommandImpl is the implementation of the GrantRoleCommand
type GrantRoleCommandImpl struct {
	v  Validator
	rr RoleRepository
}

// Execute executes the GrantRoleCommand with the given parameters
func (c *GrantRoleCommandImpl) Execute(ctx context.Context, params GrantRoleParams) error {
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
		userID, err := uuid.Parse(params.UserID)
		if err != nil {
			return err
		}
		//Check if the user is an admin user
		isAdmin, err := c.rr.HasPrivillage(ctx, core.HasPrivillageParams{
			UserID:   callerID,
			RoleName: "admin",
		})
		if err != nil {
			return err
		}
		if !isAdmin {
			return errs.B().Code(errs.Forbidden).Msg("not enough privillage to grant role").Err()
		}
		// Grant user the given role
		err = c.rr.GrantRole(ctx, core.GrantRoleToUserParams{UserID: userID, RoleName: params.RoleName})
		if err != nil {
			return err
		}
		return nil
	})
}

// NewGrantRoleCommand returns a new GrantRoleCommand with the passed dependencies
func NewGrantRoleCommand(v Validator, rr RoleRepository) GrantRoleCommand {
	return &GrantRoleCommandImpl{v: v, rr: rr}
}
