package core

import (
	"github.com/google/uuid"
	"time"
)

type UserToken struct {
	User      User
	SessionID uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
}
