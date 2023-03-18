package grpc

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/application"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/pkgCore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	uc *application.UseCases
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(uc *application.UseCases) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	err := h.uc.Signup.Execute(ctx, application.SignupParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SignupResponse{Message: "user created successfully"}, nil
}

func (h *AuthHandler) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	clientIP, userAgent := pkgCore.GetMDFromContext(ctx)
	response, err := h.uc.Signin.Execute(ctx, application.SigninParams{
		Email:     req.Email,
		Password:  req.Password,
		ClientIP:  clientIP,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SigninResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.uc.Logout.Execute(ctx, application.LogoutParams{SessionID: req.SessionId})
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{Message: "session deleted successfully"}, nil
}

func (h *AuthHandler) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	response, err := h.uc.RenewToken.Execute(ctx, application.RenewTokenParams{RefreshToken: req.RefreshToken})
	if err != nil {
		return nil, err
	}
	return &pb.RenewAccessTokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *AuthHandler) GetUserDevices(ctx context.Context, _ *pb.GetUserDevicesRequest) (*pb.GetUserDevicesResponse, error) {
	sessions, err := h.uc.GetUserDevices.Execute(ctx, application.GetUserDevicesParams{})
	if err != nil {
		return nil, err
	}
	// Convert response to *pb.Session
	pbSessions := make([]*pb.Session, len(sessions))
	for i, v := range sessions {
		pbSessions[i] = fromCoreToPbSession(v)
	}
	return &pb.GetUserDevicesResponse{DevicesSessions: pbSessions}, nil
}

// fromCoreToPbSession convert regular core.Session to *pb.Session
func fromCoreToPbSession(session core.Session) *pb.Session {
	return &pb.Session{
		Id: session.ID.String(),
		UserDevice: &pb.Session_UserDevice{
			ClientIp:  session.UserDevice.ClientIP,
			UserAgent: session.UserDevice.UserAgent,
		},
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
