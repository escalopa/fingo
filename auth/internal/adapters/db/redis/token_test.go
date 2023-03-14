package redis

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenRepository_Store(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tokenDuration := 2 * time.Second
	tr := NewTokenRepository(testRedis, WithTokenDuration(tokenDuration))
	// Test cases
	testCases := []struct {
		name  string
		token string
		arg   core.TokenCache
	}{
		{
			name:  "store token",
			token: gofakeit.UUID(),
			arg: core.TokenCache{
				UserID: uuid.New().String(),
				Roles:  []string{"admin"},
			},
		},
	}
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Store(ctx, tc.token, tc.arg)
			require.NoError(t, err)
		})
	}
	// Wait for token to expire
	time.Sleep(tokenDuration)
	// Run test cases again to check if the token is deleted
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tr.r.Get(ctx, tc.token).Result()
			require.ErrorIs(t, err, redis.Nil)
		})
	}
}

func TestTokenRepository_Delete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tr := NewTokenRepository(testRedis, WithTokenDuration(1*time.Minute))
	// Test cases
	testCases := []struct {
		name  string
		token string
		arg   core.TokenCache
	}{
		{
			name:  "delete token",
			token: gofakeit.UUID(),
			arg: core.TokenCache{
				UserID: uuid.New().String(),
				Roles:  []string{"admin"},
			},
		},
	}
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Store(ctx, tc.token, tc.arg)
			require.NoError(t, err)
		})
	}
	// Run test cases again to check if the token is deleted
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Delete(ctx, tc.token)
			require.NoError(t, err)
		})
	}
	// Run test cases again to check if the token is deleted
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tr.r.Get(ctx, tc.token).Result()
			require.ErrorIs(t, err, redis.Nil)
		})
	}
}
