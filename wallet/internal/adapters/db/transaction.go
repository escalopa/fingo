package db

import (
	"context"
	"database/sql"

	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, params core.CreateTransactionParams) error {
	ctx, span := oteltracer.Tracer().Start(ctx, "TransactionRepository.CreateTransaction")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// If it's not a withdrawal, add money to destination account
	if params.Type != core.TransactionTypeWithdrawal {
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
	}
	// If it's not a deposit, subtract money from source account
	if params.Type != core.TransactionTypeDeposit {
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
	}
	// Create transaction
	err = q.CreateTransaction(ctx, tx, sqlc.CreateTransactionParams{
		SourceAccountID:      sql.NullInt64{Int64: params.FromAccountID, Valid: true},
		DestinationAccountID: sql.NullInt64{Int64: params.ToAccountID, Valid: true},
		Amount:               params.Amount,
		Type:                 fromTransactionTypeToDBTransactionType(params.Type),
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
	ctx, span := oteltracer.Tracer().Start(ctx, "TransactionRepository.GetTransaction")
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
	coreTx, _, _ := fromDBTransactionRowToTransaction(transaction)
	return coreTx, nil
}

// GetTransactions returns a list of transactions for a given account with filters(pagination, date range, etc)
func (r *TransactionRepository) GetTransactions(ctx context.Context, params core.GetTransactionsParams) ([]core.Transaction, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "TransactionRepository.GetTransactions")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Convert params to db params
	dbParam := sqlc.GetTransactionsParams{AccountID: params.AccountID, Limit: params.Limit, Offset: params.Offset}
	if params.IsRolledBack {
		dbParam.IsRolledBack = sql.NullBool{Bool: true, Valid: true}
	}
	if !params.FromDate.IsZero() {
		dbParam.FromDate = sql.NullTime{Time: params.FromDate, Valid: true}
	}
	if !params.ToDate.IsZero() {
		dbParam.ToDate = sql.NullTime{Time: params.ToDate, Valid: true}
	}
	if params.TransactionType != core.TransactionTypeUnknown {
		dbParam.TransactionType = sqlc.NullTransactionType{TransactionType: fromTransactionTypeToDBTransactionType(params.TransactionType), Valid: true}
	}
	if params.FromAmount > 0 {
		dbParam.FromAmount = sql.NullFloat64{Float64: params.FromAmount, Valid: true}
	}
	if params.ToAmount > 0 {
		dbParam.ToAmount = sql.NullFloat64{Float64: params.ToAmount, Valid: true}
	}
	// Get transactions
	transactions, err := q.GetTransactions(ctx, tx, dbParam)
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
	ctx, span := oteltracer.Tracer().Start(ctx, "TransactionRepository.RollbackTransaction")
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
	if transaction.Type == sqlc.TransactionTypeDeposit || transaction.Type == sqlc.TransactionTypeWithdrawal {
		return errorRollabckUnsupported
	}
	// Set transaction as rolled back
	err = q.SetTransactionRolledBack(ctx, tx, transactionID)
	if err != nil {
		return errorQuery(err, "failed to set transaction as rolled back")
	}
	// Add money to source account
	err = q.AddAccountBalance(ctx, tx, sqlc.AddAccountBalanceParams{
		ID:      transaction.FromAccountID,
		Balance: transaction.Amount,
	})
	if err != nil {
		return errorQuery(err, "failed to add money to source account")
	}
	// Subtract money from destination account
	err = q.SubAccountBalance(ctx, tx, sqlc.SubAccountBalanceParams{
		ID:      transaction.ToAccountID,
		Balance: transaction.Amount,
	})
	if err != nil {
		return errorQuery(err, "failed to subtract money from destination account")
	}
	return nil
}

// fromDBTransactionRowToTransaction converts a sqlc.GetTransactionRow to a core.Transaction
func fromDBTransactionRowToTransaction(t sqlc.GetTransactionRow) (core.Transaction, int64, int64) {
	tx := fromDBTransactionsRowToTransaction(sqlc.GetTransactionsRow{
		ID:              t.ID,
		Amount:          t.Amount,
		Type:            t.Type,
		FromAccountName: t.FromAccountName,
		ToAccountName:   t.ToAccountName,
		CreatedAt:       t.CreatedAt,
	})
	return tx, t.FromAccountID, t.ToAccountID
}

// fromDBTransactionsRowToTransaction converts a sqlc.GetTransactionsRow to a core.Transaction
func fromDBTransactionsRowToTransaction(t sqlc.GetTransactionsRow) core.Transaction {
	tx := core.Transaction{
		ID:        t.ID,
		Amount:    t.Amount,
		Type:      fromDBTransactionTypeToTransactionType(t.Type),
		CreatedAt: t.CreatedAt,
	}
	if t.FromAccountName != "" {
		tx.FromAccount = t.FromAccountName
	} else { // Deposit
		tx.FromAccount = "ATM"
	}
	if t.ToAccountName != "" {
		tx.ToAccount = t.ToAccountName
	} else { // Withdrawal
		tx.ToAccount = "ATM"
	}
	return tx
}

// fromTransactionTypeToDBTransactionType converts a core.TransactionType to a sqlc.TransactionType
func fromTransactionTypeToDBTransactionType(t core.TransactionType) sqlc.TransactionType {
	switch t {
	case core.TransactionTypeDeposit:
		return sqlc.TransactionTypeDeposit
	case core.TransactionTypeWithdrawal:
		return sqlc.TransactionTypeWithdrawal
	case core.TransactionTypeTransfer:
		return sqlc.TransactionTypeTransfer
	}
	return ""
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
	}
	return ""
}
