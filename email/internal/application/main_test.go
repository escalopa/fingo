package application

import (
	"context"
	"testing"
	"time"

	"github.com/escalopa/fingo/email/internal/adapters/validator"
)

var testUseCases *UseCases

type emailSenderMock struct {
}

func (esm *emailSenderMock) SendVerificationCode(ctx context.Context, email string, name string, code string) error {
	return nil
}
func (esm *emailSenderMock) SendResetPasswordToken(ctx context.Context, email string, name string, token string) error {
	return nil
}
func (esm *emailSenderMock) SendNewSignInSession(ctx context.Context, email string, name string, clientIP string, userAgent string) error {
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
