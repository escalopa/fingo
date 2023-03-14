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
		name      string
		token     string
		arg       core.TokenCache
		wantError bool
	}{
		{
			name:  "success",
			token: gofakeit.UUID(),
			arg: core.TokenCache{
				UserID: uuid.New().String(),
				Roles:  []string{"admin"},
			},
			wantError: false,
		},
		{
			name:  "empty token",
			token: "",
			arg: core.TokenCache{
				UserID: uuid.New().String(),
				Roles:  []string{"admin"},
			},
			wantError: true,
		},
	}
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Store(ctx, tc.token, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantError {
				t.Error("expected error but got nil")
			}
			if err == nil {
				_, err = tr.r.Get(ctx, tc.token).Result()
				require.NoError(t, err)
			}
		})
	}
	// Wait for token to expire
	time.Sleep(tokenDuration)
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
		name      string
		token     string
		arg       core.TokenCache
		wantError bool
	}{
		{
			name:  "success",
			token: gofakeit.UUID(),
			arg: core.TokenCache{
				UserID: uuid.New().String(),
				Roles:  []string{"admin"},
			},
			wantError: false,
		},
		{
			name:      "empty token",
			token:     "",
			arg:       core.TokenCache{},
			wantError: true,
		},
	}
	// Run test cases again to check if the token is deleted
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Store(ctx, tc.token, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantError {
				t.Error("expected error but got nil")
			}
			if err == nil {
				_, err = tr.r.Get(ctx, tc.token).Result()
				require.NoError(t, err)
				err = tr.Delete(ctx, tc.token)
				require.NoError(t, err)
				_, err = tr.r.Get(ctx, tc.token).Result()
				require.ErrorIs(t, err, redis.Nil)
			}
		})
	}
}
