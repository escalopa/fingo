package mypostgres

import (
	"context"
	"database/sql"
	db "github.com/escalopa/gochat/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type UserRepository struct {
	db *sql.DB
	q  db.Querier
}

func NewUserRepository(conn *sql.DB) (*UserRepository, error) {
	return &UserRepository{db: conn}, nil
}

type CreateUserParams struct {
	ID             uuid.UUID
	Name           string
	Username       string
	Email          string
	HashedPassword string
}

type ChangeUserEmailParams struct {
	ID    uuid.UUID
	Email string
}

type ChangePasswordParams struct {
	ID             uuid.UUID
	HashedPassword string
}

func (ur *UserRepository) CreateUser(ctx context.Context, arg CreateUserParams) error {
	err := ur.q.CreateUser(ctx, db.CreateUserParams{
		ID:             arg.ID,
		Name:           arg.Name,
		Username:       arg.Username,
		Email:          arg.Email,
		HashedPassword: arg.HashedPassword,
	})
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (core.User, error) {
	user, err := ur.q.GetUserByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return core.User{}, errs.B(err).Msgf("no user found with the given id, id: %s", id).Err()
		}
		return core.User{}, err
	}
	return fromDbUserToCore(user), nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (core.User, error) {
	user, err := ur.q.GetUserByEmail(ctx, email)
	if err != nil {
		if isNotFoundError(err) {
			return core.User{}, errs.B(err).Msgf("no user found with the given email, email: %s", email).Err()
		}
		return core.User{}, err
	}
	return fromDbUserToCore(user), nil
}

func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (core.User, error) {
	user, err := ur.q.GetUserByUsername(ctx, username)
	if err != nil {
		if isNotFoundError(err) {
			return core.User{},
				errs.B(err).Msgf("no user found with the given username, username: %s", username).Err()
		}
		return core.User{}, err
	}
	return fromDbUserToCore(user), nil
}

func (ur *UserRepository) ChangeUserEmail(ctx context.Context, arg ChangeUserEmailParams) error {
	err := ur.q.ChangeUserEmail(ctx, db.ChangeUserEmailParams{
		ID:    arg.ID,
		Email: arg.Email,
	})
	if err != nil {
		if isNotFoundError(err) {
			return errs.B(err).Msgf("no user found with the given id, id: %s", arg.ID).Err()
		}
		return err
	}
	return nil
}

func (ur *UserRepository) ChangePassword(ctx context.Context, arg ChangePasswordParams) error {
	err := ur.q.ChangePassword(ctx, db.ChangePasswordParams{
		ID:             arg.ID,
		HashedPassword: arg.HashedPassword,
	})
	if err != nil {
		if isNotFoundError(err) {
			return errs.B(err).Msgf("no user found with the given id, id: %s", arg.ID).Err()
		}
		return err
	}
	return nil
}

func (ur *UserRepository) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	err := ur.q.DeleteUserByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return errs.B(err).Msgf("no user found with the given id, id: %s", id).Err()
		}
		return err
	}
	return nil
}

func fromDbUserToCore(user db.User) core.User {
	return core.User{
		ID:         user.ID,
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Password:   user.HashedPassword,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}
}
