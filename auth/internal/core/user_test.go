package core

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUser_Binary(t *testing.T) {
	user := User{
		ID:         gofakeit.UUID(),
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
