package application

import (
	"context"
)

// ---------------------- Verify User ---------------------- //

type SendUserCodeParam struct {
	Email string `validate:"required,email"`
}

type SendUserCodeCommand interface {
	Execute(ctx context.Context, params SendUserCodeParam) error
}

type SendUserCodeCommandImpl struct {
	v Validator
}

func (vu *SendUserCodeCommandImpl) Execute(ctx context.Context, params SendUserCodeParam) error {
	//if err := vu.v.Validate(params); err != nil {
	//	return err
	//}
	//user, err := vu.ur.GetUserByEmail(ctx, params.Email)
	//if err != nil {
	//	return err
	//}
	//// Check if the user is already verified
	//if user.IsVerified {
	//	return errs.B().Code(errs.InvalidArgument).Msg("user is already verified").Err()
	//}
	//// VerifyToken code with email service
	//_, err = vu.esc.SendCode(ctx, &pb.SendCodeRequest{
	//	Email: params.Email,
	//})
	//// Handle error
	//if err != nil {
	//	if er, ok := err.(*errs.Error); ok {
	//		return errs.B(er).Code(er.Code).Msg("failed to send code").Err()
	//	}
	//	return errs.B(err).Code(errs.Internal).Msg("failed to send code").Err()
	//}
	return nil
}

func NewSendUserCodeCommand(v Validator) SendUserCodeCommand {
	return &SendUserCodeCommandImpl{v: v}
}
