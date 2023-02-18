package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"reflect"
	"testing"
)

var (
	testRedis   *redis.Client
	testContext context.Context
)

func TestMain(m *testing.M) {
	testContext = context.Background()
	redisC, err := startContainer(testContext)
	if err != nil {
		log.Fatalf("failed to setup redis: %s", err.Error())
	}
	m.Run()
	err = terminateRedis(redisC)
}

func startContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())
	testRedis, err = New(uri)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func terminateRedis(c testcontainers.Container) error {
	return c.Terminate(testContext)
}

func compareErrors(t *testing.T, e1, e2 error) {
	if e1 != nil && e2 == nil || e1 == nil && e2 != nil {
		t.Errorf("e1ors are not the same actual:%s, excpected:%s", e1, e2)
	}
	if e1 != nil && e2 != nil {
		if er1, ok := e1.(*errs.Error); ok {
			if er2, ok := e2.(*errs.Error); ok {
				require.False(t, reflect.DeepEqual(er1.Msg, er2.Msg), "e1ors are not equal actual:%s, expected:%s", e1, e2)
			} else {
				t.Errorf("e2 is not of type *e1s.Err: actual:%T, excpected:%T", e1, e2)
			}
		} else {
			t.Errorf("e1 is not of type *e1s.Err: actual:%T, excpected:%T", e1, e2)
		}
	}
}
