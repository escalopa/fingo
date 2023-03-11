package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestSession_CreateSession(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  CreateSessionParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.CreateSession(ctx, tc.arg)
		})
	}
}

func TestSession_GetSessionByID(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  uuid.UUID // sessionID
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = q.GetSessionByID(ctx, tc.arg)
		})
	}
}

func TestSession_GetUserDevices(t *testing.T) {
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
			_, _ = q.GetUserDevices(ctx, tc.arg)
		})
	}
}

func TestSession_GetUserSessions(t *testing.T) {
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
			_, _ = q.GetUserSessions(ctx, tc.arg)
		})
	}
}

func TestSession_SetSessionIsBlocked(t *testing.T) {
	ctx := context.Background()
	q := New(dbSQL)

	testCases := []struct {
		name string
		arg  SetSessionIsBlockedParams
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = q.SetSessionIsBlocked(ctx, tc.arg)
		})
	}
}
