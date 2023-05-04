package db

import (
	"context"
	"database/sql"

	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/adapters/db/sql/sqlc"
	"github.com/google/uuid"
)

type UserRepository struct {
	q  *sqlc.Queries
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db, q: sqlc.New()}
}

// CreateUser creates a new user in the database with the given uuid
func (r *UserRepository) CreateUser(ctx context.Context, uuid uuid.UUID) error {
	ctx, span := tracer.Tracer().Start(ctx, "UserRepo.CreateUser")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Create user
	err = r.q.CreateUser(ctx, tx, uuid)
	if err != nil {
		if IsUniqueViolationError(err) {
			return errorUniqueViolation(err, "user with this uuid already exists")
		} else {
			return errorQuery(err, "failed to create user")
		}
	}
	return nil
}

// GetUser returns the user id for the given uuid, Where uuid is the global user id between services
func (r *UserRepository) GetUser(ctx context.Context, uuid uuid.UUID) (int64, error) {
	ctx, span := tracer.Tracer().Start(ctx, "UserRepo.GetUser")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, errorTxNotStarted(err)
	}
	defer func() { err = deferTx(tx, err) }()
	// Get user id
	userID, err := r.q.GetUserByExternalID(ctx, tx, uuid)
	if err != nil {
		if IsNotFoundError(err) {
			return 0, errorNotFound(err, "user not found")
		} else {
			return 0, errorQuery(err, "failed to get user")
		}
	}
	return userID, nil
}
