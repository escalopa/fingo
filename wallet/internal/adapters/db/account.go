package db

import (
	"context"
	"database/sql"

	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/lordvidex/errs"
)

type AccountRepository struct {
	q  *sqlc.Queries
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db, q: sqlc.New()}
}

// CreateAccount creates an account for a user
func (r *AccountRepository) CreateAccount(ctx context.Context, params core.CreateAccountParams) error {
	ctx, span := tracer.Tracer().Start(ctx, "AccountRepository.CreateAccount")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Get currency id
	currencyID, err := r.q.GetCurrencyByName(ctx, tx, params.Currency.String())
	if err != nil {
		return errs.B(err).Code(errs.NotFound).Msg("failed to get currency id").Err()
	}
	// Create account
	err = r.q.CreateAccount(ctx, tx, sqlc.CreateAccountParams{
		UserID:     params.UserID,
		Name:       params.Name,
		CurrencyID: currencyID,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return errorUniqueViolation(err, "account with this uuid already exists")
		} else {
			return errorQuery(err, "failed to create account")
		}
	}
	return nil
}

// GetAccount returns account for given account id
func (r *AccountRepository) GetAccount(ctx context.Context, accountID int64) (core.Account, error) {
	ctx, span := tracer.Tracer().Start(ctx, "AccountRepository.GetAccount")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Account{}, errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	account, err := r.q.GetAccount(ctx, tx, accountID)
	if err != nil {
		if IsNotFoundError(err) {
			return core.Account{}, errorNotFound(err, "account not found")
		} else {
			return core.Account{}, errorQuery(err, "failed to get account")
		}
	}
	return fromDBAccountToAccount(account), nil
}

// GetAccounts returns all accounts for given user
func (r *AccountRepository) GetAccounts(ctx context.Context, userID int64) ([]core.Account, error) {
	ctx, span := tracer.Tracer().Start(ctx, "AccountRepository.GetAccounts")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Get accounts of a user by his id
	accounts, err := r.q.GetAccounts(ctx, tx, userID)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, errorNotFound(err, "no accounts found")
		}
		return nil, errorQuery(err, "failed to get accounts")
	}
	res := make([]core.Account, len(accounts))
	for i, account := range accounts {
		res[i] = fromDBAccountsToAccount(account)
	}
	return res, nil
}

// DeleteAccount deletes account by given id
func (r *AccountRepository) DeleteAccount(ctx context.Context, accountID int64) error {
	ctx, span := tracer.Tracer().Start(ctx, "AccountRepository.DeleteAccount")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Delete account
	err = r.q.DeleteAccount(ctx, tx, accountID)
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "account not found")
		}
		return errorQuery(err, "failed to delete account")
	}
	return nil
}

// fromDBAccountToAccount converts sqlc.Account to core.Account
func fromDBAccountToAccount(account sqlc.GetAccountRow) core.Account {
	return core.Account{
		ID:       account.ID,
		OwnerID:  account.UserID,
		Name:     account.Name,
		Currency: core.Currency(account.CurrencyName),
		Balance:  account.Balance,
	}
}

func fromDBAccountsToAccount(account sqlc.GetAccountsRow) core.Account {
	return core.Account{
		ID:       account.ID,
		OwnerID:  account.UserID,
		Name:     account.Name,
		Currency: core.Currency(account.CurrencyName),
		Balance:  account.Balance,
	}
}
