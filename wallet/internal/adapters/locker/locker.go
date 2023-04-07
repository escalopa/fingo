package locker

import (
	"context"
	"sync"
	"time"
)

type Locker struct {
	mx sync.Map
}

func NewLocker(ctx context.Context, frequency time.Duration) *Locker {
	l := &Locker{}
	go l.cleanup(ctx, frequency)
	return l
}

func (l *Locker) cleanup(ctx context.Context, frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	for {
		select {
		case <-ticker.C:
			l.mx.Range(func(key, value interface{}) bool {
				lock, ok := value.(*sync.Mutex)
				if !ok {
					return true
				}
				if lock.TryLock() {
					l.mx.Delete(key)
					value.(*sync.Mutex).Unlock()
				}
				return true
			})
		case <-ctx.Done():
		}
	}
}

func (l *Locker) Lock(ctx context.Context, x any, y ...any) func() {
	if ctx.Err() != nil {
		return func() {}
	}
	locks := make([]*sync.Mutex, len(y)+1)
	for i, key := range append([]any{x}, y...) {
		lock, _ := l.mx.LoadOrStore(key, &sync.Mutex{})
		lock.(*sync.Mutex).Lock()
		locks[i] = lock.(*sync.Mutex)
	}
	return func() {
		for _, lock := range locks {
			lock.Unlock()
		}
	}
}
