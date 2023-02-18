package application

import (
	"github.com/lordvidex/errs"
)

// ---------------------- Signin ---------------------- //

type SigninParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type SigninCommand interface {
	Execute(params SigninParams) (string, error)
}

type SigninCommandImpl struct {
	v  Validator
	h  PasswordHasher
	ur UserRepository
	tg TokenGenerator
}

func (s *SigninCommandImpl) Execute(params SigninParams) (string, error) {
	if err := s.v.Validate(params); err != nil {
		return "", err
	}
	// Get user from cache
	user, err := s.ur.Get(params.Email)
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
