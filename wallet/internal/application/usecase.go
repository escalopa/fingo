package application

type UseCases struct {
	v   Validator
	ur  UserRepository
	ar  AccountRepository
	cr  CardRepository
	tr  TransactionRepository
	ss  SmsSender
	cng CardNumberGenerator

	Command
	Query
}

type UseCasesOption func(*UseCases)

func NewUseCases(opts ...UseCasesOption) *UseCases {
	uc := &UseCases{}
	for _, opt := range opts {
		opt(uc)
	}
	uc.Command = Command{}
	uc.Query = Query{}
	return uc
}

func WithValidator(v Validator) UseCasesOption {
	return func(uc *UseCases) {
		uc.v = v
	}
}

func WithUserRepository(ur UserRepository) UseCasesOption {
	return func(uc *UseCases) {
		uc.ur = ur
	}
}

func WithAccountRepository(ar AccountRepository) UseCasesOption {
	return func(uc *UseCases) {
		uc.ar = ar
	}
}

func WithCardRepository(cr CardRepository) UseCasesOption {
	return func(uc *UseCases) {
		uc.cr = cr
	}
}

func WithTransactionRepository(tr TransactionRepository) UseCasesOption {
	return func(uc *UseCases) {
		uc.tr = tr
	}
}

func WithSmsSender(ss SmsSender) UseCasesOption {
	return func(uc *UseCases) {
		uc.ss = ss
	}
}

func WithCardNumberGenerator(cng CardNumberGenerator) UseCasesOption {
	return func(uc *UseCases) {
		uc.cng = cng
	}
}

type Command struct {
}

type Query struct {
}
