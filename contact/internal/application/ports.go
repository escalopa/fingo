package application

import (
	"context"
	"github.com/escalopa/fingo/contact/internal/core"
)

type EmailSender interface {
	SendVerificationCode(ctx context.Context, params core.SendVerificationCodeMessage) error
	SendResetPasswordToken(ctx context.Context, params core.SendResetPasswordTokenMessage) error
	SendNewSignInSession(ctx context.Context, params core.SendNewSignInSessionMessage) error
	Close() error
}

type Validator interface {
	Validate(i interface{}) error
}

type MessageConsumer interface {
	HandleSendVerificationsCode(handler func(ctx context.Context, params core.SendVerificationCodeMessage) error) error
	HandleSendResetPasswordToken(handler func(ctx context.Context, params core.SendResetPasswordTokenMessage) error) error
	HandleSendNewSignInSession(handler func(ctx context.Context, params core.SendNewSignInSessionMessage) error) error
	Close() error
}

type Server interface {
	Start() error
	Stop()
}
