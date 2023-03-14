package mypostgres

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	sr, err := NewSessionRepository(testPGConn, WithSessionDuration(1*time.Hour))
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create user
	user := randomUser()
	// Create test cases
	testCases := []struct {
		name      string
		arg       core.CreateSessionParams
		wantError bool
	}{
		{
			name:      "success",
			arg:       randomSession(user.ID),
			wantError: false,
		},
		{
			name:      "invalid params",
			arg:       core.CreateSessionParams{},
			wantError: true,
		},
	}
	// Create user to bind sessions to
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = sr.CreateSession(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestSessionRepository_GetSessionByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	sr, err := NewSessionRepository(testPGConn, WithSessionDuration(1*time.Hour))
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create user
	user := randomUser()
	session := randomSession(user.ID)
	// Create test cases
	testCases := []struct {
		name      string
		arg       uuid.UUID
		wantError bool
	}{
		{
			name:      "success",
			arg:       session.ID,
			wantError: false,
		},
		{
			name:      "invalid sessions id",
			arg:       uuid.New(),
			wantError: true,
		},
	}
	// Create user to bind sessions to
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create sessions
	err = sr.CreateSession(ctx, session)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := sr.GetSessionByID(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				require.Equal(t, tc.arg, s.ID)
			}
		})
	}
}

func TestSessionRepository_GetUserSessions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	sr, err := NewSessionRepository(testPGConn, WithSessionDuration(1*time.Hour))
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create user
	user := randomUser()
	userSessions := []core.CreateSessionParams{
		randomSession(user.ID),
		randomSession(user.ID),
		randomSession(user.ID),
	}
	// Create user to bind userSessions to
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create test cases
	testCases := []struct {
		name      string
		arg       uuid.UUID // userID
		sessions  []core.CreateSessionParams
		wantError bool
	}{
		{
			name:      "success",
			arg:       user.ID,
			sessions:  userSessions,
			wantError: false,
		},
		{
			name:      "invalid user id",
			arg:       uuid.New(),
			sessions:  []core.CreateSessionParams{},
			wantError: false,
		},
	}
	// Create userSessions
	for _, s := range userSessions {
		err = sr.CreateSession(ctx, s)
		if err != nil {
			t.Errorf("unexpected error, got %s", err)
		}
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := sr.GetUserSessions(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			require.Len(t, s, len(tc.sessions))
		})
	}
}

func TestSessionRepository_SetSessionIsBlocked(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	sr, err := NewSessionRepository(testPGConn, WithSessionDuration(1*time.Hour))
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create user
	user := randomUser()
	session := randomSession(user.ID)
	// Create test cases
	testCases := []struct {
		name      string
		arg       core.SetSessionIsBlockedParams
		wantError bool
	}{
		{
			name: "valid sessions id",
			arg: core.SetSessionIsBlockedParams{
				ID:        session.ID,
				IsBlocked: true,
			},
			wantError: false,
		},
		{
			name: "invalid sessions id",
			arg: core.SetSessionIsBlockedParams{
				ID:        uuid.New(),
				IsBlocked: true,
			},
			wantError: true,
		},
	}
	// Create user to bind sessions to
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create sessions
	err = sr.CreateSession(ctx, session)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = sr.SetSessionIsBlocked(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				s, err := sr.GetSessionByID(ctx, tc.arg.ID)
				require.NoError(t, err)
				require.Equal(t, tc.arg.IsBlocked, s.IsBlocked)
			}
		})
	}
}

func TestSessionRepository_UpdateSessionRefreshToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	sr, err := NewSessionRepository(testPGConn, WithSessionDuration(1*time.Hour))
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create user
	user := randomUser()
	session := randomSession(user.ID)
	// Create test cases
	testCases := []struct {
		name      string
		arg       core.UpdateSessionRefreshTokenParams
		wantError bool
	}{
		{
			name: "success",
			arg: core.UpdateSessionRefreshTokenParams{
				ID:           session.ID,
				RefreshToken: "new-refresh-token",
			},
			wantError: false,
		},
		{
			name: "invalid sessions id",
			arg: core.UpdateSessionRefreshTokenParams{
				ID:           uuid.New(),
				RefreshToken: "new-refresh-token",
			},
			wantError: true,
		},
	}
	// Create user to bind sessions to
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create sessions
	err = sr.CreateSession(ctx, session)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = sr.UpdateSessionRefreshToken(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)

			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				s, err := sr.GetSessionByID(ctx, tc.arg.ID)
				require.NoError(t, err)
				require.Equal(t, tc.arg.RefreshToken, s.RefreshToken)
			}
		})
	}
}

func TestSessionRepository_DeleteSessionByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ur := NewUserRepository(testPGConn)
	sr, err := NewSessionRepository(testPGConn, WithSessionDuration(1*time.Hour))
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create user
	user := randomUser()
	session := randomSession(user.ID)
	testCases := []struct {
		name      string
		arg       uuid.UUID // sessionID
		wantError bool
	}{
		{
			name:      "success",
			arg:       session.ID,
			wantError: false,
		},
		{
			name:      "invalid sessions id",
			arg:       uuid.New(),
			wantError: true,
		},
	}
	// Create user to bind sessions to
	err = ur.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Create sessions
	err = sr.CreateSession(ctx, session)
	if err != nil {
		t.Errorf("unexpected error, got %s", err)
	}
	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = sr.DeleteSessionByID(ctx, tc.arg)
			if err != nil && !tc.wantError {
				t.Errorf("unexpected error, got %s", err)
			}
			if err == nil && tc.wantError {
				t.Errorf("expected error, got nil")
			}
			if err == nil {
				s, err := sr.GetSessionByID(ctx, tc.arg)
				require.Error(t, err)
				require.Empty(t, s)
			}
		})
	}
}

func randomSession(userID uuid.UUID) core.CreateSessionParams {
	return core.CreateSessionParams{
		ID:           uuid.New(),
		UserID:       userID,
		AccessToken:  gofakeit.UUID(),
		RefreshToken: gofakeit.UUID(),
		UserDevice: core.UserDevice{
			ClientIP:  gofakeit.IPv4Address(),
			UserAgent: gofakeit.UserAgent(),
		},
	}
}
