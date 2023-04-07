package core

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

func ParseTransactionType(t string) TransactionType {
	switch strings.ToLower(t) {
	case TransactionTypeTransfer.String():
		return TransactionTypeTransfer
	case TransactionTypeDeposit.String():
		return TransactionTypeDeposit
	case TransactionTypeWithdrawal.String():
		return TransactionTypeWithdrawal
	default:
		return ""
	}
}

func (t TransactionType) String() string {
	return string(t)
}

type CreateTransactionParams struct {
	Amount        float64
	FromAccountID int64
	ToAccountID   int64
}

type GetTransactionsParams struct {
	AccountID int64
	Offset    int32
	Limit     int32
	// Additional filters
	Type      TransactionType
	MinAmount float64
	MaxAmount float64
}

type SendTransactionSmsParams struct {
	CardNumber    string  `json:"card_number"`
	RecipientName string  `json:"recipient_name"`
	Amount        float64 `json:"amount"`
	Balance       float64 `json:"balance"`
}

type Transaction struct {
	ID              uuid.UUID       `json:"id"`
	Amount          float64         `json:"amount"`
	Type            TransactionType `json:"type"`
	FromAccountID   int64           `json:"from_account_id"`
	FromAccountName string          `json:"from_account_name"`
	ToAccountID     int64           `json:"to_account_id"`
	ToAccountName   string          `json:"to_account_name"`
	CreatedAt       time.Time       `json:"created_at"`
	IsRolledBack    bool            `json:"is_rolled_back"`
}
