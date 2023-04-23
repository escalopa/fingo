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

func TestNewTokenRepository(t *testing.T) {
	ctx := context.Background()

	tr := NewTokenRepositoryV1(testRedisClient)
	require.NotNil(t, tr)

	testCases := []struct {
		name          string
		token         string
		payload       core.TokenPayload
		storeDuration time.Duration
		sleepDuration time.Duration
		err           bool
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
			storeDuration: 2 * time.Second,
			sleepDuration: 0 * time.Second,
			err:           false,
		},
		{
			name:  "expires",
			token: gofakeit.UUID(),
			payload: core.TokenPayload{
				UserID:    uuid.New(),
				SessionID: uuid.New(),
				ClientIP:  gofakeit.IPv4Address(),
				UserAgent: gofakeit.UserAgent(),
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
			storeDuration: 1 * time.Second,
			sleepDuration: 2 * time.Second,
			err:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set token
			testRedisClient.Set(ctx, tc.token, tc.payload, tc.storeDuration)
			// get token
			time.Sleep(tc.sleepDuration)
			payload, err := tr.GetTokenPayload(ctx, tc.token)
			require.True(t, tc.err == (err != nil))
			if err == nil {
				require.NotNil(t, payload)
				require.Equal(t, tc.payload.UserID, payload.UserID)
				require.Equal(t, tc.payload.SessionID, payload.SessionID)
				require.Equal(t, tc.payload.ClientIP, payload.ClientIP)
				require.Equal(t, tc.payload.UserAgent, payload.UserAgent)
				require.WithinDurationf(t, tc.payload.IssuedAt, payload.IssuedAt, 1*time.Second, "issued at should be within 1 second")
				require.WithinDurationf(t, tc.payload.ExpiresAt, payload.ExpiresAt, 1*time.Second, "expires at should be within 1 second")
			}
		})
	}
}
