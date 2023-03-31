package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRabbitMQContainer(ctx context.Context) (usi string, terminate func() error, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.10.19",
		ExposedPorts: []string{"5672/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete").WithStartupTimeout(time.Minute * 2),
	}
	rabbitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to start rabbit container").Err()
	}
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to start rabbitmq container").Err()
	}
	mappedPort, err := rabbitContainer.MappedPort(ctx, "5672")
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to get container port").Err()
	}

	hostIP, err := rabbitContainer.Host(ctx)
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to get container host").Err()
	}
	// Create a new connection for
	uri := fmt.Sprintf("amqp://guest:guest@%s:%s/", hostIP, mappedPort.Port())
	// Create terminate function
	terminate = func() error {
		if err := rabbitContainer.Terminate(ctx); err != nil {
			return errs.B(err).Code(errs.Unknown).Msg("failed to terminate rabbit container").Err()
		}
		return nil
	}
	return uri, terminate, nil
}
