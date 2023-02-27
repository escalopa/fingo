package redis

import (
	"context"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

func TestSaveUser(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(testRedis)

	// Test cases
	testCases := []struct {
		name string
		user core.User
		err  error
	}{
		{
			name: "save user successfully",
			user: core.User{
				ID:         randomUserID(t),
				Email:      gofakeit.Email(),
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		}, {
			name: "save user a user to prepare for duplicate user test",
			user: core.User{
				ID:         randomUserID(t),
				Email:      "ahmad@gmail.com",
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		}, {
			name: "save duplicate user with same email",
			user: core.User{
				ID:         randomUserID(t),
				Email:      "ahmad@gmail.com",
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: errs.B().Code(errs.AlreadyExists).Msg("user already exists").Err(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ur.Save(ctx, tc.user)
			compareErrors(t, err, tc.err)
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(testRedis)
	testCases := []struct {
		name string
		user core.User
		err  error
	}{
		{
			name: "get user successfully",
			user: core.User{
				ID:         randomUserID(t),
				Email:      gofakeit.Email(),
				Password:   gofakeit.Password(true, true, true, true, true, 32),
				IsVerified: false,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ur.Save(ctx, tc.user)
			require.NoError(t, err)
			_, err = ur.Get(ctx, tc.user.Email)
			compareErrors(t, err, tc.err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(testRedis)
	testCases := []struct {
		name string
		user core.User
		err  error
	}{
		{
			name: "update user successfully",
			user: core.User{
				ID:         randomUserID(t),
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
			err := ur.Save(ctx, tc.user)
			require.NoError(t, err)
			// get user
			u1, err := ur.Get(ctx, tc.user.Email)
			require.NoError(t, err)
			// update user
			u1.IsVerified = true
			err = ur.Update(ctx, u1)
			require.NoError(t, err)
			// get user
			u2, err := ur.Get(ctx, tc.user.Email)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(u1, u2),
				"users are not equal actual:%s, expected:%s", u1, u2)
		})
	}
}

func randomUserID(t *testing.T) uuid.UUID {
	id, err := newUserID()
	require.NoError(t, err)
	return id
}
