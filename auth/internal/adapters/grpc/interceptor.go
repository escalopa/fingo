package grpc

import (
	"context"

	"github.com/escalopa/fingo/pb"
	"github.com/escalopa/fingo/pkg/pkgCore"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	unauthorizedRequests = []string{
		"/pb.AuthService/Signup",
		"/pb.AuthService/Signin",
	}
)

type AuthInterceptor struct {
	c pb.TokenServiceClient
}

func NewAuthInterceptor(url string, creds credentials.TransportCredentials) (*AuthInterceptor, error) {
	client, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to connect to token service").Err()
	}
	ai := &AuthInterceptor{pb.NewTokenServiceClient(client)}
	return ai, nil
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// check if request is unauthorized and skip auth
		for _, path := range unauthorizedRequests {
			if info.FullMethod == path {
				return handler(ctx, req)
			}
		}
		// Get access token from context
		accessToken, err := pkgCore.GetAccessTokenFromContext(ctx)
		clientIP, userAgent := pkgCore.GetMDFromContext(ctx)
		out := metadata.NewOutgoingContext(ctx, metadata.Pairs(
			pkgCore.NewContextKeyClientIP, clientIP,
			pkgCore.NewContextKeyUserAgent, userAgent,
		))
		if err != nil {
			return nil, err
		}
		userID, err := ai.c.ValidateToken(out, &pb.ValidateTokenRequest{AccessToken: accessToken})
		if err != nil {
			return nil, errs.B(err).Code(errs.Unauthenticated).Msg("failed to validate token").Err()
		}
		// Set user-id in context
		out = context.WithValue(out, pkgCore.ContextKeyUserID, userID.UserId) // nolint:all
		return handler(out, req)
	}
}
