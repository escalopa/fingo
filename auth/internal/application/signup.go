package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
)

// ---------------------- Signup ---------------------- //

type SignupParams struct {
	Name     string `validate:"required,alpha"`
	Username string `validate:"required,alphanum"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type SignupCommand interface {
	Execute(ctx context.Context, params SignupParams) error
}

type SignupCommandImpl struct {
	v  Validator
	h  PasswordHasher
	ur UserRepository
}

func (l *SignupCommandImpl) Execute(ctx context.Context, params SignupParams) error {
	if err := l.v.Validate(params); err != nil {
		return err
	}
	// Hash password
	hashedPassword, err := l.h.Hash(params.Password)
	if err != nil {
		return err
	}
	// Save user to cache
	err = l.ur.CreateUser(ctx, core.CreateUserParams{
		ID:             uuid.New(),
		Name:           params.Name,
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewSignupCommand(v Validator, h PasswordHasher, ur UserRepository) SignupCommand {
	return &SignupCommandImpl{v: v, h: h, ur: ur}
}
