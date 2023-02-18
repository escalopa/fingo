package redis

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/gofly/auth/internal/core"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSaveUser(t *testing.T) {
	ur := NewUserRepository(testRedis, WithUserContext(testContext))

	genID := func() string {
		id, err := newUserID()
		require.NoError(t, err)
		return id
	}

	testCases := []struct {
		name string
		user core.User
		err  error
	}{
		{
			name: "save user successfully",
			user: core.User{
				ID:         genID(),
				Email:      gofakeit.Email(),
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		}, {
			name: "save user a user to prepare for duplicate user test",
			user: core.User{
				ID:         genID(),
				Email:      "ahmad@gmail.com",
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		}, {
			name: "save duplicate user with same email",
			user: core.User{
				ID:         genID(),
				Email:      "ahmad@gmail.com",
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: errs.B().Code(errs.AlreadyExists).Msg("user already exists").Err(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ur.Save(tc.user)
			compareErrors(t, err, tc.err)
		})
	}
}

func TestGetUser(t *testing.T) {
	ur := NewUserRepository(testRedis, WithUserContext(testContext))

	testCases := []struct {
		name string
		user core.User
		err  error
	}{
		{
			name: "get user successfully",
			user: core.User{
				ID:         "1",
				Email:      gofakeit.Email(),
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ur.Save(tc.user)
			require.NoError(t, err)
			_, err = ur.Get(tc.user.Email)
			compareErrors(t, err, tc.err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ur := NewUserRepository(testRedis, WithUserContext(testContext))

	testCases := []struct {
		name string
		user core.User
		err  error
	}{
		{
			name: "update user successfully",
			user: core.User{
				ID:         "1",
				Email:      gofakeit.Email(),
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// save user
			err := ur.Save(tc.user)
			require.NoError(t, err)
			// get user
			u1, err := ur.Get(tc.user.Email)
			require.NoError(t, err)
			// update user
			tc.user.IsVerified = true
			err = ur.Update(tc.user)
			require.NoError(t, err)
			// get user
			u2, err := ur.Get(tc.user.Email)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(u1, u2),
				"users are not equal actual:%s, expected:%s", u1, u2)
		})
	}
}
