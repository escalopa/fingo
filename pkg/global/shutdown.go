package global

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CatchSignal() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	return sig
}

func Shutdown(ctx context.Context, timeout time.Duration, gracefulStop func(), forceStop func()) {
	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), timeout)
	log.Println("Server shutdown initiated, Press Ctrl+C to force")
	go func() {
		defer stopCancel()
		gracefulStop()
	}()

	select {
	case <-stopCtx.Done():
		switch stopCtx.Err() {
		case context.DeadlineExceeded:
			log.Println("Server shutdown timeout")
		case context.Canceled:
			log.Println("Server shutdown done")
		}
	case <-CatchSignal():
		fmt.Println("Shutdown force")
		forceStop()
	}
}
