package redis

import (
	"context"
	"github.com/google/uuid"
	"time"

	ac "github.com/escalopa/gofly/auth/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
)

type UserRepository struct {
	r *redis.Client
	c context.Context
}

func NewUserRepository(client *redis.Client, opts ...func(*UserRepository)) *UserRepository {
	ur := &UserRepository{r: client}
	for _, opt := range opts {
		opt(ur)
	}
	return ur
}

func WithUserContext(ctx context.Context) func(*UserRepository) {
	return func(ur *UserRepository) {
		ur.c = ctx
	}
}

func WithTimeout(timeout time.Duration) func(*UserRepository) {
	return func(ur *UserRepository) {
		ur.r = ur.r.WithTimeout(timeout)
	}
}

func (ur *UserRepository) Save(u ac.User) error {
	var err error
	// Check if user already exists
	_, err = ur.Get(u.Email)
	if err == nil {
		return errs.B().Code(errs.AlreadyExists).Msg("user already exists").Err()
	}
	// Generate user id
	u.ID, err = newUserID()
	if err != nil {
		return err
	}
	// Save user to cache
	err = ur.r.Set(ur.c, u.Email, u, 0).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("error saving user").Err()
	}
	return nil
}

func (ur *UserRepository) Get(id string) (ac.User, error) {
	var u ac.User
	err := ur.r.Get(ur.c, id).Scan(&u)
	if err != nil {
		if err == redis.Nil {
			return u, errs.B(err).Code(errs.NotFound).Msg("user not found").Err()
		}
		return u, errs.B(err).Code(errs.Internal).Msg("error getting user").Err()
	}
	return u, nil
}

func (ur *UserRepository) Update(u ac.User) error {
	err := ur.r.Set(ur.c, u.Email, u, 0).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("error updating user").Err()
	}
	return nil
}

func newUserID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", errs.B(err).Code(errs.Internal).Msg("error generating user id").Err()
	}
	return "user-" + id.String(), nil
}
