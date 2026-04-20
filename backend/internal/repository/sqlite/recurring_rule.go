package sqlite

import (
	"database/sql"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type SQLiteRecurringRuleRepository struct {
	db *sql.DB
}

func (r *SQLiteRecurringRuleRepository) Save(rule domain.RecurringRule) error {
	_, err := r.db.Exec(`
		INSERT INTO recurring_rules (source_id, type, amount, day, start_date, label)
        VALUES (?, ?, ?, ?, ?, ?)
		`, rule.SourceID, rule.Type, rule.Amount, rule.Day, rule.StartDate, rule.Label,
	)
	if err != nil {
		return errors.DBError("failed to save recurring rule", err)
	}

	return nil
}

func (r *SQLiteRecurringRuleRepository) FindActive(today time.Time) ([]domain.RecurringRule, error) {
	query := `SELECT id, source_id, type, amount, day, start_date, label
		FROM recurring_rules r
		WHERE start_date = (
		SELECT MAX(start_date) 
		FROM recurring_rules
		WHERE source_id = r.source_id AND start_date <= ?
		)`

	rows, err := r.db.Query(query, today)
	if err != nil {
		return nil, errors.DBError("failed to find recurring rule", err)
	}
	defer rows.Close()

	var result []domain.RecurringRule

	for rows.Next() {
		var id int
		var source_id sql.NullString
		var typeo sql.NullString
		var amount float64
		var day int
		var start_date time.Time
		var label sql.NullString

		if err := rows.Scan(&id, &source_id, &typeo, &amount, &day, &start_date, &label); err != nil {
			return nil, errors.DBError("failed to scan recurring rule row", err)
		}

		trans := domain.RecurringRule{ID: id, SourceID: source_id.String, Type: typeo.String, Amount: amount, Day: day, StartDate: start_date, Label: label.String}
		result = append(result, trans)
	}

	return result, nil
}

func (r *SQLiteRecurringRuleRepository) Delete(sourceID string) error {
	querry := "DELETE FROM recurring_rules WHERE source_id = ?"

	_, err := r.db.Exec(querry, sourceID)

	if err != nil {
		return errors.DBError("failed to delete recurring rule", err)
	}

	return nil
}
