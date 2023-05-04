package testcontainer

import (
	"context"
	"fmt"

	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRedisContainer(ctx context.Context) (url string, terminate func() error, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to start redis container").Err()
	}

	mappedPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to get container port").Err()
	}

	hostIP, err := redisContainer.Host(ctx)
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to get container host").Err()
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to parse redis url").Err()
	}
	terminate = func() error {
		return redisContainer.Terminate(ctx)
	}
	return uri, terminate, nil
}
