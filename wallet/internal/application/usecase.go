package application

import "github.com/lordvidex/errs"

type UseCases struct {
	v   Validator
	l   Locker
	ur  UserRepository
	ar  AccountRepository
	cr  CardRepository
	tr  TransactionRepository
	ss  SmsSender
	cng CardNumberGenerator

	command
	query
}

var (
	errorNotAccountOwner = errs.B().Code(errs.Forbidden).Msg("failed to delete account, not account owner").Err()
)

type UseCasesOption func(*UseCases)

func NewUseCases(opts ...UseCasesOption) *UseCases {
	uc := &UseCases{}
	for _, opt := range opts {
		opt(uc)
	}
	uc.command = command{
		CreateWallet:      NewCreateWalletCommand(uc.v, uc.ur),
		CreateAccount:     NewCreateAccountCommand(uc.v, uc.ur, uc.ar),
		DeleteAccount:     NewDeleteAccountCommand(uc.v, uc.ur, uc.ar),
		CreateCard:        NewCreateCardCommand(uc.v, uc.ur, uc.ar, uc.cr, uc.cng),
		DeleteCard:        NewDeleteCardCommand(uc.v, uc.ur, uc.ar, uc.cr),
		CreateTransaction: NewCreateTransactionCommand(uc.v, uc.l, uc.ur, uc.ar, uc.cr, uc.tr),
		TransferRollback:  NewTransferRollbackCommand(uc.v, uc.l, uc.ur, uc.ar, uc.tr),
	}
	uc.query = query{
		GetAccounts:           NewGetAccountsCommand(uc.v, uc.ur, uc.ar),
		GetCards:              NewGetCardsCommand(uc.v, uc.ur, uc.ar, uc.cr),
		GetTransactionHistory: NewGetTransactionHistoryCommand(uc.v, uc.ur, uc.ar, uc.tr),
	}
	return uc
}

func WithValidator(v Validator) UseCasesOption {
	return func(uc *UseCases) {
		uc.v = v
	}
}

func WithLocker(l Locker) UseCasesOption {
	return func(uc *UseCases) {
		uc.l = l
	}
}

func WithSmsSender(ss SmsSender) UseCasesOption {
	return func(uc *UseCases) {
		uc.ss = ss
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

func WithCardNumberGenerator(cng CardNumberGenerator) UseCasesOption {
	return func(uc *UseCases) {
		uc.cng = cng
	}
}

type command struct {
	CreateWallet      CreateWalletCommand
	CreateAccount     CreateAccountCommand
	DeleteAccount     DeleteAccountCommand
	CreateCard        CreateCardCommand
	DeleteCard        DeleteCardCommand
	CreateTransaction CreateTransactionCommand
	TransferRollback  TransferRollbackCommand
}

type query struct {
	GetAccounts           GetAccountsCommand
	GetCards              GetCardsCommand
	GetTransactionHistory GetTransactionHistoryCommand
}
