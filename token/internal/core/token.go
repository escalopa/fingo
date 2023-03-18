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
}

func (t TokenPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t TokenPayload) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &t)
}
