package mypostgres

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/escalopa/fingo/auth/internal/adapters/db/postgres/sqlc"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type UserRepository struct {
	q db.Querier
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{q: db.New(conn)}
}

func (ur *UserRepository) CreateUser(ctx context.Context, arg core.CreateUserParams) error {
	err := ur.q.CreateUser(ctx, db.CreateUserParams{
		ID:             arg.ID,
		FirstName:      arg.FirstName,
		LastName:       arg.LastName,
		Username:       arg.Username,
		Gender:         arg.Gender,
		Email:          arg.Email,
		PhoneNumber:    arg.Phone,
		HashedPassword: arg.HashedPassword,
		Birthday:       arg.BirthDate,
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
	gender, err := parseGender(user.Gender)
	if err != nil {
		return core.User{}, err
	}
	return core.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Phone:           user.PhoneNumber,
		Username:        user.Username,
		Gender:          gender,
		Email:           user.Email,
		IsEmailVerified: user.IsVerifiedEmail,
		IsPhoneVerified: user.IsVerifiedPhone,
		BirthDate:       user.Birthday,
		CreatedAt:       user.CreatedAt,
	}, nil
}

func parseGender(gender interface{}) (string, error) {
	// Read gender as bytes
	byteValue, ok := gender.([]uint8)
	if !ok {
		return "", errs.B().Msg(fmt.Sprintf("invalid gender type: %v", gender)).Err()
	}
	// Convert gender to string
	var strValue string
	for _, v := range byteValue {
		strValue += string(rune(v))
	}
	// Check gender acceptable values
	switch strValue {
	case "MALE":
		return "MALE", nil
	case "FEMALE":
		return "FEMALE", nil
	default:
		return "", errs.B().Msg("unknown gender type").Err()
	}
}
