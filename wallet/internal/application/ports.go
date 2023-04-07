package application

import (
	"context"

	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, uuid uuid.UUID) error
	GetUser(ctx context.Context, uuid uuid.UUID) (int64, error)
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, params core.CreateAccountParams) error
	GetAccount(ctx context.Context, accountID int64) (core.Account, error)
	GetAccounts(ctx context.Context, userID int64) ([]core.Account, error)
	DeleteAccount(ctx context.Context, accountID int64) error
}

type CardRepository interface {
	CreateCard(ctx context.Context, params core.CreateCardParams) error
	GetCard(ctx context.Context, cardNumber string) (core.Card, error)
	GetCardAccount(ctx context.Context, cardNumber string) (core.Account, error)
	GetCards(ctx context.Context, accountID int64) ([]core.Card, error)
	DeleteCard(ctx context.Context, cardNumber string) error
}

type TransactionRepository interface {
	Transfer(ctx context.Context, params core.CreateTransactionParams) error
	Deposit(ctx context.Context, params core.CreateTransactionParams) error
	Withdraw(ctx context.Context, params core.CreateTransactionParams) error
	GetTransaction(ctx context.Context, transactionID uuid.UUID) (core.Transaction, error)
	GetTransactions(ctx context.Context, params core.GetTransactionsParams) ([]core.Transaction, error)
	RollbackTransaction(ctx context.Context, transactionID uuid.UUID) error
}

// CardNumberGenerator is an interface for generating card numbers
type CardNumberGenerator interface {
	GenCardNumber(ctx context.Context) (string, error)
}

// SmsSender is an interface for sending sms about transactions completion
type SmsSender interface {
	SendTransactionSms(ctx context.Context, params core.SendTransactionSmsParams) error
}

// Locker locks account ids for applying transactions
// The function returned is used to free the lock
type Locker interface {
	Lock(ctx context.Context, x any, y ...any) func()
}

// Validator is an interface for validating structs using tags
type Validator interface {
	Validate(ctx context.Context, s interface{}) (err error)
}
