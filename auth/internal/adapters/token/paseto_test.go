package token

import (
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

func TestNewPaseto(t *testing.T) {
	t.Parallel()
	_, err := NewPaseto("", 1*time.Minute, 2*time.Minute)
	require.Error(t, err)
	er, ok := err.(*errs.Error)
	require.True(t, ok)
	require.Equal(t, er.Code, errs.InvalidArgument)
}
func TestPasetoTokenizer_GenerateAccessToken(t *testing.T) {
	t.Parallel()
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	user := randomUser()
	sessionID := uuid.New()
	// Generate user access token
	token, err := p.GenerateAccessToken(core.GenerateTokenParam{
		UserID:    user.ID,
		SessionID: sessionID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// DecryptToken token
	payload, err := p.DecryptToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(user.ID, payload.UserID))
	require.True(t, reflect.DeepEqual(sessionID, payload.SessionID))
}

func TestPasetoTokenizer_GenerateRefreshToken(t *testing.T) {
	t.Parallel()
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	user := randomUser()
	sessionID := uuid.New()
	// Generate user refresh token
	token, err := p.GenerateRefreshToken(core.GenerateTokenParam{
		UserID:    user.ID,
		SessionID: sessionID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// DecryptToken token
	payload, err := p.DecryptToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(user.ID, payload.UserID))
	require.True(t, reflect.DeepEqual(sessionID, payload.SessionID))
}

func randomUser() core.User {
	return core.User{Email: gofakeit.Email()}
}
