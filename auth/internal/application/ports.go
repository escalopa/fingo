package application

import (
	"context"
	ac "github.com/escalopa/gochat/auth/internal/core"
)

type UserRepository interface {
	Save(ctx context.Context, user ac.User) error
	Get(ctx context.Context, email string) (ac.User, error)
	Update(ctx context.Context, user ac.User) error
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
