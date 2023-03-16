package mypostgres

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRolesRepository_CreateRole(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	rr := NewRolesRepository(testPGConn)
	testCases := []struct {
		name      string
		role      string
		wantError bool
	}{
		{
			name:      "create role",
			role:      "test_role",
			wantError: false,
		},
		{
			name:      "create role with empty name",
			role:      "",
			wantError: true,
		},
		{
			name:      "create role with duplicate name",
			role:      "test_role",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := rr.CreateRole(ctx, tc.role)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestRolesRepository_GetRoleByName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	rr := NewRolesRepository(testPGConn)
	// Create user roles
	roleNames := []string{gofakeit.Name(), gofakeit.Name()}
	// Create user for test
	user := randomUser()
	err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	testCases := []struct {
		name      string
		role      string
		wantError bool
	}{
		{
			name:      "get role 1",
			role:      roleNames[0],
			wantError: false,
		},
		{
			name:      "get role 2",
			role:      roleNames[1],
			wantError: false,
		},
		{
			name:      "get role with empty name",
			role:      "",
			wantError: true,
		},
	}
	// Create roles
	for _, rn := range roleNames {
		err = rr.CreateRole(ctx, rn)
		require.NoError(t, err)
	}
	// Run add user to roles test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := rr.GetRoleByName(ctx, tc.role)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if err != nil && r.Name != tc.role {
				t.Errorf("expected equal role names, expected: %s, got:%s", tc.name, r.Name)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestRolesRepository_GetUserRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	rr := NewRolesRepository(testPGConn)
	// Create user roles
	roleNames := []string{gofakeit.Name(), gofakeit.Name(), gofakeit.Name()}
	sort.Strings(roleNames)
	// Create user for test
	user := randomUser()
	err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	testCases := []struct {
		name      string
		userID    uuid.UUID
		roles     []string
		wantError bool
	}{
		{
			name:      "get user roles",
			userID:    user.ID,
			roles:     roleNames,
			wantError: false,
		},
	}
	// Create roles
	for _, rn := range roleNames {
		err = rr.CreateRole(ctx, rn)
		require.NoError(t, err)
	}
	// Add user to roles
	for _, rn := range roleNames {
		err = rr.GrantRole(ctx, core.GrantRoleToUserParams{
			RoleName: rn,
			UserID:   user.ID,
		})
		require.NoError(t, err)
	}
	// Run add user to roles test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			roles, err := rr.GetUserRoles(ctx, tc.userID)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err != nil {
				if len(roles) != len(tc.roles) {
					t.Errorf("expected equal roles length, expected: %d, got:%d", len(tc.roles), len(roles))
				}
				// Sort roles by name
				sort.Strings(tc.roles)
				sort.Slice(roles, func(i, j int) bool { return strings.Compare(roles[i].Name, roles[j].Name) < 0 })
				for i, r := range roles {
					if r.Name != tc.roles[i] {
						t.Errorf("expected equal role names, expected: %s, got:%s", tc.roles[i], r.Name)
					}
				}
			}
		})
	}
}

func TestRolesRepository_GrantRole(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	rr := NewRolesRepository(testPGConn)
	// Create user roles
	roleNames := []string{gofakeit.Name(), gofakeit.Name()}
	// Create user for test
	user := randomUser()
	err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	testCases := []struct {
		name      string
		userID    uuid.UUID
		role      string
		wantError bool
	}{
		{
			name:      "grant role",
			userID:    user.ID,
			role:      roleNames[0],
			wantError: false,
		},
		{
			name:      "grant role twice",
			userID:    user.ID,
			role:      roleNames[0],
			wantError: true,
		},
		{
			name:      "grant role with empty name",
			userID:    user.ID,
			role:      "",
			wantError: true,
		},
		{
			name:      "grant role with empty user id",
			userID:    uuid.Nil,
			role:      roleNames[1],
			wantError: true,
		},
		{
			name:      "grant role with duplicate name",
			userID:    user.ID,
			role:      roleNames[0],
			wantError: true,
		},
	}
	// Create roles
	for _, rn := range roleNames {
		err = rr.CreateRole(ctx, rn)
		require.NoError(t, err)
	}
	// Run add user to roles test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = rr.GrantRole(ctx, core.GrantRoleToUserParams{
				RoleName: tc.role,
				UserID:   tc.userID,
			})
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestRolesRepository_RevokeRole(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	rr := NewRolesRepository(testPGConn)
	// Create user roles
	roleNames := []string{gofakeit.Name(), gofakeit.Name()}
	// Create user for test
	user := randomUser()
	err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	testCases := []struct {
		name      string
		userID    uuid.UUID
		role      string
		wantError bool
	}{
		{
			name:      "revoke role",
			userID:    user.ID,
			role:      roleNames[0],
			wantError: false,
		},
		{
			name:      "revoke role with empty name",
			userID:    user.ID,
			role:      "",
			wantError: true,
		},
		{
			name:      "revoke role with empty user id",
			userID:    uuid.Nil,
			role:      roleNames[1],
			wantError: true,
		},
		{
			name:      "revoke role twice",
			userID:    user.ID,
			role:      roleNames[0],
			wantError: true,
		},
	}
	// Create roles
	for _, rn := range roleNames {
		err = rr.CreateRole(ctx, rn)
		require.NoError(t, err)
	}
	// Add user to roles
	for _, rn := range roleNames {
		err = rr.GrantRole(ctx, core.GrantRoleToUserParams{
			RoleName: rn,
			UserID:   user.ID,
		})
		require.NoError(t, err)
	}
	// Run add user to roles test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = rr.RevokeRole(ctx, core.RevokeRoleFromUserParams{
				RoleName: tc.role,
				UserID:   tc.userID,
			})
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
