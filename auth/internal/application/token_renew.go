package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/lordvidex/errs"
)

// RenewTokenParams is the input for the RenewTokenCommand
type RenewTokenParams struct {
	RefreshToken string `validate:"required"`
}

// RenewTokenResponse is the output for the RenewTokenCommand
type RenewTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

// RenewTokenCommand is the interface for the RenewTokenCommand
type RenewTokenCommand interface {
	Execute(ctx context.Context, params RenewTokenParams) (*RenewTokenResponse, error)
}

// RenewTokenCommandImpl is the implementation of the RenewTokenCommand
type RenewTokenCommandImpl struct {
	v  Validator
	tg TokenGenerator
	sr SessionRepository
	rr RoleRepository
	tr TokenRepository
}

// Execute executes the RenewTokenCommand with the given params
func (c *RenewTokenCommandImpl) Execute(ctx context.Context, params RenewTokenParams) (*RenewTokenResponse, error) {
	var response RenewTokenResponse
	err := executeWithContextTimeout(ctx, 5*time.Second, func() error {
		// Decrypt token to get sessionID
		payload, err := c.tg.DecryptToken(params.RefreshToken)
		if err != nil {
			return err
		}
		// Check if session has expired
		if time.Now().After(payload.ExpiresAt) {
			return errs.B().Msg("user session has expired").Err()
		}
		// Get user's session from database to check if session is blocked
		session, err := c.sr.GetSessionByID(ctx, payload.SessionID)
		if err != nil {
			return err
		}
		// Validate refresh token
		if session.RefreshToken != params.RefreshToken {
			return errs.B().Msg("refresh token doesn't match the one stored in session database").Err()
		}
		// Generate a new access & refresh tokens
		accessToken, err := c.tg.GenerateAccessToken(core.GenerateTokenParam{
			UserID:    payload.UserID,
			SessionID: payload.SessionID,
			Roles:     payload.Roles,
		})
		refreshToken, err := c.tg.GenerateRefreshToken(core.GenerateTokenParam{
			UserID:    payload.UserID,
			SessionID: payload.SessionID,
			Roles:     payload.Roles,
		})
		if err != nil {
			return errs.B(err).Code(errs.Internal).Msg("failed to create access token").Err()
		}
		// Update session in database
		err = c.sr.UpdateSessionTokens(ctx, core.UpdateSessionTokenParams{
			ID:           payload.SessionID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		// Get token payload after encryption to save the new values
		payload, err = c.tg.DecryptToken(accessToken)
		if err != nil {
			return err
		}
		// Store access token in cache repository
		err = c.tr.Store(ctx, accessToken, payload)
		response = RenewTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		return nil
	})
	return &response, err
}

// NewRenewTokenCommand creates a new RenewTokenCommand with the passed dependencies
func NewRenewTokenCommand(v Validator, tg TokenGenerator, sr SessionRepository, rr RoleRepository, tr TokenRepository) RenewTokenCommand {
	return &RenewTokenCommandImpl{v: v, tg: tg, sr: sr, rr: rr, tr: tr}
}
