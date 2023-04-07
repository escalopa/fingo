package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type CreateAccountParams struct {
	Currency string `validate:"required"`
	Name     string `validate:"required"`
}

type CreateAccountCommand interface {
	Execute(ctx context.Context, params CreateAccountParams) error
}

type CreateAccountCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
}

func (c *CreateAccountCommandImpl) Execute(ctx context.Context, params CreateAccountParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := oteltracer.Tracer().Start(ctx, "CreateAccountCommand.Execute")
		defer span.End()
		// Validate params
		if err := c.v.Validate(ctx, params); err != nil {
			return err
		}
		// Get user external id
		userID, err := contextutils.GetUserID(ctx)
		if err != nil {
			return err
		}
		// Get the inner user id
		innerID, err := c.ur.GetUser(ctx, userID)
		if err != nil {
			return err
		}
		// Create a new account
		err = c.ar.CreateAccount(ctx, core.CreateAccountParams{
			UserID:   innerID,
			Name:     params.Name,
			Currency: core.ParseCurrency(params.Currency),
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func NewCreateAccountCommand(v Validator, ur UserRepository, ar AccountRepository) CreateAccountCommand {
	return &CreateAccountCommandImpl{v: v, ur: ur, ar: ar}
}
