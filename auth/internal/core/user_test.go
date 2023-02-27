package core

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUser_Binary(t *testing.T) {
	id, err := uuid.Parse(gofakeit.UUID())
	require.NoError(t, err)
	user := User{
		ID:         id,
		Email:      gofakeit.Email(),
		Password:   gofakeit.Password(false, false, false, false, false, 10),
		IsVerified: true,
	}
	b, err := user.MarshalBinary()
	require.NoError(t, err)
	var userB User
	err = userB.UnmarshalBinary(b)
	require.NoError(t, err)
	//require.Equal(t, user, userB)
}
