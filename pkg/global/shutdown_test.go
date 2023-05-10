package global

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCatchSignal(t *testing.T) {
	tests := []struct {
		name string
		sig  syscall.Signal
	}{
		{
			name: "success on SIGTERM",
			sig:  syscall.SIGTERM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				time.Sleep(1 * time.Millisecond)
				syscall.Kill(syscall.Getpid(), tt.sig)
			}()
			got := <-CatchSignal()
			require.Equal(t, tt.sig, got)
		})
	}
}

func TestShutdown(t *testing.T) {
	type args struct {
		ctx          func() context.Context
		timeout      time.Duration
		gracefulStop func()
		forceStop    func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success with graceful stop",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(1 * time.Millisecond)
						cancel()
					}()
					return ctx
				},
				timeout: 1 * time.Second,
				gracefulStop: func() {
					time.Sleep(1 * time.Millisecond)
				},
				forceStop: func() {
					time.Sleep(1 * time.Millisecond)
				},
			},
		},
		{
			name: "success with deadline exceeded",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(1 * time.Millisecond)
						cancel()
					}()
					return ctx
				},
				timeout: 1 * time.Millisecond,
				gracefulStop: func() {
					time.Sleep(1 * time.Second)
				},
				forceStop: func() {
					time.Sleep(1 * time.Second)
				},
			},
		},
		{
			name: "success with forced stop",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(1 * time.Millisecond)
						cancel()
					}()
					go func() {
						time.Sleep(10 * time.Millisecond)
						syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					}()
					return ctx
				},
				timeout: 1 * time.Second,
				gracefulStop: func() {
					time.Sleep(1 * time.Second)
				},
				forceStop: func() {
					time.Sleep(1 * time.Millisecond)
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Shutdown(tt.args.ctx(), tt.args.timeout, tt.args.gracefulStop, tt.args.forceStop)
		})
	}
}
