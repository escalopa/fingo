package application

import (
	"context"
	"reflect"
	"time"

	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/lordvidex/errs"
)

type VerifyTokenParams struct {
	AccessToken string `validate:"required"`
}

type VerifyTokenCommand interface {
	Execute(ctx context.Context, params VerifyTokenParams) (core.User, error)
}

type VerifyTokenCommandImpl struct {
	v  Validator
	tg TokenGenerator
	sr SessionRepository
}

func (v *VerifyTokenCommandImpl) Execute(ctx context.Context, params VerifyTokenParams) (core.User, error) {
	if err := v.v.Validate(params); err != nil {
		return core.User{}, err
	}
	// TODO: Check client-ip & user-agent matches the ones stored on the database,
	// 	 if the passed ip & user-agent doesn't match, return unrecognized device & send to the user
	//	 email about the call signin attempt

	// Verify token & get user
	user, sessionID, err := v.tg.VerifyToken(params.AccessToken)
	if err != nil {
		return core.User{}, err
	}
	// Check empty user token check
	if reflect.DeepEqual(user, core.User{}) {
		return core.User{}, errs.B().Code(errs.Unauthenticated).Msg("invalid token, token not assigned").Err()
	}
	// Get session by ID from database
	session, err := v.sr.GetSessionByID(ctx, sessionID)
	if err != nil {
		return core.User{}, errs.B(err).Msg("user session not found").Err()
	}
	// Validate session
	err = validateSession(session)
	if err != nil {
		return core.User{}, err
	}
	if !user.IsVerified {
		return core.User{}, errs.B().Code(errs.Unauthenticated).Msg("invalid token, user not verified").Err()
	}
	return user, nil
}

func NewVerifyTokenCommand(v Validator, tg TokenGenerator, sr SessionRepository) VerifyTokenCommand {
	return &VerifyTokenCommandImpl{v: v, tg: tg, sr: sr}
}

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
