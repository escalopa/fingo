package redis

import (
	"context"
	"encoding/json"
	"github.com/escalopa/gochat/contact/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
	"time"
)

type CodeRepository struct {
	r   *redis.Client
	exp time.Duration
}

// NewCodeRepository creates a new code repository
// The default expiration time is 5 minutes
func NewCodeRepository(client *redis.Client, opts ...func(*CodeRepository)) *CodeRepository {
	cr := &CodeRepository{
		r:   client,
		exp: 5 * time.Minute,
	}
	for _, opt := range opts {
		opt(cr)
	}
	return cr
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
func (cr *CodeRepository) Save(ctx context.Context, email string, vc core.VerificationCode) error {
	err := cr.r.Set(ctx, email, vc, cr.exp).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to save code").Err()
	}
	return nil
}

// Get returns the user id associated with the given code if it exists
func (cr *CodeRepository) Get(ctx context.Context, email string) (core.VerificationCode, error) {
	value, err := cr.r.Get(ctx, email).Result()
	if err != nil {
		// if the email does not exist, return a not found error
		if err == redis.Nil {
			return core.VerificationCode{}, errs.B(err).Code(errs.NotFound).Msg("email not found").Err()
		}
		// otherwise, return an internal error
		return core.VerificationCode{}, errs.B(err).Code(errs.Internal).Msg("failed to get email").Err()
	}
	var vc core.VerificationCode
	err = json.Unmarshal([]byte(value), &vc)
	if err != nil {
		return core.VerificationCode{}, errs.B(err).Code(errs.Internal).Msg("failed to unmarshal code").Err()
	}
	return vc, nil
}

func (cr *CodeRepository) Close() error {
	return cr.r.Close()
}
