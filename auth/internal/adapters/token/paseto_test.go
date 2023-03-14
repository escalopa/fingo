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
	u := randomUser()
	sID := uuid.New()
	// Generate user access token
	token, err := p.GenerateAccessToken(core.GenerateTokenParam{
		User:      u,
		SessionID: sID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// DecryptToken token
	u2, sID2, err := p.DecryptToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(u, u2))
	require.True(t, reflect.DeepEqual(sID, sID2))
}

func TestPasetoTokenizer_GenerateRefreshToken(t *testing.T) {
	t.Parallel()
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u := randomUser()
	sID := uuid.New()
	// Generate user refresh token
	token, err := p.GenerateRefreshToken(core.GenerateTokenParam{
		User:      u,
		SessionID: sID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// DecryptToken token
	u2, sID2, err := p.DecryptToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(u, u2))
	require.True(t, reflect.DeepEqual(sID, sID2))
}

func TestPasetoTokenizer_VerifyToken(t *testing.T) {
	t.Parallel()
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 0*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u := randomUser()
	// Create new user session
	sessionID := uuid.New()
	// Generate user token
	token, err := p.GenerateAccessToken(core.GenerateTokenParam{
		User:      u,
		SessionID: sessionID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// DecryptToken token
	u2, _, err := p.DecryptToken(token)
	require.Error(t, err)
	require.Empty(t, u2)
	// Check error
	er, ok := err.(*errs.Error)
	require.True(t, ok)
	require.Equal(t, []string{"token expired"}, er.Msg)
	require.Equal(t, errs.Unauthenticated, er.Code)
}

func randomUser() core.User {
	return core.User{Email: gofakeit.Email()}
}
