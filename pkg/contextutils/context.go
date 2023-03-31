package contextutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type contextKey string

func (ck contextKey) String() string {
	return string(ck)
}

const (
	contextKeyUserID contextKey = "user-uuid"

	contextKeyClientIP  contextKey = "client-ip"
	contextKeyUserAgent contextKey = "user-agent"

	forwardContextKeyClientIP  contextKey = "forward-for-client-ip"
	forwardContextKeyUserAgent contextKey = "forward-for-user-agent"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if ids := md.Get(contextKeyUserID.String()); len(ids) > 0 {
			id, err := uuid.Parse(ids[0])
			if err != nil {
				return uuid.UUID{}, errs.B(err).Code(errs.Internal).Msg("invalid user-id").Err()
			}
			return id, nil
		} else {
			return uuid.UUID{}, errs.B().Code(errs.Internal).Msg("failed to parse user id from headers").Err()
		}
	}
	return uuid.UUID{}, errs.B().Code(errs.Internal).Msg("missing metadata to parse user id from headers").Err()
}

func SetUserID(ctx context.Context, userID string) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Append(contextKeyUserID.String(), userID)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(contextKeyUserID.String(), userID))
	}
	return ctx
}

func GetMetadata(ctx context.Context) (clientIP, userAgent string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgents := md.Get(contextKeyUserAgent.String()); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}
		if clientIPs := md.Get(contextKeyClientIP.String()); len(clientIPs) > 0 {
			clientIP = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok && len(clientIP) == 0 {
		clientIP = p.Addr.String()
	}
	clientIP = formatClientIP(clientIP)
	return
}

func GetForwardMetadata(ctx context.Context) (clientIP, userAgent string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if clientIPs := md.Get(forwardContextKeyClientIP.String()); len(clientIPs) > 0 {
			clientIP = clientIPs[0]
		}
		if userAgents := md.Get(forwardContextKeyUserAgent.String()); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}
		return
	}
	return "", ""
}

func SetForwardMetadata(ctx context.Context, clientIP, userAgent string) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Append(forwardContextKeyClientIP.String(), clientIP)
		md.Append(forwardContextKeyUserAgent.String(), userAgent)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
			forwardContextKeyClientIP.String(), clientIP,
			forwardContextKeyUserAgent.String(), userAgent,
		))
	}
	return ctx
}

func formatClientIP(ip string) string {
	fmt.Println(ip)
	twoDots := strings.Count(ip, ":")
	if twoDots > 1 && strings.Contains(ip, "[") { // IPV6
		ip = ip[:strings.LastIndex(ip, ":")]
		if ip[0] == '[' {
			ip = ip[1:]
		}
		if ip[len(ip)-1] == ']' {
			ip = ip[:len(ip)-1]
		}
	} else if twoDots == 1 { // IPV4
		ip = ip[:strings.LastIndex(ip, ":")]
	}
	return ip
}
