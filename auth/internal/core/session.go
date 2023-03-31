package core

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	AccessToken  string
	RefreshToken string
	UserDevice   UserDevice
	UpdatedAt    time.Time
	ExpiresAt    time.Time
}

type UserDevice struct {
	UserAgent string
	ClientIP  string
}

// ------------------------- Params -------------------------

type CreateSessionParams struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	AccessToken  string
	RefreshToken string
	UserDevice   UserDevice
}

type UpdateSessionTokenParams struct {
	ID           uuid.UUID
	AccessToken  string
	RefreshToken string
}

type SetSessionIsBlockedParams struct {
	ID        uuid.UUID
	IsBlocked bool
}

type SendNewSignInSessionParams struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	ClientIP  string `json:"client-ip"`
	UserAgent string `json:"user-agent"`
}
