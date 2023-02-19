package application

import "github.com/escalopa/gofly/contact/internal/core"

type CodeRepository interface {
	Save(email string, vc core.VerificationCode) error
	Get(email string) (core.VerificationCode, error)
	Close() error
}

type EmailSender interface {
	SendVerificationCode(email string, vc core.VerificationCode) error
	Close() error
}

type CodeGenerator interface {
	GenerateCode() (string, error)
	VerifyCode(code string) bool
}
