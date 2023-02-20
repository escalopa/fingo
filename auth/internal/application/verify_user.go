package application

import (
	"context"

	"github.com/escalopa/gofly/pb"
	"github.com/lordvidex/errs"
)

// ---------------------- Verify User ---------------------- //

type VerifyUserParam struct {
	Email string `validate:"required,email"`
	Code  string `validate:"required,numeric"`
}

type VerifyUserCommand interface {
	Execute(ctx context.Context, params VerifyUserParam) error
}

type VerifyUserCommandImpl struct {
	v   Validator
	ur  UserRepository
	esc pb.EmailServiceClient
}

func (vu *VerifyUserCommandImpl) Execute(ctx context.Context, params VerifyUserParam) error {
	if err := vu.v.Validate(params); err != nil {
		return err
	}
	user, err := vu.ur.Get(ctx, params.Email)
	if err != nil {
		return err
	}
	if user.IsVerified {
		return errs.B().Code(errs.InvalidArgument).Msg("user is already verified").Err()
	}
	// VerifyToken code with email service
	_, err = vu.esc.VerifyCode(ctx, &pb.VerifyCodeRequest{
		Email: user.Email,
		Code:  params.Code,
	})
	// Handle error
	if err != nil {
		if errsError, ok := err.(*errs.Error); ok {
			return errs.B(errsError).Code(errsError.Code).Msg("failed to verify code").Err()
		}
		return errs.B(err).Code(errs.Internal).Msg("failed to verify code").Err()
	}
	user.IsVerified = true
	return vu.ur.Update(ctx, user)
}

func NewVerifyUserCommand(v Validator, ur UserRepository, esc pb.EmailServiceClient) VerifyUserCommand {
	return &VerifyUserCommandImpl{v: v, ur: ur, esc: esc}
}
