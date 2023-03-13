package application

import (
	"context"
	"github.com/escalopa/fingo/email/internal/core"
	"time"
)

type SendVerificationCodeCommandParam struct {
	Name  string `validate:"required,alpha,min=2,max=50"`
	Email string `validate:"required,email"`
	Code  string `validate:"required"`
}

type SendVerificationCodeCommand interface {
	Execute(ctx context.Context, params SendVerificationCodeCommandParam) error
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

func (c *SendVerificationCodeCommandImpl) Execute(ctx context.Context, params SendVerificationCodeCommandParam) error {
	// TODO: check if the user has sent verification code request in the last `c.sci` if so, return an error
	if err := c.v.Validate(params); err != nil {
		return err
	}
	err := c.es.SendVerificationCode(ctx, core.SendVerificationCodeMessage{
		Name:  params.Name,
		Email: params.Email,
		Code:  params.Code,
	})
	if err != nil {
		return err
	}
	return nil
}
