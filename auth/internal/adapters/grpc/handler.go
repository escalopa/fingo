package grpc

import (
	"context"
	"github.com/escalopa/gochat/auth/internal/application"
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/escalopa/gochat/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ----------------------------------------- //
// -------------- AuthHandler -------------- //
// ----------------------------------------- //

type AuthHandler struct {
	uc *application.UseCases
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(uc *application.UseCases) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	err := h.uc.Signup.Execute(ctx, application.SignupParams{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SignupResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Successfully signed-up"},
	}, nil
}

func (h *AuthHandler) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	token, err := h.uc.Signin.Execute(ctx, application.SigninParams{
		Email: req.Email, Password: req.Password, MetaData: req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SigninResponse{
		AccessToken: token,
		Response:    &pb.BasicResponse{Status: 200, Message: "Successfully signed-in"},
	}, nil
}

func (h *AuthHandler) SendUserCode(ctx context.Context, req *pb.SendUserCodeRequest) (*pb.SendUserCodeResponse, error) {
	err := h.uc.SendUserCode.Execute(ctx, application.SendUserCodeParam{Email: req.Email})
	if err != nil {
		return nil, err
	}
	return &pb.SendUserCodeResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Successfully sent code to user email"},
	}, nil
}

func (h *AuthHandler) VerifyUserCode(ctx context.Context, req *pb.VerifyUserCodeRequest) (*pb.VerifyUserCodeResponse, error) {
	err := h.uc.VerifyUserCode.Execute(ctx, application.VerifyUserCodeParam{Email: req.Email, Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyUserCodeResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Successfully verified user email"},
	}, nil
}

func (h *AuthHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	user, err := h.uc.VerifyToken.Execute(ctx, application.VerifyTokenParams{AccessToken: req.AccessToken})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResponse{
		User:     fromUserToPb(user),
		Response: &pb.BasicResponse{Status: 200, Message: "Successfully verified access token"},
	}, nil
}

func (h *AuthHandler) RenewToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	newAccessToken, err := h.uc.RenewToken.Execute(ctx, application.RenewTokenParams{RefreshToken: req.RefreshToken})
	if err != nil {
		return nil, err
	}
	return &pb.RenewAccessTokenResponse{
		NewAccessToken: newAccessToken,
		Response:       &pb.BasicResponse{Status: 200, Message: "Successfully renewed access token"},
	}, nil
}

func fromUserToPb(u core.User) *pb.User {
	return &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
}
