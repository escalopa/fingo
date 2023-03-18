package redis

import (
	"context"
	"log"
	"testing"

	"github.com/escalopa/fingo/utils/testcontainer"
	"github.com/go-redis/redis/v9"
)

var (
	redisClient *redis.Client
	testContext context.Context
)

func TestMain(m *testing.M) {
	testContext = context.Background()
	client, terminate, err := testcontainer.NewRedisContainer(testContext)
	if err != nil {
		log.Fatalf("failed to setup redis: %s", err.Error())
	}
	redisClient, err = New(client)
	if err != nil {
		log.Fatalf("failed to create redis client: %s", err.Error())
	}
	m.Run()
	defer func() {
		err := terminate()
		if err != nil {
			return
		}
	}()
}
