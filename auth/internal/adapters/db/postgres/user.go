package mypostgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	db "github.com/escalopa/fingo/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type UserRepository struct {
	q db.Querier
}

func NewUserRepository(conn *sql.DB) (*UserRepository, error) {
	return &UserRepository{q: db.New(conn)}, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, arg core.CreateUserParams) error {
	ctx, span := tracer.Tracer().Start(ctx, "UserRepository.CreateUser")
	defer span.End()
	err := ur.q.CreateUser(ctx, db.CreateUserParams{
		ID:             arg.ID,
		FirstName:      arg.FirstName,
		LastName:       arg.LastName,
		Username:       arg.Username,
		Email:          arg.Email,
		HashedPassword: arg.HashedPassword,
	})
	if err != nil {
		if IsUniqueViolationError(err) {
			return errs.B(err).Code(errs.AlreadyExists).Msg("user already exists").Err()
		}
		return errs.B(err).Code(errs.Internal).Msg("failed to create user").Err()
	}
	return nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (core.User, error) {
	ctx, span := tracer.Tracer().Start(ctx, "UserRepository.GetUserByID")
	defer span.End()
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
	ctx, span := tracer.Tracer().Start(ctx, "UserRepository.GetUserByEmail")
	defer span.End()
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
	ctx, span := tracer.Tracer().Start(ctx, "UserRepository.DeleteUserByID")
	defer span.End()
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

// IsUniqueViolationError checks if an error is a unique violation error
func IsUniqueViolationError(err error) bool {
	er, ok := err.(*pq.Error)
	return ok && er.Code == "23505"
}
