package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/lordvidex/errs"
)

// DecryptTokenParams is the input for DecryptTokenCommand
type DecryptTokenParams struct {
	AccessToken string `validate:"required"`
}

// DecryptTokenCommand is the interface for decrypting a token
type DecryptTokenCommand interface {
	Execute(ctx context.Context, params DecryptTokenParams) (string, error)
}

// DecryptTokenCommandImpl implements DecryptTokenCommand
type DecryptTokenCommandImpl struct {
	v  Validator
	tg TokenGenerator
	sr SessionRepository
}

// Execute decrypts the token and returns the user ID
func (v *DecryptTokenCommandImpl) Execute(ctx context.Context, params DecryptTokenParams) (string, error) {
	if err := v.v.Validate(params); err != nil {
		return "", err
	}
	// Verify token & get user
	user, sessionID, err := v.tg.VerifyToken(params.AccessToken)
	if err != nil {
		return "", err
	}
	// Get session by ID from database
	session, err := v.sr.GetSessionByID(ctx, sessionID)
	if err != nil {
		return "", errs.B(err).Msg("user session not found").Err()
	}
	// Validate session
	err = validateSession(session)
	if err != nil {
		return "", err
	}
	return user.ID.String(), nil
}

// NewDecryptTokenCommand returns a new instance of DecryptTokenCommand
func NewDecryptTokenCommand(v Validator, tg TokenGenerator, sr SessionRepository) DecryptTokenCommand {
	return &DecryptTokenCommandImpl{v: v, tg: tg, sr: sr}
}

// validateSession checks if the session is blocked or expired
func validateSession(session core.Session) error {
	// Check is session is blocked
	if session.IsBlocked {
		return errs.B().Msg("user session is blocked").Err()
	}
	// Check is session has expired
	if time.Now().After(session.ExpiresAt) {
		return errs.B().Msg("user session has expired").Err()
	}
	return nil
}
