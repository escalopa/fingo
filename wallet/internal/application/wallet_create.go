package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
)

type CreateWalletParams struct{}

type CreateWalletCommand interface {
	Execute(ctx context.Context, params CreateWalletParams) error
}

type CreateWalletCommandImpl struct {
	v  Validator
	ur UserRepository
}

func (c *CreateWalletCommandImpl) Execute(ctx context.Context, params CreateWalletParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "CreateWalletCommand.Execute")
		defer span.End()
		// Validate params
		if err := c.v.Validate(ctx, params); err != nil {
			return err
		}
		// Parse userID from context
		userID, err := contextutils.GetUserID(ctx)
		if err != nil {
			return err
		}
		// Create a new wallet for the given user-id
		err = c.ur.CreateUser(ctx, userID)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewCreateWalletCommand(v Validator, ur UserRepository) CreateWalletCommand {
	return &CreateWalletCommandImpl{v: v, ur: ur}
}
