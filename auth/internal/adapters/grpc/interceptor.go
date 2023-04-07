package grpc

import (
	"context"

	oteltracer "github.com/escalopa/fingo/auth/internal/adapters/tracer"

	"github.com/escalopa/fingo/pkg/interceptors"

	"github.com/escalopa/fingo/pb"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	// List of requests that do not require authentication
	unauthorizedRequests = []string{
		"/pb.AuthService/Signup",
		"/pb.AuthService/Signin",
	}
)

type AuthInterceptor struct {
	c pb.TokenServiceClient
}

// NewAuthInterceptor returns a new AuthInterceptor
func NewAuthInterceptor(url string, creds credentials.TransportCredentials) (*AuthInterceptor, error) {
	client, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errs.B(err).Code(errs.InvalidArgument).Msg("failed to connect to token service").Err()
	}
	ai := &AuthInterceptor{pb.NewTokenServiceClient(client)}
	return ai, nil
}

// Unary returns a UnaryServerInterceptor that validates the access token and set the user id in the context
func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return interceptors.TokenUnaryInterceptor(
		unauthorizedRequests,
		func(ctx context.Context, accessToken string) (string, error) {
			ctx, span := oteltracer.Tracer().Start(ctx, "ValidateToken")
			defer span.End()
			response, err := ai.c.ValidateToken(ctx, &pb.ValidateTokenRequest{AccessToken: accessToken})
			if err != nil {
				return "", err
			}
			return response.GetUserId(), err
		},
	)
}
