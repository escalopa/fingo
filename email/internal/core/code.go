package core

import (
	"encoding/json"
	"time"
)

type VerificationCode struct {
	Code   string
	SentAt time.Time
}

func (vc VerificationCode) MarshalBinary() ([]byte, error) {
	return json.Marshal(vc)
}

func (vc VerificationCode) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &vc)
}
