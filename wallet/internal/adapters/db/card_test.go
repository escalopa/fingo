package db

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	db "github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/escalopa/fingo/wallet/internal/core"
)

func TestNewCardRepository(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *CardRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardRepository_CreateCard(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		params core.CreateCardParams
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
			r := &CardRepository{
				db: tt.fields.db,
			}
			if err := r.CreateCard(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("CardRepository.CreateCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCardRepository_GetCard(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx        context.Context
		cardNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    core.Card
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardRepository{
				db: tt.fields.db,
			}
			got, err := r.GetCard(tt.args.ctx, tt.args.cardNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardRepository.GetCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardRepository.GetCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardRepository_GetCardAccount(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx        context.Context
		cardNumber string
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
			r := &CardRepository{
				db: tt.fields.db,
			}
			got, err := r.GetCardAccount(tt.args.ctx, tt.args.cardNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardRepository.GetCardAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardRepository.GetCardAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardRepository_GetCards(t *testing.T) {
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
		want    []core.Card
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardRepository{
				db: tt.fields.db,
			}
			got, err := r.GetCards(tt.args.ctx, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardRepository.GetCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardRepository.GetCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardRepository_DeleteCard(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx        context.Context
		cardNumber string
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
			r := &CardRepository{
				db: tt.fields.db,
			}
			if err := r.DeleteCard(tt.args.ctx, tt.args.cardNumber); (err != nil) != tt.wantErr {
				t.Errorf("CardRepository.DeleteCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fromDBCardToCard(t *testing.T) {
	type args struct {
		card db.Card
	}
	tests := []struct {
		name string
		args args
		want core.Card
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDBCardToCard(tt.args.card); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDBCardToCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
