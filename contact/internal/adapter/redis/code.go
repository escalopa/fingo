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

func NewCodeRepository(client *redis.Client, opts ...func(*CodeRepository)) *CodeRepository {
	cr := &CodeRepository{r: client}
	for _, opt := range opts {
		opt(cr)
	}
	return cr
}

func WithCodeContext(ctx context.Context) func(*CodeRepository) {
	return func(cr *CodeRepository) {
		cr.c = ctx
	}
}

func WithExpiration(exp time.Duration) func(*CodeRepository) {
	return func(cr *CodeRepository) {
		cr.exp = exp
	}
}

func (cr *CodeRepository) Save(code string, userID string) error {
	err := cr.r.Set(cr.c, code, userID, cr.exp).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to save code").Err()
	}
	return nil
}

func (cr *CodeRepository) Get(code string) (string, error) {
	value, err := cr.r.Get(cr.c, code).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errs.B(err).Code(errs.NotFound).Msg("code not found").Err()
		}
		return "", errs.B(err).Code(errs.Internal).Msg("failed to get code").Err()
	}
	return value, nil
}
