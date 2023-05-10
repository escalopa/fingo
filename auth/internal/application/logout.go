package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/escalopa/fingo/pkg/tracer"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type LogoutParams struct {
	SessionID string `validate:"required,uuid"`
}

type LogoutCommand interface {
	Execute(ctx context.Context, params LogoutParams) error
}

type LogoutCommandImpl struct {
	v  Validator
	sr SessionRepository
	tr TokenRepository
}

func (c *LogoutCommandImpl) Execute(ctx context.Context, params LogoutParams) error {
	return contextutils.ExecuteWithContextTimeout(ctx, 5*time.Second, func() error {
		ctx, span := tracer.Tracer().Start(ctx, "SignupCommand.Execute")
		defer span.End()
		if err := c.v.Validate(ctx, params); err != nil {
			return err
		}
		// Read user id from context
		callerID, err := contextutils.GetUserID(ctx)
		if err != nil {
			return err
		}
		// Parse sessionUUID
		sessionUUID, _ := uuid.Parse(params.SessionID)
		// Get session form DB
		var session core.Session
		session, err = c.sr.GetSessionByID(ctx, sessionUUID)
		if err != nil {
			return err
		}
		// Check session owner is the caller
		if callerID != session.UserID {
			return errs.B().Code(errs.Forbidden).Msg("not session owner").Err()
		}
		// Invalidate the access token from cache storage
		err = c.tr.Delete(ctx, session.AccessToken)
		if err != nil {
			// If token is already deleted (expired) then skip this error
			if errErrs, ok := err.(*errs.Error); ok {
				if errErrs.Code != errs.NotFound {
					return err
				}
			} else {
				return err
			}
		}
		// Delete Session
		err = c.sr.DeleteSessionByID(ctx, sessionUUID)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewLogoutCommand(v Validator, sr SessionRepository, tr TokenRepository) LogoutCommand {
	return &LogoutCommandImpl{v: v, sr: sr, tr: tr}
}
