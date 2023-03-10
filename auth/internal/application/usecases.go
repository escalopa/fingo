package application

import "github.com/escalopa/fingo/pb"

type UseCases struct {
	v   Validator
	h   PasswordHasher
	tg  TokenGenerator
	ur  UserRepository
	sr  SessionRepository
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
		Signin:         NewSigninCommand(u.v, u.h, u.tg, u.ur, u.sr),
		Signup:         NewSignupCommand(u.v, u.h, u.ur),
		Logout:         NewLogoutCommand(u.v, u.tg, u.ur, u.sr),
		SendUserCode:   NewSendUserCodeCommand(u.v, u.ur, u.esc),
		VerifyUserCode: NewVerifyUserCodeCommand(u.v, u.ur, u.esc),
		VerifyToken:    NewVerifyTokenCommand(u.v, u.tg, u.sr),
		RenewToken:     NewRenewTokenCommand(u.v, u.tg, u.sr),
	}
	return u
}

func WithUserRepository(ur UserRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.ur = ur
	}
}

func WithSessionRepository(sr SessionRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.sr = sr
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
	Logout         LogoutCommand
	SendUserCode   SendUserCodeCommand
	VerifyUserCode VerifyUserCodeCommand
	VerifyToken    VerifyTokenCommand
	RenewToken     RenewTokenCommand
}
