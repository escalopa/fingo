package grpc

import (
	"context"
	"github.com/escalopa/gofly/auth/internal/application"
	"github.com/escalopa/gofly/pb"
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

func (h *AuthHandler) Signup(_ context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	err := h.uc.Signup.Execute(application.SignupParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SignupResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Signup successful"},
	}, nil
}

func (h *AuthHandler) Signin(_ context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	token, err := h.uc.Signin.Execute(application.SigninParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SigninResponse{
		Response:    &pb.BasicResponse{Status: 200, Message: "Signin successful"},
		AccessToken: toStrPtr(token),
	}, nil
}

func (h *AuthHandler) Verify(_ context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	err := h.uc.VerifyUser.Execute(application.VerifyUserParam{Email: req.Email, Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyUserResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Verification successful"},
	}, nil
}

// ------------------------------------------ //
// -------------- TokenHandler -------------- //
// ------------------------------------------ //

type TokenHandler struct {
	uc *application.UseCases
	pb.UnimplementedTokenServiceServer
}

func NewTokenHandler(uc *application.UseCases) *TokenHandler {
	return &TokenHandler{uc: uc}
}

func (h *TokenHandler) Verify(_ context.Context, req *pb.VerifyTokenRequest) (*pb.User, error) {
	token := req.Token
	user, err := h.uc.VerifyToken.Execute(application.VerifyTokenParams{Token: token})
	if err != nil {
		return nil, err
	}
	return &pb.User{Id: user.ID, Email: user.Email}, nil
}

func toStrPtr(s string) *string {
	return &s
}
