package locker

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestConcurrentLocker_Lock(t *testing.T) {
	l := NewLocker(context.TODO(), time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unlock := l.Lock(context.TODO(), "x", "y", "z")
			time.Sleep(time.Second * 2) // let the other goroutines try to lock the same keys
			defer unlock()
		}()
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
	type args struct {
		ctx  context.Context
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
				ctx: context.Background(),
				x:   "x",
				y:   []any{"y", "z"},
				f:   time.Minute,
			},
		},
		{
			name: "test cleanup",
			args: args{
				ctx:  context.Background(),
				x:    "x",
				y:    []any{"y", "z"},
				f:    time.Second,
				wait: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocker(tt.args.ctx, tt.args.f)
			unlock := l.Lock(tt.args.ctx, tt.args.x, tt.args.y...)
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
