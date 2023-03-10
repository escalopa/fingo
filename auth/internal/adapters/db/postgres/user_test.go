package mypostgres

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"testing"
)

func TestUserRepository_CreateUser(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  core.CreateUserParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.CreateUser(ctx, tc.arg)
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  string // emailarg

	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.GetUserByEmail(ctx, tc.arg)
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  uuid.UUID // userID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.GetUserByID(ctx, tc.arg)
		})
	}
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  string // username
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.GetUserByUsername(ctx, tc.arg)
		})
	}
}

func TestUserRepository_SetUserIsVerified(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  core.SetUserIsVerifiedParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.SetUserIsVerified(ctx, tc.arg)
		})
	}
}

func TestUserRepository_ChangeNames(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  core.ChangeNamesParam
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.ChangeNames(ctx, tc.arg)
		})
	}
}

func TestUserRepository_ChangePassword(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  core.ChangePasswordParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.ChangePassword(ctx, tc.arg)
		})
	}
}

func TestUserRepository_ChangeUserEmail(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  core.ChangeUserEmailParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.ChangeUserEmail(ctx, tc.arg)
		})
	}
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	ctx := context.Background()
	ur := NewUserRepository(dbSQL)

	testCases := []struct {
		name string
		arg  uuid.UUID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur.DeleteUserByID(ctx, tc.arg)
		})
	}
}
