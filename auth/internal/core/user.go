package core

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID
	Name       string
	Username   string
	Email      string
	Password   string
	IsVerified bool
	CreatedAt  time.Time
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

type CreateUserParams struct {
	ID             uuid.UUID
	Name           string
	Username       string
	Email          string
	HashedPassword string
}

type SetUserIsVerifiedParams struct {
	ID         uuid.UUID
	IsVerified bool
}

type ChangeUserEmailParams struct {
	ID    uuid.UUID
	Email string
}

type ChangePasswordParams struct {
	ID             uuid.UUID
	HashedPassword string
}
