package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"

	"github.com/escalopa/fingo/pb"
	"github.com/lordvidex/errs"
)

// ---------------------- Verify User ---------------------- //

type VerifyUserCodeParam struct {
	Email string `validate:"required,email"`
	Code  string `validate:"required,numeric"`
}

type VerifyUserCodeCommand interface {
	Execute(ctx context.Context, params VerifyUserCodeParam) error
}

type VerifyUserCodeCommandImpl struct {
	v   Validator
	ur  UserRepository
	esc pb.EmailServiceClient
}

func (vu *VerifyUserCodeCommandImpl) Execute(ctx context.Context, params VerifyUserCodeParam) error {
	if err := vu.v.Validate(params); err != nil {
		return err
	}
	user, err := vu.ur.GetUserByEmail(ctx, params.Email)
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
	// Update user to be verified
	err = vu.ur.SetUserIsVerified(ctx, core.SetUserIsVerifiedParams{
		ID:         user.ID,
		IsVerified: true,
	})
	if err != nil {
		return err
	}
	return err
}

func NewVerifyUserCodeCommand(v Validator, ur UserRepository, esc pb.EmailServiceClient) VerifyUserCodeCommand {
	return &VerifyUserCodeCommandImpl{v: v, ur: ur, esc: esc}
}
