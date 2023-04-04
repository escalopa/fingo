package db

import (
	"context"
	"database/sql"

	db "github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	oteltracer "github.com/escalopa/fingo/wallet/internal/adapters/tracer"
	"github.com/escalopa/fingo/wallet/internal/core"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// CreateCard creates a new card in the database
func (r *CardRepository) CreateCard(ctx context.Context, params core.CreateCardParams) (int64, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "CardRepo.CreateCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := db.New()
	// Create card
	res, err := q.CreateCard(ctx, tx, db.CreateCardParams{
		AccountID: params.AccountID,
		Number:    params.Number,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return 0, errorUniqueViolation(err, "card with this number already exists")
		} else {
			return 0, errorQuery(err, "failed to create card")
		}
	}
	// Get card id
	cardID, err := res.LastInsertId()
	if err != nil {
		return 0, errorQuery(err, "failed to get card id from result")
	}
	return int64(cardID), nil
}

// GetCard returns a card for a given number
func (r *CardRepository) GetCard(ctx context.Context, cardNumber string) (core.Card, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "CardRepo.GetCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Card{}, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := db.New()
	// Get card
	card, err := q.GetCard(ctx, tx, cardNumber)
	if err != nil {
		if IsNotFoundError(err) {
			return core.Card{}, errorNotFound(err, "card not found")
		} else {
			return core.Card{}, errorQuery(err, "failed to get card")
		}
	}
	return fromDBCardToCard(card), nil
}

// GetCards returns all cards for a given account
func (r *CardRepository) GetCards(ctx context.Context, accountID int64) ([]core.Card, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "CardRepo.GetCards")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := db.New()
	// Get cards
	cards, err := q.GetAccountCards(ctx, tx, accountID)
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
	ctx, span := oteltracer.Tracer().Start(ctx, "CardRepo.DeleteCard")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer deferTx(tx, &err)
	q := db.New()
	// Delete card
	err = q.DeleteCard(ctx, tx, cardNumber)
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
func fromDBCardToCard(card db.Card) core.Card {
	return core.Card{
		AccountID: card.AccountID,
		Number:    card.Number,
	}
}
