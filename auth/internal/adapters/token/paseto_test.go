package token

import (
	"reflect"
	"testing"
	"time"

	ac "github.com/escalopa/gochat/auth/internal/core"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

func TestPaseto(t *testing.T) {
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u, err := randomUser()
	require.NoError(t, err)
	token, err := p.GenerateToken(u)
	require.NoError(t, err)
	require.NotEmpty(t, t)

	// VerifyToken token
	u2, err := p.VerifyToken(token)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(u, u2))
}

func TestPasetoExpired(t *testing.T) {
	// Create paseto
	p, err := NewPaseto("12345678901234567890123456789012", 1*time.Nanosecond)
	require.NoError(t, err)
	require.NotNil(t, p)
	// Generate token
	u, err := randomUser()
	require.NoError(t, err)
	token, err := p.GenerateToken(u)
	require.NoError(t, err)
	require.NotEmpty(t, t)
	// VerifyToken token
	u2, err := p.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, u2)
	// Check error
	errr, ok := err.(*errs.Error)
	require.True(t, ok)
	require.Equal(t, []string{"token expired"}, errr.Msg)
	require.Equal(t, errs.Unauthenticated, errr.Code)
}

func randomUser() (ac.User, error) {
	return ac.User{
		Email:    "test@gmail.com",
		Password: "123456",
	}, nil
}
