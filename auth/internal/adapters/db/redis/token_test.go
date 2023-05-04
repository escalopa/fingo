package redis

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenRepository_Store(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	exp := 5 * time.Second
	tr := NewTokenRepository(testRedis, WithTokenDuration(exp))
	// Test cases
	testCases := []struct {
		name      string
		token     string
		arg       core.TokenPayload
		wantError bool
	}{
		{
			name:  "success",
			token: gofakeit.UUID(),
			arg: core.TokenPayload{
				UserID:    uuid.New(),
				SessionID: uuid.New(),
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(exp),
			},
			wantError: false,
		},
		{
			name:  "empty token",
			token: "",
			arg: core.TokenPayload{
				UserID: uuid.New(),
			},
			wantError: true,
		},
	}
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Store(ctx, tc.token, tc.arg)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if err == nil {
				_, err = tr.r.Get(ctx, tc.token).Result()
				require.NoError(t, err)
			}
		})
	}
	// Wait for token to expire
	time.Sleep(exp + 1*time.Millisecond) // Wait for token to expire
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
	exp := 1 * time.Minute
	tr := NewTokenRepository(testRedis, WithTokenDuration(exp))
	// Test cases
	testCases := []struct {
		name      string
		token     string
		arg       core.TokenPayload
		wantError bool
	}{
		{
			name:  "success",
			token: gofakeit.UUID(),
			arg: core.TokenPayload{
				UserID:    uuid.New(),
				SessionID: uuid.New(),
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(exp),
			},
			wantError: false,
		},
		{
			name:      "empty token",
			token:     "",
			arg:       core.TokenPayload{},
			wantError: true,
		},
	}
	// Run test cases again to check if the token is deleted
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.Store(ctx, tc.token, tc.arg)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
