package application

import (
	ac "github.com/escalopa/gofly/auth/internal/core"
)

type UserRepository interface {
	Save(user ac.User) error
	Get(email string) (ac.User, error)
	Update(user ac.User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type TokenGenerator interface {
	GenerateToken(user ac.User) (string, error)
	VerifyToken(token string) (ac.User, error)
}

type Validator interface {
	Validate(s interface{}) error
}
