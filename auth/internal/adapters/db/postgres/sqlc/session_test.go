package db

import (
	"context"
	"github.com/google/uuid"
	"testing"
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
			q.CreateSession(ctx, tc.arg)
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
			q.GetSessionByID(ctx, tc.arg)
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
			q.GetUserDevices(ctx, tc.arg)
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
			q.GetUserSessions(ctx, tc.arg)
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
			q.SetSessionIsBlocked(ctx, tc.arg)
		})
	}
}
