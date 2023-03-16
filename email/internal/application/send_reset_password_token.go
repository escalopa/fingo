package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/email/internal/core"
)

type SendResetPasswordTokenCommandParam struct {
	Name  string `validate:"required,alpha,min=2,max=50"`
	Email string `validate:"required,email"`
	Token string `validate:"required"`
}

type SendResetPasswordTokenCommand interface {
	Execute(ctx context.Context, params SendResetPasswordTokenCommandParam) error
}

type SendResetPasswordTokenCommandImpl struct {
	v   Validator
	es  EmailSender
	spi time.Duration
}

func NewSendResetPasswordTokenCommand(v Validator, es EmailSender, spi time.Duration) SendResetPasswordTokenCommand {
	return &SendResetPasswordTokenCommandImpl{
		v:   v,
		es:  es,
		spi: spi,
	}
}

func (c *SendResetPasswordTokenCommandImpl) Execute(ctx context.Context, params SendResetPasswordTokenCommandParam) error {
	// TODO: check if the user has requested a password reset token in the last `c.spi` if so, return an error
	if err := c.v.Validate(params); err != nil {
		return err
	}
	err := c.es.SendResetPasswordToken(ctx, core.SendResetPasswordTokenMessage{
		Name:  params.Name,
		Email: params.Email,
		Token: params.Token,
	})
	if err != nil {
		return err
	}
	return nil
}
