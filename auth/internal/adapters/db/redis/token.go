package redis

import (
	"context"
	"time"

	"github.com/lordvidex/errs"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/go-redis/redis/v9"
)

// TokenRepository is a redis repository for tokens implementing the TokenRepository interface
type TokenRepository struct {
	r  *redis.Client
	td time.Duration // token duration
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(client *redis.Client, opts ...func(*TokenRepository)) *TokenRepository {
	tr := &TokenRepository{r: client}
	for _, opt := range opts {
		opt(tr)
	}
	return tr
}

// WithTokenDuration sets the token duration
func WithTokenDuration(td time.Duration) func(*TokenRepository) {
	return func(tr *TokenRepository) {
		tr.td = td
	}
}

// Store stores a token mapped to a session id
func (tr *TokenRepository) Store(ctx context.Context, token string, params core.TokenPayload) error {
	ctx, span := tracer.Tracer().Start(ctx, "TokenRepository.Store")
	defer span.End()
	if token == "" {
		return errs.B(nil).Code(errs.InvalidArgument).Msg("token cannot be empty").Err()
	}
	err := tr.r.Set(ctx, token, params, tr.td).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to store token").Err()
	}
	return nil
}

// Delete deletes a token from the cache
func (tr *TokenRepository) Delete(ctx context.Context, token string) error {
	ctx, span := tracer.Tracer().Start(ctx, "TokenRepository.Delete")
	defer span.End()
	if token == "" {
		return errs.B(nil).Code(errs.InvalidArgument).Msg("token cannot be empty").Err()
	}
	err := tr.r.Del(ctx, token).Err()
	if err != nil {
		return errs.B(err).Code(errs.Internal).Msg("failed to delete token").Err()
	}
	return nil
}
