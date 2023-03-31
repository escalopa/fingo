package application

import (
	"context"
	"time"

	"github.com/escalopa/fingo/pkg/contextutils"

	"github.com/escalopa/fingo/auth/internal/core"
)

// GetUserDevicesParams contains the parameters for the GetUserDevicesCommand
type GetUserDevicesParams struct{}

// GetUserDevicesCommand is the interface for the GetUserDevicesCommandImpl
type GetUserDevicesCommand interface {
	Execute(ctx context.Context, params GetUserDevicesParams) ([]core.Session, error)
}

// GetUserDevicesCommandImpl is the implementation of the GetUserDevicesCommand
type GetUserDevicesCommandImpl struct {
	v  Validator
	sr SessionRepository
}

// Execute executes the GetUserDevicesCommand with the given parameters
func (c *GetUserDevicesCommandImpl) Execute(ctx context.Context, params GetUserDevicesParams) ([]core.Session, error) {
	var response []core.Session
	err := executeWithContextTimeout(ctx, 5*time.Second, func() error {
		// Validate function can be removed since the params are empty
		// But for design patterns & logic it won't be removed
		err := c.v.Validate(params)
		if err != nil {
			return err
		}
		// Parse userID from context
		userID, err := contextutils.GetUserID(ctx)
		if err != nil {
			return err
		}
		// Get user session from db
		sessions, err := c.sr.GetUserSessions(ctx, userID)
		if err != nil {
			return err
		}
		response = sessions
		return nil
	})
	return response, err
}

// NewGetUserDevicesCommand returns a new GetUserDevicesCommand with the passed dependencies
func NewGetUserDevicesCommand(
	v Validator,
	sr SessionRepository,
) GetUserDevicesCommand {
	return &GetUserDevicesCommandImpl{v: v, sr: sr}
}
