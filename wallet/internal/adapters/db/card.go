package db

import (
	"context"
	"database/sql"

	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type CardRepository struct {
	q  *sqlc.Queries
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db, q: sqlc.New()}
}

// CreateCard creates a new card in the database
func (r *CardRepository) CreateCard(ctx context.Context, params core.CreateCardParams) error {
	ctx, span := tracer.Tracer().Start(ctx, "CardRepository.CreateCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Create card
	err = r.q.CreateCard(ctx, tx, sqlc.CreateCardParams{
		AccountID: params.AccountID,
		Number:    params.Number,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return errorUniqueViolation(err, "card with this number already exists")
		} else {
			return errorQuery(err, "failed to create card")
		}
	}
	return nil
}

// GetCard returns a card for a given number
func (r *CardRepository) GetCard(ctx context.Context, cardNumber string) (core.Card, error) {
	ctx, span := tracer.Tracer().Start(ctx, "CardRepository.GetCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Card{}, errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Get card
	card, err := r.q.GetCard(ctx, tx, cardNumber)
	if err != nil {
		if IsNotFoundError(err) {
			return core.Card{}, errorNotFound(err, "card not found")
		} else {
			return core.Card{}, errorQuery(err, "failed to get card")
		}
	}
	return fromDBCardToCard(card), nil
}

func (r *CardRepository) GetCardAccount(ctx context.Context, cardNumber string) (core.Account, error) {
	ctx, span := tracer.Tracer().Start(ctx, "CardRepository.GetCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Account{}, errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	result, err := r.q.GetCardAccount(ctx, tx, cardNumber)
	if err != nil {
		if IsNotFoundError(err) {
			return core.Account{}, errorNotFound(err, "card not found")
		} else {
			return core.Account{}, errorQuery(err, "failed to get card's account")
		}
	}
	account := core.Account{
		ID:       result.ID,
		OwnerID:  result.OwnerID,
		Name:     result.Name,
		Balance:  result.Balance,
		Currency: core.Currency(result.Currency),
	}
	return account, nil
}

// GetCards returns all cards for a given account
func (r *CardRepository) GetCards(ctx context.Context, accountID int64) ([]core.Card, error) {
	ctx, span := tracer.Tracer().Start(ctx, "CardRepository.GetCards")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Get cards
	cards, err := r.q.GetAccountCards(ctx, tx, accountID)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, errorNotFound(err, "cards not found")
		} else {
			return nil, errorQuery(err, "failed to get cards")
		}
	}
	// Convert to core.Card
	res := make([]core.Card, len(cards))
	for i, card := range cards {
		res[i] = fromDBCardToCard(card)
	}
	return res, nil
}

// DeleteCard deletes a card for a given number
func (r *CardRepository) DeleteCard(ctx context.Context, cardNumber string) error {
	ctx, span := tracer.Tracer().Start(ctx, "CardRepository.DeleteCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Delete card
	err = r.q.DeleteCard(ctx, tx, cardNumber)
	if err != nil {
		if IsNotFoundError(err) {
			return errorNotFound(err, "card not found")
		} else {
			return errorQuery(err, "failed to delete card")
		}
	}
	return nil
}

// fromDBCardToCard converts a db.Card to a core.Card
func fromDBCardToCard(card sqlc.Card) core.Card {
	return core.Card{
		AccountID: card.AccountID,
		Number:    card.Number,
	}
}
