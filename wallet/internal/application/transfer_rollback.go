package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type TransferRollbackParams struct {
	TransactionID string `validate:"required,uuid"`
}

type TransferRollbackCommand interface {
	Execute(ctx context.Context, params TransferRollbackParams) error
}

type TransferRollbackCommandImpl struct {
	v  Validator
	l  Locker
	ur UserRepository
	ar AccountRepository
	tr TransactionRepository
}

func (c *TransferRollbackCommandImpl) Execute(ctx context.Context, params TransferRollbackParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "TransferRollbackCommand.Execute")
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
		// Parse transaction id
		transactionID, err := uuid.Parse(params.TransactionID)
		if err != nil {
			return errs.B().Code(errs.InvalidArgument).Msg("invalid transaction id").Err()
		}
		transaction, err := c.tr.GetTransaction(ctx, transactionID)
		if err != nil {
			return err
		}
		// Check if the transaction type is transfer
		if transaction.Type.String() != core.TransactionTypeTransfer.String() {
			return errs.B().Code(errs.InvalidArgument).Msg("can't rollback a non transfer transaction").Err()
		}
		// Check if the transaction is already rolled back
		if transaction.IsRolledBack {
			return errs.B().Code(errs.InvalidArgument).Msg("transaction already rolled back").Err()
		}
		fromAccount, err := c.ar.GetAccount(ctx, transaction.FromAccountID)
		if err != nil {
			return err
		}
		// Check if the transaction creator is the same as caller
		if innerID != fromAccount.ID {
			return errorNotAccountOwner
		}
		// Lock accounts
		unlock := c.l.Lock(ctx, transaction.FromAccountID, transaction.ToAccountID)
		defer unlock()
		// Rollback transaction
		err = c.tr.RollbackTransaction(ctx, transactionID)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewTransferRollbackCommand(
	v Validator,
	l Locker,
	ur UserRepository,
	ar AccountRepository,
	tr TransactionRepository,
) TransferRollbackCommand {
	return &TransferRollbackCommandImpl{v: v, l: l, ur: ur, ar: ar, tr: tr}
}
