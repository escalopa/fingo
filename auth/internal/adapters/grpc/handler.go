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

func (h *AuthHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	err := h.uc.Signup.Execute(ctx, application.SignupParams{
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

func (h *AuthHandler) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	token, err := h.uc.Signin.Execute(ctx, application.SigninParams{
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

func (h *AuthHandler) SendUserCode(ctx context.Context, req *pb.SendUserCodeRequest) (*pb.SendUserCodeResponse, error) {
	err := h.uc.SendUserCode.Execute(ctx, application.SendUserCodeParam{Email: req.Email})
	if err != nil {
		return nil, err
	}
	return &pb.SendUserCodeResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Verification successful"},
	}, nil
}

func (h *AuthHandler) VerifyUserCode(ctx context.Context, req *pb.VerifyUserCodeRequest) (*pb.VerifyUserCodeResponse, error) {
	err := h.uc.VerifyUserCode.Execute(ctx, application.VerifyUserCodeParam{Email: req.Email, Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyUserCodeResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Verification successful"},
	}, nil
}

func (h *AuthHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	token := req.Token
	user, err := h.uc.VerifyToken.Execute(ctx, application.VerifyTokenParams{Token: token})
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResponse{
		Response: &pb.BasicResponse{Status: 200, Message: "Token verification successful"},
		User:     &pb.User{Email: user.Email, Id: user.ID},
	}, nil
}

func toStrPtr(s string) *string {
	return &s
}
