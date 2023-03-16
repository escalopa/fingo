package application

import (
	"context"
	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
	"reflect"
	"time"
)

type LogoutParams struct {
	SessionID string `validate:"required,uuid"`
}

type LogoutCommand interface {
	Execute(ctx context.Context, params LogoutParams) error
}

type LogoutCommandImpl struct {
	v  Validator
	rr RoleRepository
	sr SessionRepository
	tr TokenRepository
}

func (c *LogoutCommandImpl) Execute(ctx context.Context, params LogoutParams) error {
	return executeWithContextTimeout(ctx, 5*time.Second, func() error {
		if err := c.v.Validate(params); err != nil {
			return err
		}
		// Read user id from context
		callerID, err := parseUserIDFromContext(ctx)
		if err != nil {
			return err
		}
		// Parse sessionID
		sessionID, err := uuid.Parse(params.SessionID)
		if err != nil {
			return errs.B(err).Code(errs.InvalidArgument).Msg("invalid session id").Err()
		}
		// Get session form DB
		var session core.Session
		session, err = c.sr.GetSessionByID(ctx, sessionID)
		if err != nil {
			return err
		}
		// Check if the user is admin
		isAdmin, err := c.rr.HasPrivillage(ctx, core.HasPrivillageParams{
			UserID:   callerID,
			RoleName: "admin",
		})
		if err != nil {
			return err
		}
		// Check if the user has enough rights to invalidate the session
		if !isAdmin && !reflect.DeepEqual(session.UserID, callerID) {
			return errs.B().Code(errs.Forbidden).Msg("not enough privileges to invalidate session").Err()
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
		err = c.sr.DeleteSessionByID(ctx, sessionID)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewLogoutCommand(v Validator, sr SessionRepository, rr RoleRepository, tr TokenRepository) LogoutCommand {
	return &LogoutCommandImpl{v: v, sr: sr, rr: rr, tr: tr}
}
