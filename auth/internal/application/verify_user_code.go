package application

import (
	"context"

	"github.com/escalopa/gofly/pb"
	"github.com/lordvidex/errs"
)

// ---------------------- Verify User ---------------------- //

type SendUserCodeParam struct {
	Email string `validate:"required,email"`
	Code  string `validate:"required,numeric"`
}

type SendUserCodeCommand interface {
	Execute(ctx context.Context, params SendUserCodeParam) error
}

type SendUserCodeCommandImpl struct {
	v   Validator
	ur  UserRepository
	esc pb.EmailServiceClient
}

func (vu *SendUserCodeCommandImpl) Execute(ctx context.Context, params SendUserCodeParam) error {
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
		Email: params.Email,
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

func NewSendUserCodeCommand(v Validator, ur UserRepository, esc pb.EmailServiceClient) SendUserCodeCommand {
	return &SendUserCodeCommandImpl{v: v, ur: ur, esc: esc}
}
