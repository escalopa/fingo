package token

import (
	"reflect"
	"testing"
	"time"

	ac "github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

func TestNewPaseto(t *testing.T) {
	_, err := NewPaseto("", 1*time.Minute, 2*time.Minute)
	require.Error(t, err)
	er, ok := err.(*errs.Error)
	require.True(t, ok)
	require.Equal(t, er.Code, errs.InvalidArgument)
}
func TestPasetoTokenizer_GenerateAccessToken(t *testing.T) {
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u, err := randomUser()
	require.NoError(t, err)
	sID := uuid.New()
	// Generate user access token
	token, err := p.GenerateAccessToken(GenerateTokenParam{
		User:      u,
		SessionID: sID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// VerifyToken token
	u2, sID2, err := p.VerifyToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(u, u2))
	require.True(t, reflect.DeepEqual(sID, sID2))
}

func TestPasetoTokenizer_GenerateRefreshToken(t *testing.T) {
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u, err := randomUser()
	require.NoError(t, err)
	sID := uuid.New()
	// Generate user refresh token
	token, err := p.GenerateRefreshToken(GenerateTokenParam{
		User:      u,
		SessionID: sID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// VerifyToken token
	u2, sID2, err := p.VerifyToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(u, u2))
	require.True(t, reflect.DeepEqual(sID, sID2))
}

func TestPasetoTokenizer_VerifyToken(t *testing.T) {
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 0*time.Minute, 2*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u, err := randomUser()
	require.NoError(t, err)
	// Create new user session
	sessionID := uuid.New()
	// Generate user token
	token, err := p.GenerateAccessToken(GenerateTokenParam{
		User:      u,
		SessionID: sessionID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// VerifyToken token
	u2, _, err := p.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, u2)
	// Check error
	er, ok := err.(*errs.Error)
	require.True(t, ok)
	require.Equal(t, []string{"token expired"}, er.Msg)
	require.Equal(t, errs.Unauthenticated, er.Code)
}

func randomUser() (ac.User, error) {
	return ac.User{
		Email:    "test@gmail.com",
		Password: "123456",
	}, nil
}
