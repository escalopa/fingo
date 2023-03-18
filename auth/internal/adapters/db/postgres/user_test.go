package mypostgres

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur, err := NewUserRepository(testPGConn)
	require.NoError(t, err)
	// Create user
	user := randomUser()
	// Create test cases
	testCases := []struct {
		name      string
		arg       core.CreateUserParams
		wantError bool
	}{
		{
			name:      "valid user",
			arg:       user,
			wantError: false,
		},
		{
			name:      "duplicate user",
			arg:       user,
			wantError: true,
		},
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = ur.CreateUser(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				u, err := ur.GetUserByEmail(ctx, tc.arg.Email)
				if err != nil {
					t.Errorf("unexcpected error, got %s", err)
				}
				require.Equal(t, tc.arg.Email, u.Email)
				require.Equal(t, tc.arg.FirstName, u.FirstName)
				require.Equal(t, tc.arg.LastName, u.LastName)
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur, err := NewUserRepository(testPGConn)
	require.NoError(t, err)
	// Create user
	user := randomUser()
	// Create test cases
	testCases := []struct {
		name      string
		arg       uuid.UUID // userID
		wantError bool
	}{
		{
			name: "success",
			arg:  user.ID,
		},
		{
			name:      "not found",
			arg:       uuid.New(),
			wantError: true,
		},
	}
	// Create user
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexcpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ur.GetUserByID(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				require.Equal(t, tc.arg, u.ID)
			}
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur, err := NewUserRepository(testPGConn)
	require.NoError(t, err)
	// Create user
	user := randomUser()
	// Create test cases
	testCases := []struct {
		name      string
		arg       string // email
		wantError bool
	}{
		{
			name:      "valid email",
			arg:       user.Email,
			wantError: false,
		},
		{
			name:      "not found",
			arg:       gofakeit.Email(),
			wantError: true,
		},
	}
	// Create user
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexcpected error, got %s", err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ur.GetUserByEmail(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				require.Equal(t, tc.arg, u.Email)
			}
		})
	}
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur, err := NewUserRepository(testPGConn)
	require.NoError(t, err)
	// Create user
	user := randomUser()
	// Create test cases
	testCases := []struct {
		name        string
		arg         uuid.UUID
		expectError bool
	}{
		{
			name:        "success",
			arg:         user.ID,
			expectError: false,
		},
		{
			name:        "not found",
			arg:         uuid.New(),
			expectError: true,
		},
	}
	// Create user
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexcpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = ur.DeleteUserByID(ctx, tc.arg)
			if err != nil && !tc.expectError {
				t.Errorf("unexcpected error, got %s", err)
			}
			if err == nil && tc.expectError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				u, err := ur.GetUserByID(ctx, tc.arg)
				require.Error(t, err)
				require.Empty(t, u)
			}
		})
	}
}

func randomUser() core.CreateUserParams {
	return core.CreateUserParams{
		ID:        uuid.New(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Phone:     gofakeit.Phone(),
		Username:  gofakeit.Username(),
		Gender:    strings.ToUpper(gofakeit.Gender()),
		Email:     gofakeit.Email(),
		BirthDate: gofakeit.Date(),
	}
}
