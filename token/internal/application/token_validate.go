package application

import (
	"context"
	"fmt"
	"time"

	oteltracer "github.com/escalopa/fingo/token/internal/adapters/tracer"

	"github.com/escalopa/fingo/pkg/contextutils"
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
	err := contextutils.ExecuteWithContextTimeout(ctx, 10*time.Second, func() error {
		ctx, span := oteltracer.Tracer().Start(ctx, "TokenValidateCommandImpl.Execute")
		defer span.End()
		fmt.Println("Executing TokenValidateCommandImpl.Execute", params)
		if err := c.v.Validate(ctx, params); err != nil {
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
		clientIP, userAgent := contextutils.GetForwardMetadata(ctx)
		// Check if the client ip is the same
		if payload.ClientIP != clientIP {
			return errs.B().Msg("client ip mismatch, possible ip spoofing", payload.ClientIP, "|", clientIP).Err()
		}
		// Check if the user agent is the same
		if payload.UserAgent != userAgent {
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
