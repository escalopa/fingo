package application

import (
	"context"
	ac "github.com/escalopa/gochat/auth/internal/core"
)

// ---------------------- Signup ---------------------- //

type SignupParams struct {
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
	// Create user
	user := ac.User{
		Email:      params.Email,
		Password:   hashedPassword,
		IsVerified: false,
	}
	// Save user to cache
	err = l.ur.Save(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func NewSignupCommand(v Validator, h PasswordHasher, ur UserRepository) SignupCommand {
	return &SignupCommandImpl{v: v, h: h, ur: ur}
}
