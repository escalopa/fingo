package application

import (
	"context"
)

type SendNewSingInSessionCommandParam struct {
	Name      string `validate:"required,alpha,min=2,max=50"`
	Email     string `validate:"required,email"`
	ClientIP  string `validate:"required,ip"`
	UserAgent string `validate:"required"`
}

type SendNewSingInSessionCommand interface {
	Execute(ctx context.Context, param SendNewSingInSessionCommandParam) error
}

type SendNewSingInSessionCommandImpl struct {
	v  Validator
	es EmailSender
}

func NewSendNewSingInSessionCommand(v Validator, es EmailSender) SendNewSingInSessionCommand {
	return &SendNewSingInSessionCommandImpl{
		v:  v,
		es: es,
	}
}

func (c *SendNewSingInSessionCommandImpl) Execute(ctx context.Context, param SendNewSingInSessionCommandParam) error {
	if err := c.v.Validate(param); err != nil {
		return err
	}
	if err := c.es.SendNewSignInSession(ctx, param.Name, param.Email, param.ClientIP, param.UserAgent); err != nil {
		return err
	}
	return nil
}
