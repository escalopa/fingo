package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTransactionType(t *testing.T) {

	test := []struct {
		name string
		t    string
		want TransactionType
	}{
		{
			name: "transfer lower case",
			t:    "transfer",
			want: TransactionTypeTransfer,
		},
		{
			name: "transfer upper case",
			t:    "TRANSFER",
			want: TransactionTypeTransfer,
		},
		{
			name: "deposit lower case",
			t:    "deposit",
			want: TransactionTypeDeposit,
		},
		{
			name: "deposit upper case",
			t:    "DEPOSIT",
			want: TransactionTypeDeposit,
		},
		{
			name: "withdrawal lower case",
			t:    "withdrawal",
			want: TransactionTypeWithdrawal,
		},
		{
			name: "withdrawal upper case",
			t:    "WITHDRAWAL",
			want: TransactionTypeWithdrawal,
		},
		{
			name: "empty",
			t:    "",
			want: "",
		},
		{
			name: "unknown",
			t:    "unknown",
			want: "",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ParseTransactionType(tt.t))
		})
	}
}
