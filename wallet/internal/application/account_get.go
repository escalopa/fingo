package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type GetAccountsParams struct{}

type GetAccountsCommand interface {
	Execute(ctx context.Context, params GetAccountsParams) ([]core.Account, error)
}

type GetAccountsCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
}

func (c *GetAccountsCommandImpl) Execute(ctx context.Context, params GetAccountsParams) ([]core.Account, error) {
	var accounts []core.Account
	err := contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "GetAccountsCommand.Execute")
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
		innerID, err := c.ur.GetUser(ctx, userID)
		if err != nil {
			return err
		}
		// Get user accounts
		accounts, err = c.ar.GetAccounts(ctx, innerID)
		if err != nil {
			return err
		}
		return nil
	})
	return accounts, err
}

func NewGetAccountsCommand(v Validator, ur UserRepository, ar AccountRepository) GetAccountsCommand {
	return &GetAccountsCommandImpl{v: v, ur: ur, ar: ar}
}
