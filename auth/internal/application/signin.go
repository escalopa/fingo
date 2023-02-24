package application

import (
	"context"
	"github.com/lordvidex/errs"
)

// ---------------------- Signin ---------------------- //

type SigninParams struct {
	Email    string            `validate:"required,email"`
	Password string            `validate:"required,min=8"`
	MetaData map[string]string `validate:"required"`
}

type SigninCommand interface {
	Execute(ctx context.Context, params SigninParams) (string, error)
}

type SigninCommandImpl struct {
	v  Validator
	h  PasswordHasher
	ur UserRepository
	tg TokenGenerator
}

func (s *SigninCommandImpl) Execute(ctx context.Context, params SigninParams) (string, error) {
	if err := s.v.Validate(params); err != nil {
		return "", err
	}
	// Get user from cache
	user, err := s.ur.Get(ctx, params.Email)
	if err != nil {
		return "", err
	}
	// Compare password
	if !s.h.Compare(user.Password, params.Password) {
		return "", errs.B().Code(errs.InvalidArgument).Msg("password is incorrect").Err()
	}
	// if user is not verified, return error
	if !user.IsVerified {
		return "", errs.B().Code(errs.Unauthenticated).Msg("user is not verified").Err()
	}
	// Generate user token
	token, err := s.tg.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewSigninCommand(v Validator, h PasswordHasher, tg TokenGenerator, ur UserRepository) SigninCommand {
	return &SigninCommandImpl{v: v, h: h, tg: tg, ur: ur}
}
