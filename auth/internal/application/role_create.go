package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/lordvidex/errs"
	"time"
)

// CreateRoleParams contains the parameters for the CreateRoleCommand
type CreateRoleParams struct {
	Name string `validate:"required,min=3,max=30"`
}

// CreateRoleCommand is the interface for the CreateRoleCommandImpl
type CreateRoleCommand interface {
	Execute(ctx context.Context, params CreateRoleParams) error
}

// CreateRoleCommandImpl is the implementation of the CreateRoleCommand
type CreateRoleCommandImpl struct {
	v  Validator
	rr RoleRepository
}

// Execute executes the CreateRoleCommand with the given parameters
func (c *CreateRoleCommandImpl) Execute(ctx context.Context, params CreateRoleParams) error {
	return executeWithContextTimeout(ctx, 5*time.Second, func() error {
		if err := c.v.Validate(params); err != nil {
			return err
		}
		// Read user id from context
		callerID, err := parseUserIDFromContext(ctx)
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
			return errs.B().Code(errs.Forbidden).Msg("not enough privillage to create role").Err()
		}
		// Create new role
		err = c.rr.CreateRole(ctx, params.Name)
		if err != nil {
			return err
		}
		return nil
	})
}

// NewCreateRoleCommand returns a new CreateRoleCommand with the passed dependencies
func NewCreateRoleCommand(v Validator, rr RoleRepository) CreateRoleCommand {
	return &CreateRoleCommandImpl{v: v, rr: rr}
}
