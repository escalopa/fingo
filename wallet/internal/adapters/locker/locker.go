package locker

import (
	"context"
	"sort"
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

// values passed to Lock must be of the same type [either all strings or all ints]
func (l *Locker) Lock(ctx context.Context, x any, y ...any) func() {
	if ctx.Err() != nil {
		return func() { /* Return empty function since context is empty */ }
	}
	locks := make([]*sync.Mutex, len(y)+1)
	// sort the keys to avoid deadlocks
	arr := orderedLocks(append([]any{x}, y...))
	sort.Sort(arr)
	for i, key := range arr {
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

type orderedLocks []any

func (o orderedLocks) Len() int {
	return len(o)
}

func (o orderedLocks) Less(i int, j int) bool {
	switch o[i].(type) {
	case string:
		return o[i].(string) < o[j].(string)
	case int:
		return o[i].(int) < o[j].(int)
	default:
		panic("Lock key must either be strings or ints")
	}
}

func (o orderedLocks) Swap(i int, j int) {
	o[i], o[j] = o[j], o[i]
}
