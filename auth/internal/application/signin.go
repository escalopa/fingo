package application

import (
	"context"
	"log"
	"time"

	"github.com/escalopa/fingo/auth/internal/core"
	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

// SigninParams contains the parameters for the SigninCommand
type SigninParams struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
	ClientIP  string `validate:"required,ip"`
	UserAgent string `validate:"required,min=1"`
}

// SigninResponse contains the response for the SigninCommand
type SigninResponse struct {
	AccessToken  string
	RefreshToken string
}

// SigninCommand is the interface for the SigninCommandImpl
type SigninCommand interface {
	Execute(ctx context.Context, params SigninParams) (SigninResponse, error)
}

// SigninCommandImpl is the implementation of the SigninCommand
type SigninCommandImpl struct {
	v  Validator
	h  PasswordHasher
	ur UserRepository
	sr SessionRepository
	tr TokenRepository
	rr RoleRepository
	tg TokenGenerator
	mp MessageProducer
}

// Execute executes the SigninCommand with the given parameters
func (c *SigninCommandImpl) Execute(ctx context.Context, params SigninParams) (SigninResponse, error) {
	var response SigninResponse
	err := executeWithContextTimeout(ctx, 5*time.Second, func() error {
		if err := c.v.Validate(params); err != nil {
			return err
		}
		// Get user from database
		user, err := c.ur.GetUserByEmail(ctx, params.Email)
		if err != nil {
			return err
		}
		// Compare password
		if !c.h.Compare(user.HashedPassword, params.Password) {
			return errs.B().Code(errs.InvalidArgument).Msg("password is incorrect").Err()
		}
		// Create new sessionID
		sessionID := uuid.New()
		// Generate access token
		accessToken, err := c.tg.GenerateAccessToken(core.GenerateTokenParam{
			UserID:    user.ID,
			ClientIP:  params.ClientIP,
			UserAgent: params.UserAgent,
			SessionID: sessionID,
			Roles:     []string{},
		})
		if err != nil {
			return err
		}
		// Generate refresh token
		refreshToken, err := c.tg.GenerateRefreshToken(core.GenerateTokenParam{
			UserID:    user.ID,
			ClientIP:  params.ClientIP,
			UserAgent: params.UserAgent,
			SessionID: sessionID,
			Roles:     []string{},
		})
		if err != nil {
			return err
		}
		// Get user roles
		roles, err := c.rr.GetUserRoles(ctx, user.ID)
		if err != nil {
			return err
		}
		roleNames := make([]string, len(roles))
		for i := range roles {
			roleNames[i] = roles[i].Name
		}
		// Get token payload after encryption
		payload, err := c.tg.DecryptToken(accessToken)
		if err != nil {
			return err
		}
		// Store access token in cache repository
		err = c.tr.Store(ctx, accessToken, payload)
		if err != nil {
			return err
		}
		// Create a new session for user
		err = c.sr.CreateSession(ctx, core.CreateSessionParams{
			ID:     sessionID,
			UserID: user.ID,
			UserDevice: core.UserDevice{
				UserAgent: params.UserAgent,
				ClientIP:  params.ClientIP,
			},
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}
		// Publish message to queue to notify user about new session creation
		err = c.mp.SendNewSignInSessionMessage(ctx, core.SendNewSignInSessionParams{
			Name:  user.FirstName,
			Email: user.Email,
			UserDevice: core.UserDevice{
				ClientIP:  params.ClientIP,
				UserAgent: params.UserAgent,
			},
		})
		if err != nil {
			log.Printf("failed to send message about new session creation, email: %s, client-ip: %s, user-agent: %s, err: %s",
				params.Email, params.ClientIP, params.UserAgent, err)
		}
		response = SigninResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		return nil
	})
	return response, err
}

// NewSigninCommand returns a new SigninCommand with the passed dependencies
func NewSigninCommand(
	v Validator,
	h PasswordHasher,
	tg TokenGenerator,
	ur UserRepository,
	sr SessionRepository,
	mp MessageProducer,
) SigninCommand {
	return &SigninCommandImpl{v: v, h: h, tg: tg, ur: ur, sr: sr, mp: mp}
}
