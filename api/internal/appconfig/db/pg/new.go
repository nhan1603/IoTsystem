package pg

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Connect connects to database
func Connect(dbURL string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)

		return nil, errors.WithStack(fmt.Errorf("connect DB failed. Err: %w", err))
	}

	if err = conn.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("Initializing DB connection")

	return conn, nil
}
