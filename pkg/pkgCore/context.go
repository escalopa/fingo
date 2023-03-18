package pkgCore

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"strings"
)

const (
	ContextKeyUserID    = "user-uuid"
	ContextKeyClientIP  = "client-ip"
	ContextKeyUserAgent = "user-agent"

	NewContextKeyClientIP  = "my-client-ip"
	NewContextKeyUserAgent = "my-user-agent"
)

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if meta, ok := metadata.FromIncomingContext(ctx); ok {
		userIds := meta.Get(ContextKeyUserID)
		if len(userIds) > 0 {
			id := userIds[0]
			userID, err := uuid.Parse(id)
			if err != nil {
				return uuid.UUID{}, errs.B(err).Code(errs.Internal).Msg("failed to parse user id from headers").Err()
			}
			return userID, nil
		}
	}
	if id, ok := ctx.Value(ContextKeyUserID).(string); ok {
		userID, err := uuid.Parse(id)
		if err != nil {
			return uuid.UUID{}, errs.B(err).Code(errs.Internal).Msg("failed to parse user id from headers").Err()
		}
		return userID, nil
	}
	return uuid.UUID{}, errs.B().Code(errs.Unauthenticated).Msg("user id not passed in headers").Err()
}

func GetMDFromContext(ctx context.Context) (clientIP, userAgent string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgents := md.Get(ContextKeyUserAgent); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}
		if clientIPs := md.Get(ContextKeyClientIP); len(clientIPs) > 0 {
			clientIP = clientIPs[0]
		}
		//// Update values if new context key is set
		if userAgents := md.Get(NewContextKeyUserAgent); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}
		if clientIPs := md.Get(NewContextKeyClientIP); len(clientIPs) > 0 {
			clientIP = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	clientIP = formatClientIP(clientIP)
	return
}

const (
	authorizationHeader = "authorization"
	authorizationType   = "Bearer"
)

func GetAccessTokenFromContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if tokens := md.Get(authorizationHeader); len(tokens) > 0 {
			if strings.HasPrefix(tokens[0], authorizationType) {
				return strings.TrimPrefix(tokens[0], authorizationType+" "), nil
			}
			return "", errs.B().Code(errs.InvalidArgument).Msg("invalid authorization header").Err()
		}
	}
	return "", errs.B().Code(errs.Unauthenticated).Msg("missing authorization header").Err()
}

func formatClientIP(ip string) string {
	twoDots := strings.Count(ip, ":")
	if twoDots > 1 && strings.Contains(ip, "[") {
		ip = ip[:strings.LastIndex(ip, ":")]
		if ip[0] == '[' {
			ip = ip[1:]
		}
		if ip[len(ip)-1] == ']' {
			ip = ip[:len(ip)-1]
		}
	} else if twoDots == 1 {
		ip = ip[:strings.LastIndex(ip, ":")]
	}
	return ip
}
