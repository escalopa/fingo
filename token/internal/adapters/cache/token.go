package cache

import (
	"context"
	"encoding/json"

	"github.com/escalopa/fingo/token/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
)

type TokeRepositoryImplV1 struct {
	c *redis.Client
}

// NewTokenRepositoryV1 creates a new token repository
func NewTokenRepositoryV1(client *redis.Client) (*TokeRepositoryImplV1, error) {
	if client == nil {
		return nil, errs.B().Code(errs.InvalidArgument).Msg("nil client").Err()
	}
	return &TokeRepositoryImplV1{c: client}, nil
}

// GetTokenPayload gets the token payload from cache
func (tr *TokeRepositoryImplV1) GetTokenPayload(ctx context.Context, accessToken string) (*core.TokenPayload, error) {
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
