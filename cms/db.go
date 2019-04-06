package cms

import (
	"database/sql"

	// Use the PG SQL driver
	_ "github.com/lib/pq"
)

type PgStore struct {
	DB *sql.DB
}
