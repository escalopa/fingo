package core

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeUnknown    TransactionType = "unknown"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

func IsSupportedTransactionType(t TransactionType) bool {
	switch t {
	case TransactionTypeTransfer, TransactionTypeDeposit, TransactionTypeWithdrawal:
		return true
	default:
		return false
	}
}

func (t TransactionType) String() string {
	return string(t)
}

type CreateTransactionParams struct {
	Amount        float64
	Type          TransactionType
	FromAccountID int64
	ToAccountID   int64
}

type GetTransactionsParams struct {
	AccountID       int64
	TransactionType TransactionType
	FromDate        time.Time
	ToDate          time.Time
	FromAmount      float64
	ToAmount        float64
	IsRolledBack    bool
	Offset          int32
	Limit           int32
}

type SendTransactionSmsParams struct {
	CardNumber    string  `json:"card_number"`
	RecipientName string  `json:"recipient_name"`
	Amount        float64 `json:"amount"`
	Balance       float64 `json:"balance"`
}

type Transaction struct {
	ID          uuid.UUID       `json:"id"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	FromAccount string          `json:"from_account"`
	ToAccount   string          `json:"to_account"`
	CreatedAt   time.Time       `json:"created_at"`
}
