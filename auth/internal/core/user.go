package core

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Phone           string
	Username        string
	Gender          string
	Email           string
	HashedPassword  string
	IsEmailVerified bool
	IsPhoneVerified bool
	CreatedAt       time.Time
}

// ------------------------- Params -------------------------

type CreateUserParams struct {
	ID             uuid.UUID
	FirstName      string
	LastName       string
	Phone          string
	Username       string
	Gender         string
	Email          string
	BirthDate      time.Time
	HashedPassword string
}
