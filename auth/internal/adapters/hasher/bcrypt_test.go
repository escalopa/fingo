package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHash(t *testing.T) {
	var err error
	var hash string
	c := NewBcryptHasher()
	// Hash the password
	hash, err = c.Hash("password")
	require.NoError(t, err)
	// Compare the hash with the password
	eq := c.Compare(hash, "password")
	require.True(t, eq)
	// Compare the hash with a different password
	eq = c.Compare(hash, "password1")
	require.False(t, eq)
	// Empty password
	_, err = c.Hash("")
	require.Error(t, err)
}
