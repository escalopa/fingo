package mypostgres

import (
	"database/sql"
	"github.com/escalopa/gochat/utils/testcontainer"
	"log"
	"testing"
)

var (
	dbSQL *sql.DB
)

func TestMain(m *testing.M) {
	conn, terminate, err := testcontainer.StartPostgresContainer()
	if err != nil {
		log.Fatal("failed to init postgres container")
	}
	dbSQL = conn
	m.Run()
	defer func() {
		err := terminate()
		if err != nil {
			log.Fatal("failed to terminate pgContainer")
		}
	}()
}
