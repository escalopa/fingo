package application

import (
	"context"
	"reflect"

	"github.com/escalopa/gofly/auth/internal/core"
	"github.com/lordvidex/errs"
)

type VerifyTokenParams struct {
	Token string `validate:"required"`
}

type VerifyTokenCommand interface {
	Execute(ctx context.Context, params VerifyTokenParams) (core.User, error)
}

type VerifyTokenCommandImpl struct {
	v  Validator
	tg TokenGenerator
}

func (v *VerifyTokenCommandImpl) Execute(ctx context.Context, params VerifyTokenParams) (core.User, error) {
	if err := v.v.Validate(params); err != nil {
		return core.User{}, err
	}
	user, err := v.tg.VerifyToken(params.Token)
	if err != nil {
		return core.User{}, err
	}
	if reflect.DeepEqual(user, core.User{}) {
		return core.User{}, errs.B().Code(errs.Unauthenticated).Msg("invalid token, token not assigned").Err()
	}
	if user.IsVerified == false {
		return core.User{}, errs.B().Code(errs.Unauthenticated).Msg("invalid token, user not verified").Err()
	}
	return user, nil
}

func NewVerifyTokenCommand(v Validator, tg TokenGenerator) VerifyTokenCommand {
	return &VerifyTokenCommandImpl{v: v, tg: tg}
}
