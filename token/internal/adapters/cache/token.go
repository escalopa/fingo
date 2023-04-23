package cache

import (
	"context"
	"encoding/json"

	oteltracer "github.com/escalopa/fingo/token/internal/adapters/tracer"

	"github.com/escalopa/fingo/token/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
)

type TokenRepository struct {
	c *redis.Client
}

// NewTokenRepositoryV1 creates a new token repository
func NewTokenRepositoryV1(client *redis.Client) *TokenRepository {
	return &TokenRepository{c: client}
}

// GetTokenPayload gets the token payload from cache
func (tr *TokenRepository) GetTokenPayload(ctx context.Context, accessToken string) (*core.TokenPayload, error) {
	ctx, span := oteltracer.Tracer().Start(ctx, "GetTokenPayload")
	defer span.End()
	// Get token payload from cache
	bytes, err := tr.c.Get(ctx, accessToken).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.B(err).Code(errs.NotFound).Msg("token not found").Err()
		}
		return nil, errs.B(err).Code(errs.Internal).Msg("failed to get token").Err()
	}
	// Unmarshal token payload
	var payload core.TokenPayload
	err = json.Unmarshal([]byte(bytes), &payload)
	if err != nil {
		return nil, errs.B(err).Code(errs.Internal).Msg("failed to unmarshal token payload").Err()
	}
	return &payload, nil
}
