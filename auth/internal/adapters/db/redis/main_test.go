package redis

import (
	"context"
	"log"
	"testing"

	"github.com/escalopa/fingo/utils/testcontainer"

	"github.com/go-redis/redis/v9"
)

var (
	testRedis *redis.Client
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	url, terminate, err := testcontainer.NewRedisContainer(ctx)
	if err != nil {
		log.Fatalf("failed to start redis container for tests, err: %s", err)
	}
	testRedis, err = New(url)
	if err != nil {
		log.Fatalf("failed to create redis client for tests, err: %s", err)
	}
	m.Run()
	defer func() {
		if err = terminate(); err != nil {
			log.Fatalf("failed to terminate redis container")
		}
	}()
}
