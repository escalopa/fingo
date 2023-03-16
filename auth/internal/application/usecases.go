package application

import (
	"context"
	"google.golang.org/grpc/metadata"
	"time"

	"github.com/google/uuid"
	"github.com/lordvidex/errs"
)

type UseCases struct {
	v  Validator
	h  PasswordHasher
	tg TokenGenerator
	ur UserRepository
	sr SessionRepository
	rr RoleRepository
	tr TokenRepository
	mp MessageProducer

	Query
	Command
}

func NewUseCases(opts ...func(*UseCases)) *UseCases {
	u := &UseCases{}
	for _, opt := range opts {
		opt(u)
	}
	u.Query = Query{
		GetUserDevices: NewGetUserDevicesCommand(u.v, u.sr),
	}
	u.Command = Command{
		Signin:     NewSigninCommand(u.v, u.h, u.tg, u.ur, u.sr, u.tr, u.rr, u.mp),
		Signup:     NewSignupCommand(u.v, u.h, u.ur),
		Logout:     NewLogoutCommand(u.v, u.sr, u.rr, u.tr),
		RenewToken: NewRenewTokenCommand(u.v, u.tg, u.sr, u.rr, u.tr),
		CreateRole: NewCreateRoleCommand(u.v, u.rr),
		GrantRole:  NewGrantRoleCommand(u.v, u.rr),
		RevokeRole: NewRevokeRoleCommand(u.v, u.rr),
	}
	return u
}

func WithUserRepository(ur UserRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.ur = ur
	}
}

func WithSessionRepository(sr SessionRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.sr = sr
	}
}

func WithTokenRepository(tr TokenRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.tr = tr
	}
}

func WithRoleRepository(rr RoleRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.rr = rr
	}
}

func WithTokenGenerator(tg TokenGenerator) func(*UseCases) {
	return func(u *UseCases) {
		u.tg = tg
	}
}

func WithPasswordHasher(h PasswordHasher) func(*UseCases) {
	return func(u *UseCases) {
		u.h = h
	}
}

func WithMessageProducer(mp MessageProducer) func(*UseCases) {
	return func(u *UseCases) {
		u.mp = mp
	}
}

func WithValidator(v Validator) func(*UseCases) {
	return func(u *UseCases) {
		u.v = v
	}
}

type Query struct {
	GetUserDevices GetUserDevicesCommand
}

type Command struct {
	Signin SigninCommand
	Signup SignupCommand
	Logout LogoutCommand

	RenewToken RenewTokenCommand

	CreateRole CreateRoleCommand
	GrantRole  GrantRoleCommand
	RevokeRole RevokeRoleCommand
}

func executeWithContextTimeout(ctx context.Context, timeout time.Duration, handler func() error) error {
	// Create a context with a given timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	errChan := make(chan error)
	// Execute logic
	go func() {
		//defer close(errChan)
		// Check if the context has closed already before calling handler
		if ctxWithTimeout.Err() != nil {
			errChan <- ctx.Err()
			return
		}
		errChan <- handler()
	}()
	// Wait for response
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctxWithTimeout.Done():
		return errs.B().Code(errs.DeadlineExceeded).Msg("context timeout").Err()
	}
	return nil
}

func parseUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if meta, ok := metadata.FromIncomingContext(ctx); ok {
		shard := meta.Get("user-id")
		if len(shard) != 1 {
			return uuid.UUID{}, errs.B().Code(errs.Unauthenticated).Msg("user id not passed in headers").Err()
		}
		id := shard[0]
		userID, err := uuid.Parse(id)
		if err != nil {
			return uuid.UUID{}, errs.B(err).Code(errs.Internal).Msg("failed to parse user id from headers").Err()
		}
		return userID, nil
	} else {
		userIDVal, ok := ctx.Value("user-id").(string)
		if !ok {
			return uuid.UUID{}, errs.B().Code(errs.Unauthenticated).Msg("user id not passed in headers").Err()
		}
		userID, err := uuid.Parse(userIDVal)
		if err != nil {
			return uuid.UUID{}, errs.B(err).Code(errs.Internal).Msg("failed to parse user id from headers").Err()
		}
		return userID, nil
	}
}
