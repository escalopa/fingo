package application

import (
	"context"
	"log"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"

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
	Execute(ctx context.Context, params RenewTokenParams) (RenewTokenResponse, error)
}

// RenewTokenCommandImpl is the implementation of the RenewTokenCommand
type RenewTokenCommandImpl struct {
	v  Validator
	tg TokenGenerator
	sr SessionRepository
	tr TokenRepository
}

// Execute executes the RenewTokenCommand with the given params
func (c *RenewTokenCommandImpl) Execute(ctx context.Context, params RenewTokenParams) (RenewTokenResponse, error) {
	var response RenewTokenResponse
	err := contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "SignupCommand.Execute")
		defer span.End()
		// Validate params
		if err := c.v.Validate(ctx, params); err != nil {
			return err
		}
		// Decrypt token to get sessionID
		payload, err := c.tg.DecryptToken(ctx, params.RefreshToken)
		if err != nil {
			return err
		}
		// Get caller ID from context
		callerID, err := contextutils.GetUserID(ctx)
		if err != nil {
			return err
		}
		if callerID != payload.UserID {
			return errs.B().Code(errs.Unauthenticated).Msg("not token owner").Err()
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
			return errs.B().Code(errs.NotFound).Msg("refresh token doesn't match the one stored in session database").Err()
		}
		// Generate a new access & refresh tokens
		accessToken, err := c.tg.GenerateAccessToken(ctx, core.GenerateTokenParam{
			UserID:    payload.UserID,
			SessionID: payload.SessionID,
			ClientIP:  payload.ClientIP,
			UserAgent: payload.UserAgent,
		})
		if err != nil {
			return errs.B(err).Code(errs.Internal).Msg("failed to create access token").Err()
		}
		refreshToken, err := c.tg.GenerateRefreshToken(ctx, core.GenerateTokenParam{
			UserID:    payload.UserID,
			SessionID: payload.SessionID,
			ClientIP:  payload.ClientIP,
			UserAgent: payload.UserAgent,
		})
		if err != nil {
			return errs.B(err).Code(errs.Internal).Msg("failed to create refresh token").Err()
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
		payload, err = c.tg.DecryptToken(ctx, accessToken)
		if err != nil {
			return err
		}
		// Store access token in cache repository
		err = c.tr.Store(ctx, accessToken, payload)
		if err != nil {
			return err
		}
		response = RenewTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		go func() {
			// Disable old access token in cache repository
			err = c.tr.Delete(ctx, session.AccessToken)
			if err != nil {
				log.Printf("failed to remove old access token, :%s", err)
			}
		}()
		return nil
	})
	return response, err
}

// NewRenewTokenCommand creates a new RenewTokenCommand with the passed dependencies
func NewRenewTokenCommand(v Validator, tg TokenGenerator, sr SessionRepository, tr TokenRepository) RenewTokenCommand {
	return &RenewTokenCommandImpl{v: v, tg: tg, sr: sr, tr: tr}
}
