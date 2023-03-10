package redis

import (
	"context"
	"github.com/escalopa/fingo/utils/testcontainer"
	"log"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
)

var (
	testRedis *redis.Client
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	client, terminate, err := testcontainer.StartRedisContainer(ctx)
	if err != nil {
		log.Fatalf("failed to start redis container for tests, err: %s", err)
	}
	testRedis = client
	m.Run()
	defer func() {
		if err = terminate(); err != nil {
			log.Fatalf("failed to terminate redis container")
		}
	}()
}

func compareErrors(t *testing.T, e1, e2 error) {
	if e1 != nil && e2 == nil || e1 == nil && e2 != nil {
		t.Errorf("errors are not the same actual: %s, excpected: %s", e1, e2)
	}
	if e1 != nil && e2 != nil {
		er1, ok1 := e1.(*errs.Error)
		require.True(t, ok1, "er1 is not of type *errs.Error")
		er2, ok2 := e2.(*errs.Error)
		require.True(t, ok2, "er2 is not of type *errs.Error")
		require.True(t, reflect.DeepEqual(er1, er2), "er1 & er2 are not the same, expected:%s, actual:%s", er2, er1)
	}
}
