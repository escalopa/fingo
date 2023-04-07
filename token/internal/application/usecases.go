package application

type UseCases struct {
	v  Validator
	tr TokenRepository

	Command
}

func NewUseCases(opts ...func(*UseCases)) *UseCases {
	u := &UseCases{}
	for _, opt := range opts {
		opt(u)
	}
	u.Command = Command{
		TokenValidate: NewTokenValidateCommand(u.v, u.tr),
	}
	return u
}

func WithTokenRepository(tr TokenRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.tr = tr
	}
}

func WithValidator(v Validator) func(*UseCases) {
	return func(u *UseCases) {
		u.v = v
	}
}

type Command struct {
	TokenValidate TokenValidateCommand
}
