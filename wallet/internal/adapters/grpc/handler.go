package grpc

import (
	"context"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/wallet/internal/application"
)

type WalletHandler struct {
	u *application.UseCases
	pb.UnimplementedWalletServiceServer
}

func NewWalletHandler(u *application.UseCases) *WalletHandler {
	return &WalletHandler{u: u}
}

func (wh *WalletHandler) CreateWallet(context.Context, *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	return nil, nil
}

func (wh *WalletHandler) CreateAccount(context.Context, *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return nil, nil
}
func (wh *WalletHandler) GetAccounts(context.Context, *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	return nil, nil
}
func (wh *WalletHandler) DeleteAccount(context.Context, *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	return nil, nil
}

func (wh *WalletHandler) CreateCard(context.Context, *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	return nil, nil
}
func (wh *WalletHandler) GetCards(context.Context, *pb.GetCardsRequest) (*pb.GetCardsResponse, error) {
	return nil, nil
}
func (wh *WalletHandler) DeleteCard(context.Context, *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	return nil, nil
}

func (wh *WalletHandler) Transfer(context.Context, *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}
func (wh *WalletHandler) TransferRollback(context.Context, *pb.TransferRollbackRequest) (*pb.TransferRollbackResponse, error) {
	return nil, nil
}

func (wh *WalletHandler) GetTransactionHistory(*pb.GetTransactionHistoryRequest, pb.WalletService_GetTransactionHistoryServer) error {
	return nil
}
