package application

import (
	"context"
	"github.com/escalopa/gofly/contact/internal/core"
)

type CodeRepository interface {
	Save(ctx context.Context, email string, vc core.VerificationCode) error
	Get(ctx context.Context, email string) (core.VerificationCode, error)
	Close() error
}

type EmailSender interface {
	SendVerificationCode(ctx context.Context, email string, code string) error
	Close() error
}

type CodeGenerator interface {
	GenerateCode() (string, error)
	VerifyCode(code string) bool
}
