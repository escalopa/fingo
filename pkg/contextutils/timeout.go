package contextutils

import (
	"context"
	"time"

	"github.com/lordvidex/errs"
)

func ExecuteWithContextTimeout(ctx context.Context, timeout time.Duration, handler func() error) error {
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
