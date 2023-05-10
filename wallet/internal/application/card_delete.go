package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
)

type DeleteCardParams struct {
	CardNumber string `validate:"required,len=16"`
}

type DeleteCardCommand interface {
	Execute(ctx context.Context, params DeleteCardParams) error
}

type DeleteCardCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
	cr CardRepository
}

func (c *DeleteCardCommandImpl) Execute(ctx context.Context, params DeleteCardParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "DeleteCardCommand.Execute")
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
		// Get card's account from db
		account, err := c.cr.GetCardAccount(ctx, params.CardNumber)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		// Check if the caller is the account's owner
		if innerID != account.OwnerID {
			return errorNotAccountOwner
		}
		// Delete card
		err = c.cr.DeleteCard(ctx, params.CardNumber)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewDeleteCardCommand(v Validator, ur UserRepository, ar AccountRepository, cr CardRepository) DeleteCardCommand {
	return &DeleteCardCommandImpl{v: v, ur: ur, ar: ar, cr: cr}
}
