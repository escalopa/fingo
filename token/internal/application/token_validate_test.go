package application

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/pkg/pkgCore"
	"github.com/escalopa/fingo/token/internal/adapters/cache"
	"github.com/escalopa/fingo/token/internal/adapters/validator"
	"github.com/escalopa/fingo/token/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewTokenValidateCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := validator.NewValidator()
	tr, err := cache.NewTokenRepositoryV1(testRedisClient)
	require.NoError(t, err)
	c := NewTokenValidateCommand(v, tr)

	testCases := []struct {
		name      string
		token     string
		payload   *core.TokenPayload
		wait      time.Duration
		wantError bool
	}{
		{
			name:  "success",
			token: gofakeit.UUID(),
			payload: &core.TokenPayload{
				UserID:    uuid.New(),
				SessionID: uuid.New(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
			wantError: false,
		},
		{
			name:  "expired",
			token: gofakeit.UUID(),
			payload: &core.TokenPayload{
				UserID:    uuid.New(),
				SessionID: uuid.New(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
				IssuedAt:  time.Now().Add(-(2 * time.Hour)),
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Set token in cache
			testRedisClient.Set(ctx, tc.token, tc.payload, 0)
			require.NoError(t, err)
			// Set client IP & user agent in context
			ctx = context.WithValue(ctx, pkgCore.ContextKeyClientIP, tc.payload.ClientIP)
			ctx = context.WithValue(ctx, pkgCore.ContextKeyUserAgent, tc.payload.UserAgent)
			// Execute command
			err = c.Execute(ctx, TokenValidateParams{AccessToken: tc.token})
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
