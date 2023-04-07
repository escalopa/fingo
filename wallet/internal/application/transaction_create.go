package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/lordvidex/errs"
)

type CreateTransactionParams struct {
	Amount   float64              `validate:"required,min=10"`
	Type     core.TransactionType `validate:"required"`
	FromCard string               `validate:"required,len=16,number"`
	ToCard   string               `validate:"omitempty,len=16,number"`
}

var (
	errorNoSufficientFunds = errs.B().Code(errs.InvalidArgument).Msg("no sufficient balance to perform transaction").Err()
)

type CreateTransactionCommand interface {
	Execute(ctx context.Context, params CreateTransactionParams) error
}

type CreateTransactionCommandImpl struct {
	v  Validator
	l  Locker
	ur UserRepository
	ar AccountRepository
	cr CardRepository
	tr TransactionRepository
}

func (c *CreateTransactionCommandImpl) Execute(ctx context.Context, params CreateTransactionParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := oteltracer.Tracer().Start(ctx, "CreateTransactionCommand.Execute")
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
		// Get inner user id
		innerID, err := c.ur.GetUser(ctx, userID)
		if err != nil {
			return err
		}
		// Get the `from card` which
		fromAccount, err := c.cr.GetCardAccount(ctx, params.FromCard)
		if err != nil {
			return err
		}
		// Check that the caller is the account from card owner
		if innerID != fromAccount.OwnerID {
			return errorNotAccountOwner
		}
		// Set the toAccountID it the transaction is
		switch params.Type {
		case core.TransactionTypeTransfer:
			// Check if the `from account` has enough balance to preform the transaction
			if fromAccount.Balance < params.Amount {
				return errorNoSufficientFunds
			}
			if params.ToCard == "" {
				return errs.B().Code(errs.InvalidArgument).Msg("transfer card receiver not set").Err()
			}
			// Get card's account from db
			toAccount, err := c.cr.GetCardAccount(ctx, params.ToCard)
			if err != nil {
				return err
			}
			// Check that the `from account` & `to account` have the same currency type
			if toAccount.Currency != fromAccount.Currency {
				return errs.B().Code(errs.InvalidArgument).
					Msgf("accounts currency mismatch, from currency: %s, to currency: %s", fromAccount.Currency, toAccount.Currency).
					Err()
			}
			// Check that the `from account` & `to account` are not the same
			if toAccount.ID == fromAccount.ID {
				return errs.B().Code(errs.InvalidArgument).Msg("cannot transfer to the same account").Err()
			}
			// Lock both `from account` & `to account` for transaction
			unlock := c.l.Lock(ctx, fromAccount.ID, toAccount.ID)
			defer unlock()
			err = c.tr.Transfer(ctx, core.CreateTransactionParams{
				Amount:        params.Amount,
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
			})
			if err != nil {
				return err
			}
		case core.TransactionTypeDeposit:
			unlock := c.l.Lock(ctx, fromAccount.ID, -1)
			defer unlock()
			err = c.tr.Deposit(ctx, core.CreateTransactionParams{
				Amount:        params.Amount,
				FromAccountID: 0,
				ToAccountID:   fromAccount.ID,
			})
			if err != nil {
				return err
			}
		case core.TransactionTypeWithdrawal:
			unlock := c.l.Lock(ctx, fromAccount.ID, -1)
			defer unlock()
			// Check if the `from account` has enough balance to preform the transaction
			if fromAccount.Balance < params.Amount {
				return errorNoSufficientFunds
			}
			err = c.tr.Withdraw(ctx, core.CreateTransactionParams{
				Amount:        params.Amount,
				FromAccountID: fromAccount.ID,
				ToAccountID:   0,
			})
			if err != nil {
				return err
			}
		default:
			return errs.B().Code(errs.InvalidArgument).Msg("transaction type must be set").Err()
		}
		return nil
	})
}

func NewCreateTransactionCommand(
	v Validator,
	l Locker,
	ur UserRepository,
	ar AccountRepository,
	cr CardRepository,
	tr TransactionRepository,
) CreateTransactionCommand {
	return &CreateTransactionCommandImpl{v: v, l: l, ur: ur, ar: ar, cr: cr, tr: tr}
}
