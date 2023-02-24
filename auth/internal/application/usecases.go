package application

import "github.com/escalopa/gochat/pb"

type UseCases struct {
	v   Validator
	h   PasswordHasher
	tg  TokenGenerator
	ur  UserRepository
	esc pb.EmailServiceClient

	Query
	Command
}

func NewUseCases(opts ...func(*UseCases)) *UseCases {
	u := &UseCases{}
	for _, opt := range opts {
		opt(u)
	}
	u.Query = Query{}
	u.Command = Command{
		Signin:         NewSigninCommand(u.v, u.h, u.tg, u.ur),
		Signup:         NewSignupCommand(u.v, u.h, u.ur),
		SendUserCode:   NewSendUserCodeCommand(u.v, u.ur, u.esc),
		VerifyUserCode: NewVerifyUserCodeCommand(u.v, u.ur, u.esc),
		VerifyToken:    NewVerifyTokenCommand(u.v, u.tg),
		RenewToken:     NewRenewTokenCommand(u.v, u.tg),
	}
	return u
}

func WithUserRepository(ur UserRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.ur = ur
	}
}

func WithTokenGenerator(tg TokenGenerator) func(*UseCases) {
	return func(u *UseCases) {
		u.tg = tg
	}
}

func WithPasswordHasher(h PasswordHasher) func(*UseCases) {
	return func(u *UseCases) {
		u.h = h
	}
}

func WithEmailService(esc pb.EmailServiceClient) func(*UseCases) {
	return func(u *UseCases) {
		u.esc = esc
	}
}

func WithValidator(v Validator) func(*UseCases) {
	return func(u *UseCases) {
		u.v = v
	}
}

type Query struct{}

type Command struct {
	Signin         SigninCommand
	Signup         SignupCommand
	SendUserCode   SendUserCodeCommand
	VerifyUserCode VerifyUserCodeCommand
	VerifyToken    VerifyTokenCommand
	RenewToken     RenewTokenCommand
}
