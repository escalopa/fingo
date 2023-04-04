package db

import (
	"context"
	"database/sql"

	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/lordvidex/errs"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// CreateAccount creates an account for a user
func (r *AccountRepository) CreateAccount(ctx context.Context, params core.CreateAccountParams) (int64, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "AccountRepository.CreateAccount")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Get currency id
	currencyID, err := q.GetCurrencyByName(ctx, tx, params.Currency.String())
	if err != nil {
		return 0, errs.B(err).Code(errs.NotFound).Msg("failed to get currency id").Err()
	}
	// Create account
	res, err := q.CreateAccount(ctx, tx, sqlc.CreateAccountParams{
		UserID:     params.UserID,
		Name:       params.Name,
		CurrencyID: currencyID,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return 0, errorUniqueViolation(err, "account with this uuid already exists")
		} else {
			return 0, errorQuery(err, "failed to create account")
		}
	}
	// Get account id from result
	accountID, err := res.LastInsertId()
	if err != nil {
		return 0, errorQuery(err, "failed to get account id from result")
	}
	return accountID, nil
}

// GetAccount returns account for given account id
func (r *AccountRepository) GetAccount(ctx context.Context, accountID int64) (core.Account, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "AccountRepository.GetAccount")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Account{}, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	account, err := q.GetAccount(ctx, tx, accountID)
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
	ctx, span := oteltracer.Tracer().Start(ctx, "AccountRepository.GetAccounts")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Get accounts of a user by his id
	accounts, err := q.GetAccounts(ctx, tx, userID)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, errorNotFound(err, "no accounts found")
		}
		return nil, errorQuery(err, "failed to get accounts")
	}
	res := make([]core.Account, len(accounts))
	for i, account := range accounts {
		res[i] = fromDBAccountToAccount(account)
	}
	return res, nil
}

// DeleteAccount deletes account by given id
func (r *AccountRepository) DeleteAccount(ctx context.Context, accountID int64) error {
	ctx, span := oteltracer.Tracer().Start(ctx, "AccountRepository.DeleteAccount")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := sqlc.New()
	// Delete account
	err = q.DeleteAccount(ctx, tx, accountID)
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "account not found")
		}
		return errorQuery(err, "failed to delete account")
	}
	return nil
}

// fromDBAccountToAccount converts sqlc.Account to core.Account
func fromDBAccountToAccount(dbAccount sqlc.Account) core.Account {
	return core.Account{
		ID:       dbAccount.ID,
		Name:     dbAccount.Name,
		Currency: convertCurrencyByID(dbAccount.CurrencyID),
		Balance:  dbAccount.Balance,
	}
}

// convertCurrencyByID converts currency id to core.Currency
// TODO: make a sql query to get currency name by given id
func convertCurrencyByID(id int16) core.Currency {
	switch id {
	case 1:
		return core.CurrencyUSD
	case 2:
		return core.CurrencyEGP
	case 3:
		return core.CurrencyEUR
	case 4:
		return core.CurrencyGBP
	default:
		return core.CurrencyRUB
	}
}
