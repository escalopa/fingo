package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestUser_CreateUser(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  CreateUserParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.CreateUser(ctx, tc.arg)
		})
	}
}

func TestUser_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  string // emailarg

	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = q.GetUserByEmail(ctx, tc.arg)
		})
	}
}

func TestUser_GetUserByID(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  uuid.UUID // userID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = q.GetUserByID(ctx, tc.arg)
		})
	}
}

func TestUser_GetUserByUsername(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  string // username
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = q.GetUserByUsername(ctx, tc.arg)
		})
	}
}

func TestUser_SetUserIsVerified(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  SetUserIsVerifiedParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.SetUserIsVerified(ctx, tc.arg)
		})
	}
}

func TestUser_ChangeNames(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  ChangeNamesParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.ChangeNames(ctx, tc.arg)
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  ChangePasswordParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.ChangePassword(ctx, tc.arg)
		})
	}
}

func TestUser_ChangeUserEmail(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  ChangeUserEmailParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.ChangeUserEmail(ctx, tc.arg)
		})
	}
}

func TestUser_DeleteUserByID(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  uuid.UUID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.DeleteUserByID(ctx, tc.arg)
		})
	}
}
