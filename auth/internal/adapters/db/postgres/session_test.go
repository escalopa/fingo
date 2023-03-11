package mypostgres

import (
	"context"
	"testing"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	ctx := context.Background()
	sr := NewSessionRepository(dbSQL, 24*time.Hour)

	testCases := []struct {
		name string
		arg  core.CreateSessionParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = sr.CreateSession(ctx, tc.arg)
		})
	}
}

func TestSessionRepository_GetSessionByID(t *testing.T) {
	ctx := context.Background()
	sr := NewSessionRepository(dbSQL, 24*time.Hour)

	testCases := []struct {
		name string
		arg  uuid.UUID // sessionID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = sr.GetSessionByID(ctx, tc.arg)
		})
	}
}

func TestSessionRepository_GetUserSessions(t *testing.T) {
	ctx := context.Background()
	sr := NewSessionRepository(dbSQL, 24*time.Hour)

	testCases := []struct {
		name string
		arg  uuid.UUID // userID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = sr.GetUserSessions(ctx, tc.arg)
		})
	}
}

func TestSessionRepository_GetUserDevices(t *testing.T) {
	ctx := context.Background()
	sr := NewSessionRepository(dbSQL, 24*time.Hour)

	testCases := []struct {
		name string
		arg  uuid.UUID // userID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = sr.GetUserDevices(ctx, tc.arg)
		})
	}
}

func TestSessionRepository_SetSessionIsBlocked(t *testing.T) {
	ctx := context.Background()
	sr := NewSessionRepository(dbSQL, 24*time.Hour)

	testCases := []struct {
		name string
		arg  core.SetSessionIsBlockedParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = sr.SetSessionIsBlocked(ctx, tc.arg)
		})
	}
}
