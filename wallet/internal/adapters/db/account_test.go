package db

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	userID := generateRandomUser(t)

	tests := []struct {
		name    string
		args    core.CreateAccountParams
		wantErr bool
	}{
		{
			name: "sucess",
			args: core.CreateAccountParams{
				UserID:   userID,
				Name:     gofakeit.Name(),
				Currency: core.CurrencyUSD,
			},
			wantErr: false,
		},
		{
			name: "unknown currency",
			args: core.CreateAccountParams{
				UserID:   userID,
				Name:     gofakeit.Name(),
				Currency: core.Currency("UNKOWN"),
			},
			wantErr: true,
		},
		{
			name: "user not found",
			args: core.CreateAccountParams{
				UserID:   gofakeit.Int64(),
				Name:     gofakeit.Name(),
				Currency: core.CurrencyUSD,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewAccountRepository(conn)
			err := r.CreateAccount(context.Background(), tt.args)
			require.Truef(t, (err != nil) == tt.wantErr, "AccountRepository.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestAccountRepository_GetAccount(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx       context.Context
		accountID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    core.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountRepository{
				db: tt.fields.db,
			}
			got, err := r.GetAccount(tt.args.ctx, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountRepository.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountRepository.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountRepository_GetAccounts(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []core.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountRepository{
				db: tt.fields.db,
			}
			got, err := r.GetAccounts(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountRepository.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountRepository.GetAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountRepository_DeleteAccount(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx       context.Context
		accountID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountRepository{
				db: tt.fields.db,
			}
			if err := r.DeleteAccount(tt.args.ctx, tt.args.accountID); (err != nil) != tt.wantErr {
				t.Errorf("AccountRepository.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fromDBAccountToAccount(t *testing.T) {
	type args struct {
		account sqlc.GetAccountRow
	}
	tests := []struct {
		name string
		args args
		want core.Account
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDBAccountToAccount(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDBAccountToAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDBAccountsToAccount(t *testing.T) {
	type args struct {
		account sqlc.GetAccountsRow
	}
	tests := []struct {
		name string
		args args
		want core.Account
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDBAccountsToAccount(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDBAccountsToAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateRandomUser(t *testing.T) int64 {
	// Create user
	externalID := gofakeit.UUID()
	ur := NewUserRepository(conn)
	err := ur.CreateUser(context.Background(), uuid.MustParse(externalID))
	require.NoError(t, err)

	// Get user internal id
	internalID, err := ur.GetUser(context.Background(), uuid.MustParse(externalID))
	require.NoError(t, err)
	require.NotZero(t, internalID)
	return internalID
}
