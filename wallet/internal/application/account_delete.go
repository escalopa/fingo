package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/lordvidex/errs"
)

type DeleteAccountParams struct {
	AccountID int64 `validate:"required,min=1"`
}

type DeleteAccountCommand interface {
	Execute(ctx context.Context, params DeleteAccountParams) error
}

type DeleteAccountCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
}

func (c *DeleteAccountCommandImpl) Execute(ctx context.Context, params DeleteAccountParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := oteltracer.Tracer().Start(ctx, "DeleteAccountCommand.Execute")
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
		// Get the account from database
		account, err := c.ar.GetAccount(ctx, params.AccountID)
		if err != nil {
			return err
		}
		// Check if the caller is the account's owner
		if innerID != account.OwnerID {
			return errorNotAccountOwner
		}
		// Check that the account balance is not empty
		if account.Balance > 0 {
			return errs.B().Code(errs.InvalidArgument).Msg("account's balance must be empty to be deleted").Err()
		}
		// Delete account
		err = c.ar.DeleteAccount(ctx, params.AccountID)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewDeleteAccountCommand(v Validator, ur UserRepository, ar AccountRepository) DeleteAccountCommand {
	return &DeleteAccountCommandImpl{v: v, ur: ur, ar: ar}
}
