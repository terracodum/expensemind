package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

var dbTimeFormats = []string{
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func parseDBTime(s string) (time.Time, error) {
	for _, f := range dbTimeFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time: %s", s)
}

func New(s *sql.DB) (SQLiteTransactionRepository, SQLiteRecurringRuleRepository, SQLiteForecastJobRepository, error) {
	_, err := s.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
				id          INTEGER PRIMARY KEY AUTOINCREMENT,
				amount      REAL    NOT NULL,
				description TEXT,
				category    TEXT,
				date        TEXT    NOT NULL
			);
	`)

	if err != nil {
		return SQLiteTransactionRepository{}, SQLiteRecurringRuleRepository{}, SQLiteForecastJobRepository{}, errors.DBError("failed to create transactions table", err)
	}
	_, err = s.Exec(`
        CREATE TABLE IF NOT EXISTS recurring_rules (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            source_id  TEXT    NOT NULL,
            type       TEXT    NOT NULL,
            amount     REAL    NOT NULL,
            day        INTEGER NOT NULL,
            start_date TEXT    NOT NULL,
            label      TEXT,
            UNIQUE(source_id, start_date)
        );
    `)
	if err != nil {
		return SQLiteTransactionRepository{}, SQLiteRecurringRuleRepository{}, SQLiteForecastJobRepository{}, errors.DBError("failed to create recurring_rules table", err)
	}

	_, err = s.Exec(`
        CREATE TABLE IF NOT EXISTS forecast_jobs (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            status     TEXT    NOT NULL,
            result     TEXT,
            created_at TEXT    NOT NULL
        );
    `)
	if err != nil {
		return SQLiteTransactionRepository{}, SQLiteRecurringRuleRepository{}, SQLiteForecastJobRepository{}, errors.DBError("failed to create forecast_jobs table", err)
	}

	return SQLiteTransactionRepository{s}, SQLiteRecurringRuleRepository{s}, SQLiteForecastJobRepository{s}, nil
}
