package application

import (
	"context"
	"github.com/escalopa/fingo/contact/internal/core"
	"testing"
	"time"

	"github.com/escalopa/fingo/contact/internal/adapters/validator"
)

var testUseCases *UseCases

type emailSenderMock struct {
}

func (esm *emailSenderMock) SendVerificationCode(_ context.Context, _ core.SendVerificationCodeMessage) error {
	return nil
}
func (esm *emailSenderMock) SendResetPasswordToken(_ context.Context, _ core.SendResetPasswordTokenMessage) error {
	return nil
}
func (esm *emailSenderMock) SendNewSignInSession(_ context.Context, _ core.SendNewSignInSessionMessage) error {
	return nil
}

func (esm *emailSenderMock) Close() error { return nil }

func TestMain(m *testing.M) {
	v := validator.NewValidator()
	testUseCases = NewUseCases(
		WithValidator(v),
		WithEmailSender(&emailSenderMock{}),
		WithMinSendCodeInterval(1*time.Second),
		WithMinSendPasswordTokenInterval(1*time.Second),
	)
	m.Run()
}
