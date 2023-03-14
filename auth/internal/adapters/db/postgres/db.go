package mypostgres

import (
	"database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/lordvidex/errs"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"log"
)

// New creates a new postgres connection with the given connection string
func New(url string) (*sql.DB, error) {
	// Creates a new postgres conn
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errs.B(err).Msg("failed to open postgres connection").Err()
	}
	return conn, nil
}

// Migrate runs the migrations on the database
func Migrate(conn *sql.DB, migrationDir string) error {
	// Create a new pg instance
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return errs.B(err).Msg("failed to create driver for migration").Err()
	}
	// Load migration files
	m, err := migrate.NewWithDatabaseInstance(
		migrationDir,
		"postgres", driver)
	if err != nil {
		return errs.B(err).Msg("failed to create new pg instance for migration").Err()
	}
	// Push migration changes
	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Printf("Databse is up to date, No migration made")
	}
	return nil
}

// isUniqueViolationError checks if an error is a unique violation error
func isUniqueViolationError(err error) bool {
	er, ok := err.(*pq.Error)
	return ok && er.Code == "23505"
}
