package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

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
