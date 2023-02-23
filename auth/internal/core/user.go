package core

import (
	"encoding/json"
)

type User struct {
	ID         string
	Name       string
	Username   string
	Email      string
	Password   string
	IsVerified bool
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}
