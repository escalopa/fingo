package testcontainer

import (
	"context"
	"fmt"
	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func NewRabbitMQContainer(ctx context.Context) (usi string, terminate func() error, err error) {
	rabbitContainer, err := spinRabbitmqContainer(ctx)
	if err != nil {
		return "", nil, errs.B(err).Msg("failed to start rabbitmq container").Err()
	}
	mappedPort, err := rabbitContainer.MappedPort(ctx, "5672")
	if err != nil {
		return "", nil, errs.B(err).Msg("failed to get container port").Err()
	}

	hostIP, err := rabbitContainer.Host(ctx)
	if err != nil {
		return "", nil, errs.B(err).Msg("failed to get container host").Err()
	}
	// Create a new connection for
	uri := fmt.Sprintf("amqp://guest:guest@%s:%s/", hostIP, mappedPort.Port())
	// Create terminate function
	terminate = func() error {
		return rabbitContainer.Terminate(ctx)
	}
	return uri, terminate, nil
}

func spinRabbitmqContainer(ctx context.Context) (testcontainers.Container, error) {

	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.10.19",
		ExposedPorts: []string{"5672/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete").WithStartupTimeout(time.Minute * 2),
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
