package db

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/escalopa/fingo/wallet/internal/core"
	"github.com/google/uuid"
)

func TestNewTransactionRepository(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *TransactionRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_Transfer(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		params core.CreateTransactionParams
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
			r := &TransactionRepository{
				db: tt.fields.db,
			}
			if err := r.Transfer(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("TransactionRepository.Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionRepository_Deposit(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		params core.CreateTransactionParams
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
			r := &TransactionRepository{
				db: tt.fields.db,
			}
			if err := r.Deposit(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("TransactionRepository.Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionRepository_Withdraw(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		params core.CreateTransactionParams
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
			r := &TransactionRepository{
				db: tt.fields.db,
			}
			if err := r.Withdraw(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("TransactionRepository.Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionRepository_GetTransaction(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx           context.Context
		transactionID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    core.Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TransactionRepository{
				db: tt.fields.db,
			}
			got, err := r.GetTransaction(tt.args.ctx, tt.args.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionRepository.GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionRepository.GetTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_GetTransactions(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		params core.GetTransactionsParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []core.Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TransactionRepository{
				db: tt.fields.db,
			}
			got, err := r.GetTransactions(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionRepository.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionRepository.GetTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_RollbackTransaction(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx           context.Context
		transactionID uuid.UUID
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
			r := &TransactionRepository{
				db: tt.fields.db,
			}
			if err := r.RollbackTransaction(tt.args.ctx, tt.args.transactionID); (err != nil) != tt.wantErr {
				t.Errorf("TransactionRepository.RollbackTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fromDBTransactionRowToTransaction(t *testing.T) {
	type args struct {
		t sqlc.GetTransactionRow
	}
	tests := []struct {
		name string
		args args
		want core.Transaction
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDBTransactionRowToTransaction(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDBTransactionRowToTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDBTransactionsRowToTransaction(t *testing.T) {
	type args struct {
		t sqlc.GetTransactionsRow
	}
	tests := []struct {
		name string
		args args
		want core.Transaction
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDBTransactionsRowToTransaction(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDBTransactionsRowToTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDBTransactionTypeToTransactionType(t *testing.T) {
	type args struct {
		t sqlc.TransactionType
	}
	tests := []struct {
		name string
		args args
		want core.TransactionType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDBTransactionTypeToTransactionType(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDBTransactionTypeToTransactionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertATM(t *testing.T) {
	type args struct {
		s sql.NullString
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertATM(tt.args.s); got != tt.want {
				t.Errorf("convertATM() = %v, want %v", got, tt.want)
			}
		})
	}
}
