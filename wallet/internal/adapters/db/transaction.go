package db

import (
	"context"
	"database/sql"

	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Transfer transfers money from one account to another
func (r *TransactionRepository) Transfer(ctx context.Context, params core.CreateTransactionParams) error {
	ctx, span := tracer.Tracer().Start(ctx, "TransactionRepository.CreateTransaction")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Add money to destination account
	err = q.AddAccountBalance(ctx, tx, sqlc.AddAccountBalanceParams{
		ID:      params.ToAccountID,
		Balance: params.Amount,
	})
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "destination account not found")
		} else {
			return errorQuery(err, "failed to add money to destination account")
		}
	}
	// Subtract money from source account
	err = q.SubAccountBalance(ctx, tx, sqlc.SubAccountBalanceParams{
		ID:      params.FromAccountID,
		Balance: params.Amount,
	})
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "source account not found")
		} else {
			return errorQuery(err, "failed to subtract money from source account")
		}
	}
	// Create transaction
	err = q.CreateTransferTransaction(ctx, tx, sqlc.CreateTransferTransactionParams{
		SourceAccountID:      sql.NullInt64{Int64: params.FromAccountID, Valid: true},
		DestinationAccountID: sql.NullInt64{Int64: params.ToAccountID, Valid: true},
		Amount:               params.Amount,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return errorUniqueViolation(err, "transaction id already exists")
		} else {
			return errorQuery(err, "failed to create transaction")
		}
	}
	return nil
}

// Deposit adds money to an account and creates a transaction
func (r *TransactionRepository) Deposit(ctx context.Context, params core.CreateTransactionParams) error {
	ctx, span := tracer.Tracer().Start(ctx, "TransactionRepository.Deposit")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Add money to source account
	err = q.AddAccountBalance(ctx, tx, sqlc.AddAccountBalanceParams{
		ID:      params.ToAccountID,
		Balance: params.Amount,
	})
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "source account not found")
		} else {
			return errorQuery(err, "failed to subtract money from source account")
		}
	}
	// Create transaction
	err = q.CreateDepositTransaction(ctx, tx, sqlc.CreateDepositTransactionParams{
		DestinationAccountID: sql.NullInt64{Int64: params.ToAccountID, Valid: true},
		Amount:               params.Amount,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return errorUniqueViolation(err, "transaction id already exists")
		} else {
			return errorQuery(err, "failed to create transaction")
		}
	}
	return nil
}

// Withdraw subtracts money from an account and creates a transaction
func (r *TransactionRepository) Withdraw(ctx context.Context, params core.CreateTransactionParams) error {
	ctx, span := tracer.Tracer().Start(ctx, "TransactionRepository.Withdraw")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Subtract money from source account
	err = q.SubAccountBalance(ctx, tx, sqlc.SubAccountBalanceParams{
		ID:      params.FromAccountID,
		Balance: params.Amount,
	})
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "source account not found")
		} else {
			return errorQuery(err, "failed to subtract money from source account")
		}
	}
	// Create transaction
	err = q.CreateWithdrawTransaction(ctx, tx, sqlc.CreateWithdrawTransactionParams{
		SourceAccountID: sql.NullInt64{Int64: params.FromAccountID, Valid: true},
		Amount:          params.Amount,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return errorUniqueViolation(err, "transaction id already exists")
		} else {
			return errorQuery(err, "failed to create transaction")
		}
	}
	return nil
}

// GetTransaction returns a transaction by its ID
func (r *TransactionRepository) GetTransaction(ctx context.Context, transactionID uuid.UUID) (core.Transaction, error) {
	ctx, span := tracer.Tracer().Start(ctx, "TransactionRepository.GetTransaction")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Transaction{}, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	transaction, err := q.GetTransaction(ctx, tx, transactionID)
	if err != nil {
		if IsNotFoundError(err) {
			return core.Transaction{}, errorNotFound(err, "transaction not found")
		} else {
			return core.Transaction{}, errorQuery(err, "failed to get transaction")
		}
	}
	coreTx := fromDBTransactionRowToTransaction(transaction)
	return coreTx, nil
}

