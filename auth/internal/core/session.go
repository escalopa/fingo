package core

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type UserDevice struct {
	UserAgent string
	ClientIP  string
}

type CreateSessionParams struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	UserAgent    string
	ClientIp     string
}

type SetSessionIsBlockedParams struct {
	ID        uuid.UUID
	IsBlocked bool
}
