package testcontainer

import (
	"context"
	"database/sql"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"

	"fmt"

	_ "github.com/lib/pq"
	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
)

func NewPostgresContainer(ctx context.Context) (dbSQL *sql.DB, terminate func() error, err error) {
	dbUser := "postgres"
	dbPass := "postgres"
	dbDB := "fingo"
	// Run container
	// Create request object with default values
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		Env:          map[string]string{},
		ExposedPorts: []string{},
	}
	opts := []postgresContainerOption{
		withPort("5432/tcp"),
		withInitialDatabase(dbUser, dbPass, dbDB),
		withWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second)),
	}
	// Apply options
	for _, opt := range opts {
		opt(&req)
	}
	// Create container
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, errs.B().Msg("failed to start postgres container").Err()

	}
	// Get pgContainer host
	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, nil, errs.B(err).Code(errs.Unknown).Msg("failed to get postgres container host").Err()
	}
	// Get container port
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, errs.B(err).Code(errs.Unknown).Msg("failed to get postgres container port").Err()
	}
	// Create connection URL & connect to pg instance
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, host, port.Port(), dbDB)
	dbSQL, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, nil, errs.B(err).Code(errs.Unknown).Msg("failed to open postgres connection").Err()
	}
	if err = dbSQL.Ping(); err != nil {
		return nil, nil, errs.B(err).Msg("failed to ping pg test container:").Err()
	}
	// Create terminate function to terminate container when done using it
	terminate = func() error {
		return pgContainer.Terminate(ctx)
	}
	return dbSQL, terminate, nil
}

type postgresContainerOption func(req *testcontainers.ContainerRequest)

func withInitialDatabase(user string, password string, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}

func withWaitStrategy(strategies ...wait.Strategy) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.WaitingFor = wait.ForAll(strategies...).WithDeadline(1 * time.Minute)
	}
}

func withPort(port string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = append(req.ExposedPorts, port)
	}
}
