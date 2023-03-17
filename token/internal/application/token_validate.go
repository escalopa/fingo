package application

import (
	"context"
	"github.com/lordvidex/errs"
	"time"
)

// TokenValidateParams is the params for the TokenValidateCommand
// AccessToken is the access token to validate (required)
type TokenValidateParams struct {
	AccessToken string `validate:"required"`
}

type TokenValidateCommand interface {
	Execute(ctx context.Context, params TokenValidateParams) error
}

type TokenValidateCommandImpl struct {
	v  Validator       // Validator is a custom interface that validates the params
	tr TokenRepository // TokenRepository is a custom interface that gets the token payload from the database
}

// Execute executes the TokenValidateCommand with the given params
func (c *TokenValidateCommandImpl) Execute(ctx context.Context, params TokenValidateParams) error {
	return executeWithContextTimeout(ctx, 10*time.Second, func() error {
		if err := c.v.Validate(params); err != nil {
			return err
		}
		// Get the token payload from the cache
		payload, err := c.tr.GetTokenPayload(ctx, params.AccessToken)
		if err != nil {
			return err
		}
		// Check if the token has expired
		if time.Now().After(payload.ExpiresAt) {
			return errs.B().Msg("access token has expired").Err()
		}
		clientIP, userAgent := extractMetadataFromContext(ctx)
		// Check if the client ip is the same
		if payload.ClientIP != clientIP {
			return errs.B().Msg("client ip mismatch, possible ip spoofing").Err()
		}
		// Check if the user agent is the same
		if payload.UserAgent != userAgent {
			return errs.B().Msg("user agent mismatch, possible user agent spoofing").Err()
		}
		return nil
	})
}

// NewTokenValidateCommand creates a new TokenValidateCommand
func NewTokenValidateCommand(v Validator, tr TokenRepository) *TokenValidateCommandImpl {
	return &TokenValidateCommandImpl{v: v, tr: tr}
}
