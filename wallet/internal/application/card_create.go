package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type CreateCardParams struct {
	AccountID int64 `validate:"required,min=1"`
}

type CreateCardCommand interface {
	Execute(ctx context.Context, params CreateCardParams) error
}

type CreateCardCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
	cr CardRepository
	ng CardNumberGenerator
}

func (c *CreateCardCommandImpl) Execute(ctx context.Context, params CreateCardParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "CreateCardCommand.Execute")
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
		// Create any random card number
		number, err := c.ng.GenCardNumber(ctx)
		if err != nil {
			return err
		}
		// Add new card to account
		err = c.cr.CreateCard(ctx, core.CreateCardParams{
			AccountID: params.AccountID,
			Number:    number,
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func NewCreateCardCommand(
	v Validator,
	ur UserRepository,
	ar AccountRepository,
	cr CardRepository,
	ng CardNumberGenerator,
) CreateCardCommand {
	return &CreateCardCommandImpl{v: v, ur: ur, ar: ar, cr: cr, ng: ng}
}
