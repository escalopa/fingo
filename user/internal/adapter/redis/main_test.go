package redis

import (
	"context"
	"github.com/escalopa/fingo/utils/testcontainer"
	"github.com/go-redis/redis/v9"
	"log"
	"testing"
)

var (
	redisClient *redis.Client
	testContext context.Context
)

func TestMain(m *testing.M) {
	testContext = context.Background()
	client, terminate, err := testcontainer.StartRedisContainer(testContext)
	redisClient = client
	if err != nil {
		log.Fatalf("failed to setup redis: %s", err.Error())
	}
	m.Run()
	defer func() {
		err := terminate()
		if err != nil {
			return
		}
	}()
}
