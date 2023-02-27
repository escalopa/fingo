package mypostgres

import (
	"context"
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
	"testing"
	"time"
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
			sr.CreateSession(ctx, tc.arg)
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
			sr.GetSessionByID(ctx, tc.arg)
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
			sr.GetUserSessions(ctx, tc.arg)
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
			sr.GetUserDevices(ctx, tc.arg)
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
			sr.SetSessionIsBlocked(ctx, tc.arg)
		})
	}
}
