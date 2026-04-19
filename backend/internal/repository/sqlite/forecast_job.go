package sqlite

import (
	"database/sql"
)

type SQLiteForecastJobRepository struct {
	db *sql.DB
}
