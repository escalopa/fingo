package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	ac "github.com/escalopa/gochat/auth/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
)

type UserRepository struct {
	r *redis.Client
}

func NewUserRepository(client *redis.Client, opts ...func(*UserRepository)) *UserRepository {
	ur := &UserRepository{r: client}
	for _, opt := range opts {
		opt(ur)
	}
	return ur
}

func WithTimeout(timeout time.Duration) func(*UserRepository) {
	return func(ur *UserRepository) {
		ur.r = ur.r.WithTimeout(timeout)
	}
}

func (ur *UserRepository) Save(ctx context.Context, u ac.User) error {
	var err error
	// Check if user already exists
	_, err = ur.Get(ctx, u.Email)
	if err == nil {
		return errs.B().Code(errs.AlreadyExists).Msg("user already exists").Err()
	}
	// Generate user id
	u.ID, err = newUserID()
	if err != nil {
		return err
	}
	// Save user to cache
	err = ur.r.Set(ctx, u.Email, u, 0).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("error saving user").Err()
	}
	return nil
}

func (ur *UserRepository) Get(ctx context.Context, id string) (ac.User, error) {
	var u ac.User
	userStr, err := ur.r.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return u, errs.B(err).Code(errs.NotFound).Msg("user not found").Err()
		}
		return u, errs.B(err).Code(errs.Internal).Msg("error getting user").Err()
	}
	err = json.Unmarshal([]byte(userStr), &u)
	if err != nil {
		return ac.User{}, errs.B().Code(errs.Internal).Msg("failed to marshal user from cache").Err()
	}
	return u, nil
}

func (ur *UserRepository) Update(ctx context.Context, u ac.User) error {
	err := ur.r.Set(ctx, u.Email, u, 0).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("error updating user").Err()
	}
	return nil
}

func newUserID() (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return uuid.UUID{}, errs.B(err).Code(errs.Internal).Msg("error generating user id").Err()
	}
	return id, nil
}
