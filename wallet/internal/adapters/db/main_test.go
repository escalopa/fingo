package db

import (
	"context"
	"database/sql"
	"log"
	"testing"

	pkgdb "github.com/escalopa/fingo/pkg/db"
	"github.com/escalopa/fingo/utils/testcontainer"
)

var (
	conn *sql.DB
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, terminate, err := testcontainer.NewPostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	conn = dbConn
	err = pkgdb.Migrate(conn, "file://./sql/migrations")
	if err != nil {
		log.Fatal(err)
	}
	defer terminate()
	m.Run()
}
