package testcontainer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartRedisContainer(ctx context.Context) (client *redis.Client, terminate func() error, err error) {
	redisContainer, err := startRedisContainer(ctx)
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to start redis container").Err()
	}

	mappedPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to get container port").Err()
	}

	hostIP, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to get container host").Err()
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())
	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to parse redis url").Err()
	}
	client = redis.NewClient(opts)
	terminate = func() error {
		return redisContainer.Terminate(ctx)
	}
	return client, terminate, nil
}

func startRedisContainer(ctx context.Context) (testcontainers.Container, error) {
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
	return container, nil
}
