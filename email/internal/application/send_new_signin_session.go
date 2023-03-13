package application

import (
	"context"
	"github.com/escalopa/fingo/email/internal/core"
)

type SendNewSingInSessionCommandParam struct {
	Name      string `validate:"required,alpha,min=2,max=50"`
	Email     string `validate:"required,email"`
	ClientIP  string `validate:"required,ip"`
	UserAgent string `validate:"required"`
}

type SendNewSingInSessionCommand interface {
	Execute(ctx context.Context, params SendNewSingInSessionCommandParam) error
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

func (c *SendNewSingInSessionCommandImpl) Execute(ctx context.Context, params SendNewSingInSessionCommandParam) error {
	if err := c.v.Validate(params); err != nil {
		return err
	}
	err := c.es.SendNewSignInSession(ctx, core.SendNewSignInSessionMessage{
		Name:      params.Name,
		Email:     params.Email,
		ClientIP:  params.ClientIP,
		UserAgent: params.UserAgent,
	})
	if err != nil {
		return err
	}
	return nil
}
