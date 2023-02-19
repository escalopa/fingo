package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
	"time"
)

type CodeRepository struct {
	r   *redis.Client
	c   context.Context
	exp time.Duration
}

// NewCodeRepository creates a new code repository
// The default expiration time is 5 minutes
func NewCodeRepository(client *redis.Client, opts ...func(*CodeRepository)) *CodeRepository {
	cr := &CodeRepository{
		r:   client,
		c:   context.Background(),
		exp: 5 * time.Minute,
	}
	for _, opt := range opts {
		opt(cr)
	}
	return cr
}

// WithCodeContext sets the context for the code repository
func WithCodeContext(ctx context.Context) func(*CodeRepository) {
	return func(cr *CodeRepository) {
		cr.c = ctx
	}
}

// WithExpiration sets the expiration time for the code
// The default expiration time is 5 minutes
func WithExpiration(exp time.Duration) func(*CodeRepository) {
	return func(cr *CodeRepository) {
		cr.exp = exp
	}
}

// Save saves the code in the redis database for the given `user id`
// The code is set to expire after `exp` seconds
func (cr *CodeRepository) Save(code string, userID string) error {
	err := cr.r.Set(cr.c, code, userID, cr.exp).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to save code").Err()
	}
	return nil
}

// Get returns the user id associated with the given code if it exists
func (cr *CodeRepository) Get(code string) (string, error) {
	value, err := cr.r.Get(cr.c, code).Result()
	if err != nil {
		// if the code does not exist, return a not found error
		if err == redis.Nil {
			return "", errs.B(err).Code(errs.NotFound).Msg("code not found").Err()
		}
		// otherwise, return an internal error
		return "", errs.B(err).Code(errs.Internal).Msg("failed to get code").Err()
	}
	return value, nil
}

func (cr *CodeRepository) Close() error {
	return cr.r.Close()
}
