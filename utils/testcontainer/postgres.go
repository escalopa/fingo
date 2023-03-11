package testcontainer

import (
	"context"
	"database/sql"

	"fmt"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/lib/pq"
	"github.com/lordvidex/errs"
	"github.com/testcontainers/testcontainers-go"
)

func StartPostgresContainer() (dbSQL *sql.DB, terminate func() error, err error) {
	dbUser := "posgtres"
	dbPass := "postgres"
	dbDB := "fingo"

	port := strconv.Itoa(gofakeit.IntRange(20_000, 30_000))
	// Run container
	ctx := context.Background()
	pgContainer, err := spinPostgresContainer(ctx,
		withInitialDatabase(dbUser, dbPass, dbDB),
		withPort(port),
	)
	if err != nil {
		return nil, nil, errs.B().Msg("failed to start postgres container").Err()

	}
	// Get pgContainer host
	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, nil, errs.B().Msg("failed to get postgres container host").Err()

	}
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, host, port, dbDB,
	)
	dbSQL, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, nil, errs.B(err).Msg("failed to open postgres connection").Err()
	}
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

// func withWaitStrategy(strategies ...wait.Strategy) func(req *testcontainers.ContainerRequest) {
// 	return func(req *testcontainers.ContainerRequest) {
// 		req.WaitingFor = wait.ForAll(strategies...).WithDeadline(1 * time.Minute)
// 	}
// }

func withPort(port string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = append(req.ExposedPorts, port)
	}
}

func withInitialDatabase(user string, password string, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}

// spinPostgresContainer creates an instance of the postgres container type
func spinPostgresContainer(ctx context.Context, opts ...postgresContainerOption) (*postgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12.8",
		Env:          map[string]string{},
		ExposedPorts: []string{},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
	}

	for _, opt := range opts {
		opt(&req)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &postgresContainer{Container: container}, nil
}
