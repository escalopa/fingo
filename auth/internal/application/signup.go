package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"time"
)

// ---------------------- Signup ---------------------- //

type SignupParams struct {
	FirstName string `validate:"required,alpha"`
	LastName  string `validate:"required,alpha"`
	Username  string `validate:"required,alphanum"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
}

type SignupCommand interface {
	Execute(ctx context.Context, params SignupParams) error
}

type SignupCommandImpl struct {
	v  Validator
	h  PasswordHasher
	ur UserRepository
}

func (c *SignupCommandImpl) Execute(ctx context.Context, params SignupParams) error {
	return executeWithContextTimeout(ctx, 10*time.Second, func() error {
		if err := c.v.Validate(params); err != nil {
			return err
		}
		// Hash password
		hashedPassword, err := c.h.Hash(params.Password)
		if err != nil {
			return err
		}
		// Save user to db
		err = c.ur.CreateUser(ctx, core.CreateUserParams{
			ID:             uuid.New(),
			FirstName:      params.FirstName,
			LastName:       params.LastName,
			Username:       params.Username,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func NewSignupCommand(v Validator, h PasswordHasher, ur UserRepository) SignupCommand {
	return &SignupCommandImpl{v: v, h: h, ur: ur}
}
