package application

import (
	"context"
	"github.com/escalopa/fingo/pkg/pkgCore"
)

// extractMetadataFromContext extracts client ip and user agent from context metadata
func extractMetadataFromContext(ctx context.Context) (clientIP, userAgent string) {
	// Extract client IP from context
	if ctx.Value(pkgCore.ContextKeyClientIP) != nil {
		varClientIP, ok := ctx.Value(pkgCore.ContextKeyClientIP).(string)
		if ok {
			clientIP = varClientIP
		}
	}
	// Extract user agent from context
	if ctx.Value(pkgCore.ContextKeyUserAgent) != nil {
		varUserAgent, ok := ctx.Value(pkgCore.ContextKeyUserAgent).(string)
		if ok {
			userAgent = varUserAgent
		}
	}
	return
}
