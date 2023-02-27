package db

import (
	"database/sql"
	"github.com/escalopa/gochat/utils/testcontainer"
	"log"
	"testing"
)

var dbSQL *sql.DB

func TestMain(m *testing.M) {
	conn, terminate, err := testcontainer.StartPostgresContainer()
	if err != nil {
		log.Fatalf("failer to start postgres container for test, err: %s", err)
	}
	dbSQL = conn
	m.Run()
	defer func() {
		err := terminate()
		if err != nil {
			log.Fatalf("failed to terminate postgres container, err: %s", err)
		}
	}()
}
