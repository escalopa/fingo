package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type GetCardsParams struct {
	AccountID int64 `validate:"required,min=1"`
}

type GetCardsCommand interface {
	Execute(ctx context.Context, params GetCardsParams) ([]core.Card, error)
}

type GetCardsCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
	cr CardRepository
}

func (c *GetCardsCommandImpl) Execute(ctx context.Context, params GetCardsParams) ([]core.Card, error) {
	var cards []core.Card
	err := contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "GetCardsCommand.Execute")
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
		// Check if caller is account owner
		account, err := c.ar.GetAccount(ctx, params.AccountID)
		if err != nil {
			return err
		}
		if innerID != account.OwnerID {
			return errorNotAccountOwner
		}
		// Get user cards
		cards, err = c.cr.GetCards(ctx, params.AccountID)
		if err != nil {
			return err
		}
		return nil
	})
	return cards, err
}

func NewGetCardsCommand(v Validator, ur UserRepository, ar AccountRepository, cr CardRepository) GetCardsCommand {
	return &GetCardsCommandImpl{v: v, ur: ur, ar: ar, cr: cr}
}
