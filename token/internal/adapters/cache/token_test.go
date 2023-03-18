package cache

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/token/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewTokenRepositoryV1(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tr, err := NewTokenRepositoryV1(testRedisClient)
	require.NoError(t, err)
	require.NotNil(t, tr)

	testCases := []struct {
		name    string
		token   string
		payload core.TokenPayload
	}{
		{
			name:  "success",
			token: gofakeit.UUID(),
			payload: core.TokenPayload{
				UserID:    uuid.New(),
				SessionID: uuid.New(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// set token
			tc := tc
			testRedisClient.Set(ctx, tc.token, tc.payload, 2*time.Second)
			require.NoError(t, err)
			// get token
			var payload *core.TokenPayload
			payload, err = tr.GetTokenPayload(ctx, tc.token)
			require.NoError(t, err)
			require.NotNil(t, payload)
			require.Equal(t, tc.payload.UserID, payload.UserID)
			require.Equal(t, tc.payload.SessionID, payload.SessionID)
			require.Equal(t, tc.payload.ClientIP, payload.ClientIP)
			require.Equal(t, tc.payload.UserAgent, payload.UserAgent)
			require.WithinDurationf(t, tc.payload.IssuedAt, payload.IssuedAt, 1*time.Second, "issued at should be within 1 second")
			require.WithinDurationf(t, tc.payload.ExpiresAt, payload.ExpiresAt, 1*time.Second, "expires at should be within 1 second")
		})
	}
}
