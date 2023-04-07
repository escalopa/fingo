package locker

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestConcurrentLocker_Lock(t *testing.T) {
	ctx := context.TODO()
	l := NewLocker(ctx, time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			unlock := l.Lock(ctx, "x", "y", "z")
			time.Sleep(time.Second * 2) // let the other goroutines try to lock the same keys
			defer unlock()
		}(i)
	}

	time.Sleep(time.Second * 2) // wait for the cleanup to run
	for _, key := range []any{"x", "y", "z"} {
		_, ok := l.mx.Load(key)
		if !ok {
			t.Errorf("lock not found")
		}
	}
	wg.Wait()
}

func TestLocker_Lock(t *testing.T) {
	ctx := context.TODO()
	type args struct {
		x    any
		y    []any
		f    time.Duration
		wait bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "lock",
			args: args{
				x: "x",
				y: []any{"y", "z"},
				f: time.Minute,
			},
		},
		{
			name: "test cleanup",
			args: args{
				x:    "x",
				y:    []any{"y", "z"},
				f:    time.Second,
				wait: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocker(context.TODO(), tt.args.f)
			unlock := l.Lock(ctx, tt.args.x, tt.args.y...)
			unlock()
			for _, key := range append([]any{tt.args.x}, tt.args.y...) {
				_, ok := l.mx.Load(key)
				if !ok {
					t.Errorf("lock not found")
				}
			}
			if tt.args.wait {
				time.Sleep(tt.args.f * 2)
				for _, key := range append([]any{tt.args.x}, tt.args.y...) {
					_, ok := l.mx.Load(key)
					if ok {
						t.Errorf("lock not cleaned up")
					}
				}
			}
		})
	}
}
