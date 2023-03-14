package core

import (
	"encoding/json"
	"github.com/google/uuid"
)

type TokenCache struct {
	UserID    string
	ClientIP  string
	UserAgent string
	Roles     []string
}

func (t TokenCache) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t TokenCache) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &t)
}

// ------------------------- Params -------------------------

type GenerateTokenParam struct {
	User      User
	SessionID uuid.UUID
}
