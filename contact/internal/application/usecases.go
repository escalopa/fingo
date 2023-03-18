package application

import "time"

type UseCases struct {
	v   Validator
	es  EmailSender
	sci time.Duration // Send code interval
	spi time.Duration // Send password interval

	Command
}

// NewUseCases creates a new use cases instance.
func NewUseCases(opts ...func(cases *UseCases)) *UseCases {
	uc := &UseCases{}
	for _, opt := range opts {
		opt(uc)
	}
	// initialize commands
	uc.SendVerificationCode = NewSendVerificationCodeCommand(uc.v, uc.es, uc.sci)
	uc.SendResetPasswordToken = NewSendResetPasswordTokenCommand(uc.v, uc.es, uc.spi)
	uc.SendNewSignInSession = NewSendNewSingInSessionCommand(uc.v, uc.es)
	return uc
}

// WithValidator is an option function to set the validator for the use cases.
func WithValidator(v Validator) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.v = v
	}
}

// WithEmailSender is an option function to set the email sender for the use cases.
func WithEmailSender(es EmailSender) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.es = es
	}
}

// WithMinSendCodeInterval is an option function to set the minimum interval between
func WithMinSendCodeInterval(sci time.Duration) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.sci = sci
	}
}

// WithMinSendPasswordTokenInterval is an option function to set the minimum interval between
func WithMinSendPasswordTokenInterval(spi time.Duration) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.spi = spi
	}
}

// Command is a struct that contains all the commands.
// that can be executed by the use cases.
type Command struct {
	SendVerificationCode   SendVerificationCodeCommand
	SendResetPasswordToken SendResetPasswordTokenCommand
	SendNewSignInSession   SendNewSingInSessionCommand
}
