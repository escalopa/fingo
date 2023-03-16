package grpc

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"

	"github.com/escalopa/fingo/auth/internal/application"
	"github.com/escalopa/fingo/pb"
)

type AuthHandler struct {
	uc *application.UseCases
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(uc *application.UseCases) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*emptypb.Empty, error) {
	err := h.uc.Signup.Execute(ctx, application.SignupParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	clientIP, userAgent := extractMetadataFromContext(ctx)
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

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	err := h.uc.Logout.Execute(ctx, application.LogoutParams{SessionID: req.SessionId})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
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

func (h *AuthHandler) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*emptypb.Empty, error) {
	err := h.uc.CreateRole.Execute(ctx, application.CreateRoleParams{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) GrantRole(ctx context.Context, req *pb.GrantRoleRequest) (*emptypb.Empty, error) {
	err := h.uc.GrantRole.Execute(ctx, application.GrantRoleParams{
		UserID:   req.UserId,
		RoleName: req.RoleName,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) RevokeRole(ctx context.Context, req *pb.RevokeRoleRequest) (*emptypb.Empty, error) {
	err := h.uc.RevokeRole.Execute(ctx, application.RevokeRoleParams{
		UserID:   req.UserId,
		RoleName: req.RoleName,
	})
	req.GetRoleName()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
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

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

func extractMetadataFromContext(ctx context.Context) (clientIP, userAgent string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			clientIP = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	// Remove port from ip address
	// Example client ip `192.168.176.6:36234` will become `192.168.176.6`
	twoDotsCount := strings.Count(clientIP, ":")
	if len(clientIP) > 4 {
		if twoDotsCount == 1 { // ipv4
			shard := strings.Split(clientIP, ":")
			clientIP = shard[0]
		} else {
			clientIP = clientIP[1:]
			shard := strings.Split(clientIP, "]")
			clientIP = shard[0]
		}
	}
	return
}
