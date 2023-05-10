package grpc

import (
	"context"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/escalopa/fingo/wallet/internal/application"
	"github.com/escalopa/fingo/wallet/internal/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletHandler struct {
	u *application.UseCases
	pb.UnimplementedWalletServiceServer
}

func NewWalletHandler(u *application.UseCases) *WalletHandler {
	return &WalletHandler{u: u}
}

func (wh *WalletHandler) CreateWallet(ctx context.Context, _ *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.CreateWallet")
	defer span.End()
	err := wh.u.CreateWallet.Execute(ctx, application.CreateWalletParams{})
	if err != nil {
		return nil, err
	}
	return &pb.CreateWalletResponse{Success: true}, nil
}

func (wh *WalletHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.CreateAccount")
	defer span.End()
	err := wh.u.CreateAccount.Execute(ctx, application.CreateAccountParams{
		Name:     req.Name,
		Currency: req.Currency.String(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateAccountResponse{Success: true}, nil
}

func (wh *WalletHandler) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.GetAccounts")
	defer span.End()
	accounts, err := wh.u.GetAccounts.Execute(ctx, application.GetAccountsParams{})
	if err != nil {
		return nil, err
	}
	// Convert to pb type
	pbAccounts := make([]*pb.GetAccountsResponse_Account, len(accounts))
	for i, a := range accounts {
		pbAccounts[i] = fromCoreAccount(a)
	}
	return &pb.GetAccountsResponse{Accounts: pbAccounts}, nil
}

func (wh *WalletHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.DeleteAccount")
	defer span.End()
	err := wh.u.DeleteAccount.Execute(ctx, application.DeleteAccountParams{AccountID: req.AccountId})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAccountResponse{Success: true}, nil
}

func (wh *WalletHandler) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.CreateCard")
	defer span.End()
	err := wh.u.CreateCard.Execute(ctx, application.CreateCardParams{AccountID: req.AccountId})
	if err != nil {
		return nil, err
	}
	return &pb.CreateCardResponse{Success: true}, nil
}

func (wh *WalletHandler) GetCards(ctx context.Context, req *pb.GetCardsRequest) (*pb.GetCardsResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.GetCards")
	defer span.End()
	cards, err := wh.u.GetCards.Execute(ctx, application.GetCardsParams{AccountID: req.AccountId})
	if err != nil {
		return nil, err
	}
	// Convert to pb type
	pbCards := make([]*pb.GetCardsResponse_Card, len(cards))
	for i, c := range cards {
		pbCards[i] = &pb.GetCardsResponse_Card{Number: c.Number}
	}
	return &pb.GetCardsResponse{Cards: pbCards}, nil
}

func (wh *WalletHandler) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.DeleteCard")
	defer span.End()
	err := wh.u.DeleteCard.Execute(ctx, application.DeleteCardParams{CardNumber: req.CardNumber})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCardResponse{Success: true}, nil
}

func (wh *WalletHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.Transfer")
	defer span.End()
	err := wh.u.CreateTransaction.Execute(ctx, application.CreateTransactionParams{
		Amount:   req.Amount,
		Type:     core.ParseTransactionType(req.Type.String()),
		FromCard: req.CardNumber,
		ToCard:   req.GetRecipientCardNumber(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTransactionResponse{Success: true}, nil
}

func (wh *WalletHandler) TransferRollback(ctx context.Context, req *pb.TransferRollbackRequest) (*pb.TransferRollbackResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.TransferRollback")
	defer span.End()
	err := wh.u.TransferRollback.Execute(ctx, application.TransferRollbackParams{TransactionID: req.TransactionId})
	if err != nil {
		return nil, err
	}
	return &pb.TransferRollbackResponse{Success: true}, nil
}

func (wh *WalletHandler) GetTransactionHistory(ctx context.Context, req *pb.GetTransactionHistoryRequest) (*pb.GetTransactionHistoryResponse, error) {
	ctx, span := tracer.Tracer().Start(ctx, "WalletHandler.GetTransactionHistory")
	defer span.End()
	transactions, err := wh.u.GetTransactionHistory.Execute(ctx, application.GetTransactionHistoryParams{
		AccountID:       req.AccountId,
		Offset:          req.Offset,
		Limit:           req.Limit,
		MinAmount:       req.GetMinAmount(),
		MaxAmount:       req.GetMaxAmount(),
		TransactionType: core.ParseTransactionType(req.GetTransactionType().String()),
	})
	if err != nil {
		return nil, err
	}
	// Convert to pb type
	pbTransactions := make([]*pb.GetTransactionHistoryResponse_Transaction, len(transactions))
	for i, t := range transactions {
		pbTransactions[i] = fromCoreTransaction(t)
	}
	return &pb.GetTransactionHistoryResponse{Transactions: pbTransactions}, nil
}

func fromCoreTransaction(t core.Transaction) *pb.GetTransactionHistoryResponse_Transaction {
	return &pb.GetTransactionHistoryResponse_Transaction{
		Id:            t.ID.String(),
		Amount:        t.Amount,
		Type:          fromCoreTransactionType(t.Type),
		SenderName:    t.FromAccountName,
		RecipientName: t.ToAccountName,
		CreatedAt:     timestamppb.New(t.CreatedAt),
		IsRolledBack:  t.IsRolledBack,
	}
}

func fromCoreAccount(a core.Account) *pb.GetAccountsResponse_Account {
	return &pb.GetAccountsResponse_Account{
		Id:       a.ID,
		Name:     a.Name,
		Currency: fromCoreCurrency(a.Currency),
		Balance:  a.Balance,
	}
}

func fromCoreTransactionType(transactionType core.TransactionType) pb.TransactionType {
	switch transactionType {
	case core.TransactionTypeDeposit:
		return pb.TransactionType_DEPOSIT
	case core.TransactionTypeWithdrawal:
		return pb.TransactionType_WITHDRAWAL
	case core.TransactionTypeTransfer:
		return pb.TransactionType_TRANSFER
	default:
		return pb.TransactionType_UNKNOWN
	}
}

func fromCoreCurrency(c core.Currency) pb.Currency {
	switch c {
	case core.CurrencyUSD:
		return pb.Currency_USD
	case core.CurrencyEUR:
		return pb.Currency_EUR
	case core.CurrencyGBP:
		return pb.Currency_GBP
	case core.CurrencyRUB:
		return pb.Currency_RUB
	case core.CurrencyEGP:
		return pb.Currency_EGP
	default:
		return pb.Currency_UNDEFINED
	}
}
