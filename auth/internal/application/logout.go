package application

import (
	"context"
	"github.com/lordvidex/errs"
	"time"
)

type LogoutParams struct {
	AccessToken string `validate:"required,email"`
}

type LogoutCommand interface {
	Execute(ctx context.Context, params LogoutParams) error
}

type LogoutCommandImpl struct {
	v  Validator
	tg TokenGenerator
	ur UserRepository
	sr SessionRepository
}

func (l *LogoutCommandImpl) Execute(ctx context.Context, params LogoutParams) error {
	// Get sessionID from token
	_, sessionID, err := l.tg.VerifyToken(params.AccessToken)
	if err != nil {
		return nil
	}
	// Get session from database by sessionID
	session, err := l.sr.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	// Validate that session hasn't been blocked or expired
	if session.IsBlocked || time.Now().After(session.ExpiresAt) {
		return errs.B().Code(errs.Unauthenticated).Msg("session is blocked or have expired").Err()
	}
	// Delete Session
	err = l.sr.DeleteSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func NewLogoutCommand(v Validator, tg TokenGenerator, ur UserRepository, sr SessionRepository) LogoutCommand {
	return &LogoutCommandImpl{v: v, tg: tg, ur: ur, sr: sr}
}