// GetTransactions returns a list of transactions for a given account with filters(pagination, date range, etc)
func (r *TransactionRepository) GetTransactions(ctx context.Context, params core.GetTransactionsParams) ([]core.Transaction, error) {
	ctx, span := tracer.Tracer().Start(ctx, "TransactionRepository.GetTransactions")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Get transactions
	//var t sqlc.NullTransactionType
	//err = t.Scan(params.Type.String())
	//if err != nil {
	//	return nil, errorQuery(err, "failed to convert transaction type")
	//}
	transactions, err := q.GetTransactions(ctx, tx, sqlc.GetTransactionsParams{
		AccountID: params.AccountID,
		Limit:     params.Limit,
		Offset:    params.Offset,
		MinAmount: sql.NullFloat64{Float64: params.MinAmount, Valid: params.MinAmount > 0},
		MaxAmount: sql.NullFloat64{Float64: params.MaxAmount, Valid: params.MaxAmount > 0},
		//TransactionType: t,
	})
	if err != nil {
		return nil, errorQuery(err, "failed to get transactions")
	}
	// Convert to core transactions
	res := make([]core.Transaction, len(transactions))
	for i, transaction := range transactions {
		res[i] = fromDBTransactionsRowToTransaction(transaction)
	}
	return res, nil
}

// RollbackTransaction deletes a transaction and restores the balance of the involved accounts
func (r *TransactionRepository) RollbackTransaction(ctx context.Context, transactionID uuid.UUID) error {
	ctx, span := tracer.Tracer().Start(ctx, "TransactionRepository.RollbackTransaction")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Get transaction
	transaction, err := q.GetTransaction(ctx, tx, transactionID)
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "transaction not found")
		} else {
			return errorQuery(err, "failed to get transaction")
		}
	}
	if transaction.Type != sqlc.TransactionTypeTransfer {
		return errorRollbackUnsupported
	}
	// Set transaction as rolled back
	err = q.SetTransactionRolledBack(ctx, tx, transactionID)
	if err != nil {
		return errorQuery(err, "failed to set transaction as rolled back")
	}
	// Add money to source account
	err = q.AddAccountBalance(ctx, tx, sqlc.AddAccountBalanceParams{
		ID:      transaction.FromAccountID.Int64,
		Balance: transaction.Amount,
	})
	if err != nil {
		return errorQuery(err, "failed to add money to source account")
	}
	// Subtract money from destination account
	err = q.SubAccountBalance(ctx, tx, sqlc.SubAccountBalanceParams{
		ID:      transaction.ToAccountID.Int64,
		Balance: transaction.Amount,
	})
	if err != nil {
		return errorQuery(err, "failed to subtract money from destination account")
	}
	return nil
}

// fromDBTransactionRowToTransaction converts a sqlc.GetTransactionRow to a core.Transaction
func fromDBTransactionRowToTransaction(t sqlc.GetTransactionRow) core.Transaction {
	return core.Transaction{
		ID:              t.ID,
		Amount:          t.Amount,
		Type:            fromDBTransactionTypeToTransactionType(t.Type),
		ToAccountID:     t.ToAccountID.Int64,
		FromAccountName: convertATM(t.FromAccountName),
		ToAccountName:   convertATM(t.ToAccountName),
		CreatedAt:       t.CreatedAt,
		IsRolledBack:    t.IsRolledBack,
	}
}

// fromDBTransactionsRowToTransaction converts a sqlc.GetTransactionsRow to a core.Transaction
func fromDBTransactionsRowToTransaction(t sqlc.GetTransactionsRow) core.Transaction {
	return core.Transaction{
		ID:              t.ID,
		Amount:          t.Amount,
		Type:            fromDBTransactionTypeToTransactionType(t.Type),
		FromAccountName: convertATM(t.FromAccountName),
		ToAccountName:   convertATM(t.ToAccountName),
		CreatedAt:       t.CreatedAt,
		IsRolledBack:    t.IsRolledBack,
	}
}

// fromDBTransactionTypeToTransactionType converts a sqlc.TransactionType to a core.TransactionType
func fromDBTransactionTypeToTransactionType(t sqlc.TransactionType) core.TransactionType {
	switch t {
	case sqlc.TransactionTypeDeposit:
		return core.TransactionTypeDeposit
	case sqlc.TransactionTypeWithdrawal:
		return core.TransactionTypeWithdrawal
	case sqlc.TransactionTypeTransfer:
		return core.TransactionTypeTransfer
	default:
		return ""
	}
}

// convertATM converts a sql.NullString to a string. If the string is null, it returns "ATM"
func convertATM(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "ATM"
}
