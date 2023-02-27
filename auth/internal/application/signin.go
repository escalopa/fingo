package application

import (
	"context"
	"github.com/escalopa/gochat/auth/internal/adapters/token"
	"github.com/escalopa/gochat/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

// ---------------------- Signin ---------------------- //

type SigninParams struct {
	Email    string            `validate:"required,email"`
	Password string            `validate:"required,min=8"`
	MetaData map[string]string `validate:"required"` // Client-IP, User-Agent
}

type SigninResponse struct {
	AccessToken  string
	RefreshToken string
}

type SigninCommand interface {
	Execute(ctx context.Context, params SigninParams) (SigninResponse, error)
}

type SigninCommandImpl struct {
	v  Validator
	h  PasswordHasher
	ur UserRepository
	sr SessionRepository
	tg TokenGenerator
}

func (s *SigninCommandImpl) Execute(ctx context.Context, params SigninParams) (SigninResponse, error) {
	if err := s.v.Validate(params); err != nil {
		return SigninResponse{}, err
	}
	// Get user from database
	user, err := s.ur.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return SigninResponse{}, err
	}
	// Compare password
	if !s.h.Compare(user.Password, params.Password) {
		return SigninResponse{}, errs.B().Code(errs.InvalidArgument).Msg("password is incorrect").Err()
	}
	// if user is not verified, return error
	if !user.IsVerified {
		return SigninResponse{}, errs.B().Code(errs.Unauthenticated).Msg("user is not verified").Err()
	}
	// Create new sessionID
	sessionID := uuid.New()
	// Generate access token
	accessToken, err := s.tg.GenerateAccessToken(token.GenerateTokenParam{
		User:      user,
		SessionID: sessionID,
	})
	if err != nil {
		return SigninResponse{}, err
	}
	// Generate refresh token
	refreshToken, err := s.tg.GenerateRefreshToken(token.GenerateTokenParam{
		User:      user,
		SessionID: sessionID,
	})
	if err != nil {
		return SigninResponse{}, err
	}
	// Create a new session for user
	err = s.sr.CreateSession(ctx, core.CreateSessionParams{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    params.MetaData["USER-AGENT"],
		ClientIp:     params.MetaData["CLIENT-IP"],
	})
	if err != nil {
		return SigninResponse{}, err
	}
	// Return response
	return SigninResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewSigninCommand(
	v Validator,
	h PasswordHasher,
	tg TokenGenerator,
	ur UserRepository,
	sr SessionRepository) SigninCommand {
	return &SigninCommandImpl{v: v, h: h, tg: tg, ur: ur, sr: sr}
}
