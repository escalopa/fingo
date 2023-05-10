package core

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

	// Unmarshal
	var payload2 TokenPayload
	err = payload2.UnmarshalBinary(b)
	require.NoError(t, err)
}
