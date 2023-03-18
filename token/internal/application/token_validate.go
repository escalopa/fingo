package application

import (
	"context"
	"fmt"
	"github.com/escalopa/fingo/pkg/pkgCore"
	"time"

	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

// TokenValidateParams is the params for the TokenValidateCommand
// AccessToken is the access token to validate (required)
type TokenValidateParams struct {
	AccessToken string `validate:"required"`
}

type TokenValidateCommand interface {
	Execute(ctx context.Context, params TokenValidateParams) (uuid.UUID, error)
}

type TokenValidateCommandImpl struct {
	v  Validator       // Validator is a custom interface that validates the params
	tr TokenRepository // TokenRepository is a custom interface that gets the token payload from the database
}

// Execute executes the TokenValidateCommand with the given params
func (c *TokenValidateCommandImpl) Execute(ctx context.Context, params TokenValidateParams) (uuid.UUID, error) {
	var id uuid.UUID
	err := executeWithContextTimeout(ctx, 10*time.Second, func() error {
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
		clientIP, userAgent := pkgCore.GetMDFromContext(ctx)
		fmt.Println(clientIP, userAgent)
		fmt.Println(payload.ClientIP, payload.UserAgent)
		// Check if the client ip is the same
		if payload.ClientIP != clientIP {
			fmt.Println(payload.ClientIP, clientIP)
			return errs.B().Msg("client ip mismatch, possible ip spoofing", payload.ClientIP, "|", clientIP).Err()
		}
		// Check if the user agent is the same
		if payload.UserAgent != userAgent {
			fmt.Println(payload.UserAgent, userAgent)
			return errs.B().Msg("user agent mismatch, possible user agent spoofing", payload.UserAgent, "|", userAgent).Err()
		}
		id = payload.UserID
		return nil
	})
	return id, err
}

// NewTokenValidateCommand creates a new TokenValidateCommand
func NewTokenValidateCommand(v Validator, tr TokenRepository) *TokenValidateCommandImpl {
	return &TokenValidateCommandImpl{v: v, tr: tr}
}
