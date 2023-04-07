package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type GetTransactionHistoryParams struct {
	AccountID       int64                `validate:"required,min=1"`
	Offset          int32                `validate:"min=0"`
	Limit           int32                `validate:"min=1,max=100"`
	MinAmount       float64              `validate:"omitempty,min=10"`
	MaxAmount       float64              `validate:"omitempty,min=10"`
	TransactionType core.TransactionType `validate:"omitempty"`
}

type GetTransactionHistoryCommand interface {
	Execute(ctx context.Context, params GetTransactionHistoryParams) ([]core.Transaction, error)
}

type GetTransactionHistoryCommandImpl struct {
	v  Validator
	ur UserRepository
	ar AccountRepository
	tr TransactionRepository
}

func (c *GetTransactionHistoryCommandImpl) Execute(ctx context.Context, params GetTransactionHistoryParams) ([]core.Transaction, error) {
	var transactions []core.Transaction
	err := contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := oteltracer.Tracer().Start(ctx, "GetTransactionHistory.Execute")
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
		account, err := c.ar.GetAccount(ctx, params.AccountID)
		if err != nil {
			return err
		}
		// Check if the caller is the account's owner
		if innerID != account.OwnerID {
			return errorNotAccountOwner
		}
		// Get user transactions
		transactions, err = c.tr.GetTransactions(ctx, core.GetTransactionsParams{
			AccountID: params.AccountID,
			Offset:    params.Offset,
			Limit:     params.Limit,
			Type:      params.TransactionType,
			MinAmount: params.MinAmount,
			MaxAmount: params.MaxAmount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return transactions, err
}

func NewGetTransactionHistoryCommand(v Validator, ur UserRepository, ar AccountRepository, tr TransactionRepository) GetTransactionHistoryCommand {
	return &GetTransactionHistoryCommandImpl{v: v, ur: ur, ar: ar, tr: tr}
}
