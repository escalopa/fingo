package hasher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHash(t *testing.T) {
	var err error
	var hash string
	ctx := context.Background()
	c := NewBcryptHasher()
	// Hash the password
	hash, err = c.Hash(ctx, "password")
	require.NoError(t, err)
	// Compare the hash with the password
	eq := c.Compare(ctx, hash, "password")
	require.True(t, eq)
	// Compare the hash with a different password
	eq = c.Compare(ctx, hash, "password1")
	require.False(t, eq)
	// Empty password
	_, err = c.Hash(ctx, "")
	require.Error(t, err)
}
