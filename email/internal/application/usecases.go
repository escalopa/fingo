package application

import "time"

type UseCases struct {
	cr  CodeRepository
	cs  CodeGenerator
	es  EmailSender
	mti time.Duration // min time interval between sending codes

	SendCode   SendCodeCommand
	VerifyCode VerifyCodeCommand
}

// NewUseCases creates a new use cases instance.
func NewUseCases(opts ...func(cases *UseCases)) *UseCases {
	// initialize use cases
	cases := &UseCases{}
	// apply functional options
	for _, opt := range opts {
		opt(cases)
	}

	// initialize commands
	cases.SendCode = NewSendCodeCommand(cases.mti, cases.cr, cases.cs, cases.es)
	cases.VerifyCode = NewVerifyCodeCommand(cases.cr, cases.cs)
	return cases
}

// WithCodeRepository is a functional option to set the code repository
// for the use cases.
func WithCodeRepository(cr CodeRepository) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.cr = cr
	}
}

// WithCodeGenerator is a functional option to set the code generator
// for the use cases.
func WithCodeGenerator(cs CodeGenerator) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.cs = cs
	}
}

// WithEmailSender is a functional option to set the email sender
// for the use cases.
func WithEmailSender(es EmailSender) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.es = es
	}
}

// WithMinTimeInterval is a functional option to set the minimum time interval
// between sending codes.
func WithMinTimeInterval(mti time.Duration) func(cases *UseCases) {
	return func(cases *UseCases) {
		cases.mti = mti
	}
}

// Command is a struct that contains all the commands.
// that can be executed by the use cases.
type Command struct {
	SendCode   SendCodeCommand
	VerifyCode VerifyCodeCommand
}
