package contextutils

import (
	"context"
	"strings"

	"github.com/lordvidex/errs"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "Bearer"
)

func GetAccessToken(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		tokens := md.Get(authorizationHeader)
		if len(tokens) == 0 {
			return "", errs.B().Code(errs.Unauthenticated).Msg("missing authorization header").Err()
		}
		token := tokens[0]
		if strings.HasPrefix(token, authorizationType) {
			return strings.TrimPrefix(token, authorizationType+" "), nil
		} else {
			return "", errs.B().Code(errs.Unauthenticated).Msg("invalid authorization header").Err()
		}
	}
	return "", errs.B().Code(errs.Unauthenticated).Msg("missing metadata for authorization").Err()
}
