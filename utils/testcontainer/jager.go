package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewJagerContainer(ctx context.Context) (url string, terminate func() error, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "jaegertracing/all-in-one:1.6",
		ExposedPorts: []string{"14268/tcp"},
		WaitingFor:   wait.ForLog("Starting jaeger-collector").WithStartupTimeout(time.Second * 10),
	}

	jagerContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to start redis container").Err()
	}

	mappedPort, err := jagerContainer.MappedPort(ctx, "14268")
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to get container port").Err()
	}

	hostIP, err := jagerContainer.Host(ctx)
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to get container host").Err()
	}

	//http://jaeger:14268/api/traces
	uri := fmt.Sprintf("http://%s:%s/api/traces", hostIP, mappedPort.Port())
	if err != nil {
		return "", nil, errs.B(err).Code(errs.Unknown).Msg("failed to parse redis url").Err()
	}
	terminate = func() error {
		return jagerContainer.Terminate(ctx)
	}
	return uri, terminate, nil
}
