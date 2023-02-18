package core

import "time"

type UserToken struct {
	User      User
	IssuedAt  time.Time
	ExpiresAt time.Time
}
