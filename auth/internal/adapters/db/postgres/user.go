package mypostgres

import (
	"context"
	"database/sql"

	db "github.com/escalopa/fingo/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type UserRepository struct {
	q db.Querier
}

func NewUserRepository(conn *sql.DB) (*UserRepository, error) {
	if conn == nil {
		return nil, errs.B().Msg("passed connection is nil").Err()
	}
	return &UserRepository{q: db.New(conn)}, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, arg core.CreateUserParams) error {
	err := ur.q.CreateUser(ctx, db.CreateUserParams{
		ID:             arg.ID,
		FirstName:      arg.FirstName,
		LastName:       arg.LastName,
		Username:       arg.Username,
		Email:          arg.Email,
		HashedPassword: arg.HashedPassword,
	})
	if err != nil {
		if isUniqueViolationError(err) {
			return errs.B(err).Code(errs.AlreadyExists).Msg("user already exists").Err()
		}
		return errs.B(err).Code(errs.Internal).Msg("failed to create user").Err()
	}
	return nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (core.User, error) {
	user, err := ur.q.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.User{}, errs.B(err).Code(errs.NotFound).Msgf("no user found with the given id, id: %s", id).Err()
		}
		return core.User{}, errs.B(err).Code(errs.Internal).Msgf("failed to get user with id: %s", id).Err()
	}
	return fromDbUserToCore(user)
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (core.User, error) {
	user, err := ur.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.User{}, errs.B(err).Code(errs.NotFound).Msgf("no user found with the given email, email: %s", email).Err()
		}
		return core.User{}, errs.B(err).Code(errs.Internal).Msgf("failed to get user with email: %s", email).Err()
	}
	return fromDbUserToCore(user)
}

func (ur *UserRepository) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	rows, err := ur.q.DeleteUserByID(ctx, id)
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msgf("failed to delete user with id: %s", id).Err()
	}
	if rows == 0 {
		return errs.B(err).Code(errs.NotFound).Msgf("no user found with the given id, id: %s", id).Err()
	}
	return nil
}

func fromDbUserToCore(user db.User) (core.User, error) {
	return core.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Username:        user.Username,
		Email:           user.Email,
		HashedPassword:  user.HashedPassword,
		IsEmailVerified: user.IsVerifiedEmail,
		CreatedAt:       user.CreatedAt,
	}, nil
}
