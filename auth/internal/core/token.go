package core

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type TokenPayload struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	ClientIP  string
	UserAgent string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Roles     []string
}

func (t TokenPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t TokenPayload) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &t)
}

// ------------------------- Params -------------------------

type GenerateTokenParam struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	ClientIP  string
	UserAgent string
	Roles     []string
}
