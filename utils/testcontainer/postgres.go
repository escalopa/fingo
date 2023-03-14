package testcontainer

import (
	"context"
	"database/sql"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"

	"fmt"
	_ "github.com/lib/pq"
	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
)

func NewPostgresContainer() (dbSQL *sql.DB, terminate func() error, err error) {
	dbUser := "postgres"
	dbPass := "postgres"
	dbDB := "fingo"
	// Run container
	ctx := context.Background()
	pgContainer, err := spinPostgresContainer(ctx,
		withPort("5432/tcp"),
		withInitialDatabase(dbUser, dbPass, dbDB),
		withWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, errs.B().Msg("failed to start postgres container").Err()

	}
	// Get pgContainer host
	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to get postgres container host").Err()

	}
	// Get container port
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to get postgres container port").Err()
	}
	// Create connection URL & connect to pg instance
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, host, port.Port(), dbDB)
	dbSQL, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to open postgres connection").Err()
	}
	if err = dbSQL.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("failed to ping pg test container: %s", err))
	}
	// Create terminate function to terminate container when done using it
	terminate = func() error {
		return pgContainer.Terminate(ctx)
	}
	return dbSQL, terminate, nil
}

// postgresContainer represents the postgres container type used in the module
type postgresContainer struct {
	testcontainers.Container
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

// spinPostgresContainer creates an instance of the postgres container type
func spinPostgresContainer(ctx context.Context, opts ...postgresContainerOption) (*postgresContainer, error) {
	// Create request object with default values
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		Env:          map[string]string{},
		ExposedPorts: []string{},
	}
	// Apply options
	for _, opt := range opts {
		opt(&req)
	}
	// Create container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	// Return container
	return &postgresContainer{Container: container}, nil
}
