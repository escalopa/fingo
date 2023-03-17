package core

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenPayload_MarshalBinary(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()
	payload := TokenPayload{
		UserID:    userID,
		SessionID: sessionID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	// Marshal
	b, err := payload.MarshalBinary()
	require.NoError(t, err)
	require.NotEmpty(t, b)
}
