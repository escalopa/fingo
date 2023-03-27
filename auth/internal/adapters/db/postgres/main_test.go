package mypostgres

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/escalopa/fingo/utils/testcontainer"
)

var (
	testPGConn *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	// Create a new connection with postgres test container
	conn, terminate, err := testcontainer.NewPostgresContainer(ctx)
	if err != nil {
		log.Fatal("failed to init postgres container")
	}
	testPGConn = conn
	// Terminate the container with defer
	defer func() {
		log.Println("terminating postgres container")
		err := terminate()
		if err != nil {
			log.Fatal("failed to terminate pgContainer")
		}
	}()
	// Migrate database
	err = Migrate(conn, "file://./migrations")
	if err != nil {
		log.Fatalf("failed to migrate database for tests: %s", err)
	}
	// Run tests
	m.Run()
}
