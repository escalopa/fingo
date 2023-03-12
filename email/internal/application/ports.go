package application

import (
	"context"
)

type EmailSender interface {
	SendVerificationCode(ctx context.Context, email string, name string, code string) error
	SendResetPasswordToken(ctx context.Context, email string, name string, token string) error
	SendNewSignInSession(ctx context.Context, email string, name string, clientIP string, userAgent string) error
	Close() error
}

type Validator interface {
	Validate(i interface{}) error
}

type MessageQueueConsumer interface {
	HandleSendVerificationsCode(handler func(ctx context.Context, email string, code string) error) error
	HandleSendResetPasswordToken(handler func(ctx context.Context, email string, token string) error) error
	HandleSendNewSignInSession(handler func(ctx context.Context, email string, clientIP string, userAgent string) error) error
}

type Server interface {
	Start() error
	Stop()
}
