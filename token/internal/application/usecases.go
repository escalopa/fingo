package application

import (
	"context"
	"github.com/lordvidex/errs"
	"time"
)

type UseCases struct {
	v  Validator
	tr TokenRepository

	Command
}

func NewUseCases(opts ...func(*UseCases)) *UseCases {
	u := &UseCases{}
	for _, opt := range opts {
		opt(u)
	}
	u.Command = Command{
		TokenValidate: NewTokenValidateCommand(u.v, u.tr),
	}
	return u
}

func WithTokenRepository(tr TokenRepository) func(*UseCases) {
	return func(u *UseCases) {
		u.tr = tr
	}
}

func WithValidator(v Validator) func(*UseCases) {
	return func(u *UseCases) {
		u.v = v
	}
}

func executeWithContextTimeout(ctx context.Context, timeout time.Duration, handler func() error) error {
	// Create a context with a given timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	errChan := make(chan error)
	// Execute logic
	go func() {
		defer close(errChan)
		if err := ctx.Err(); err != nil {
			return // err is sent to errChan by the select below
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

type Command struct {
	TokenValidate TokenValidateCommand
}
