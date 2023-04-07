package application

type UseCases struct {
	v  Validator
	h  PasswordHasher
	tg TokenGenerator
	ur UserRepository
	sr SessionRepository
	tr TokenRepository
	mp MessageProducer

	Query
	Command
}

func NewUseCases(opts ...func(*UseCases)) *UseCases {
	u := &UseCases{}
	for _, opt := range opts {
		opt(u)
	}
	u.Query = Query{
		GetUserDevices: NewGetUserDevicesCommand(u.v, u.sr),
	}
	u.Command = Command{
		Signin:     NewSigninCommand(u.v, u.h, u.tg, u.ur, u.sr, u.tr, u.mp),
		Signup:     NewSignupCommand(u.v, u.h, u.ur),
		Logout:     NewLogoutCommand(u.v, u.sr, u.tr),
		RenewToken: NewRenewTokenCommand(u.v, u.tg, u.sr, u.tr),
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

func WithTokenRepository(tr TokenRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.tr = tr
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

func WithMessageProducer(mp MessageProducer) func(*UseCases) {
	return func(u *UseCases) {
		u.mp = mp
	}
}

func WithValidator(v Validator) func(*UseCases) {
	return func(u *UseCases) {
		u.v = v
	}
}

type Query struct {
	GetUserDevices GetUserDevicesCommand
}

type Command struct {
	Signin     SigninCommand
	Signup     SignupCommand
	Logout     LogoutCommand
	RenewToken RenewTokenCommand
}
