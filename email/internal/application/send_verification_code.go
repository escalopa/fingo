package application

import (
	"context"
	"time"
)

type SendVerificationCodeCommandParam struct {
	Name  string `validate:"required,alpha,min=2,max=50"`
	Email string `validate:"required,email"`
	Code  string `validate:"required"`
}

type SendVerificationCodeCommand interface {
	Execute(ctx context.Context, param SendVerificationCodeCommandParam) error
}

type SendVerificationCodeCommandImpl struct {
	v   Validator
	es  EmailSender
	sci time.Duration
}

func NewSendVerificationCodeCommand(v Validator, es EmailSender, sci time.Duration) SendVerificationCodeCommand {
	return &SendVerificationCodeCommandImpl{
		v:   v,
		es:  es,
		sci: sci,
	}
}

func (c *SendVerificationCodeCommandImpl) Execute(ctx context.Context, param SendVerificationCodeCommandParam) error {
	// TODO: check if the user has sent verification code request in the last `c.sci` if so, return an error
	if err := c.v.Validate(param); err != nil {
		return err
	}
	if err := c.es.SendVerificationCode(ctx, param.Email, param.Name, param.Code); err != nil {
		return err
	}
	return nil
}
