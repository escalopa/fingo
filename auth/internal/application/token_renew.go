package application

import (
	"context"

	"github.com/escalopa/fingo/auth/internal/adapters/token"
	"github.com/lordvidex/errs"
)

type RenewTokenParams struct {
	RefreshToken string `validate:"required"`
}

type RenewTokenCommand interface {
	Execute(ctx context.Context, params RenewTokenParams) (string, error)
}

type RenewTokenCommandImpl struct {
	v  Validator
	tg TokenGenerator
	sr SessionRepository
}

func (v *RenewTokenCommandImpl) Execute(ctx context.Context, params RenewTokenParams) (string, error) {
	// Decrypt token to get sessionID
	user, sessionID, err := v.tg.VerifyToken(params.RefreshToken)
	if err != nil {
		return "", err
	}
	// Get user session from database to check if session is blocked
	session, err := v.sr.GetSessionByID(ctx, sessionID)
	if err != nil {
		return "", err
	}
	// Validate user session
	err = validateSession(session)
	if err != nil {
		return "", err
	}
	// Validate refresh token
	if session.RefreshToken != params.RefreshToken {
		return "", errs.B().Msg("refresh token doesn't match the one stored in session database").Err()
	}
	// Generate a new access token
	accessToken, err := v.tg.GenerateAccessToken(token.GenerateTokenParam{
		User:      user,
		SessionID: sessionID,
	})
	if err != nil {
		return "", errs.B(err).Msg("failed to create access token").Err()
	}
	return accessToken, nil
}

func NewRenewTokenCommand(v Validator, tg TokenGenerator, sr SessionRepository) RenewTokenCommand {
	return &RenewTokenCommandImpl{v: v, tg: tg, sr: sr}
}
