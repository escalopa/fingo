package interceptors

import (
	"context"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc"
)

// TokenUnaryInterceptor returns a UnaryServerInterceptor that validates the access token and set the user id in the context.
// unauthorizedRequests is a list of requests that do not require authentication.
// tokenValidator is a function that validates the access token and returns the user id.
// This function is going to be used by all the services that need to validate the access token in a unary request.
func TokenUnaryInterceptor(
	unauthorizedRequests []string,
	tokenValidator func(ctx context.Context, token string) (string, error),
) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract clientIP & userAgent of the current request
		clientIP, userAgent := contextutils.GetMetadata(ctx)
		ctx = contextutils.SetForwardMetadata(ctx, clientIP, userAgent)
		// check if request is unauthorized and skip auth
		for _, request := range unauthorizedRequests {
			if info.FullMethod == request {
				return handler(ctx, req)
			}
		}
		// Get access token from context
		accessToken, err := contextutils.GetAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		userID, err := tokenValidator(ctx, accessToken)
		if err != nil {
			return nil, errs.B(err).Code(errs.Unauthenticated).Msg("failed to validate token").Err()
		}
		// Set user-id in context
		ctx = contextutils.SetUserID(ctx, userID)
		return handler(ctx, req)
	}
}
