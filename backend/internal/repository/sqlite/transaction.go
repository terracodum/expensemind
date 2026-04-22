package sqlite

import (
	"database/sql"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type SQLiteTransactionRepository struct {
	db *sql.DB
}

func (r *SQLiteTransactionRepository) scanRows(rows *sql.Rows) ([]domain.Transaction, error) {
	var result []domain.Transaction

	for rows.Next() {
		var id int
		var dateStr string
		var amount float64
		var description sql.NullString
		var category sql.NullString

		if err := rows.Scan(&id, &amount, &description, &category, &dateStr); err != nil {
			return nil, errors.DBError("failed to scan transaction row", err)
		}

		date, err := parseDBTime(dateStr)
		if err != nil {
			return nil, errors.DBError("failed to parse transaction date", err)
		}

		trans := domain.Transaction{ID: id, Date: date, Amount: amount, Description: description.String, Category: category.String}
		result = append(result, trans)
	}

	return result, nil
}

func (r *SQLiteTransactionRepository) Save(tx domain.Transaction) error {
	_, err := r.db.Exec(`
		INSERT INTO transactions (amount, description, category, date) VALUES (?, ?, ?, ?)`,
		tx.Amount, tx.Description, tx.Category, tx.Date,
	)

	if err != nil {
		return errors.DBError("failed to save transaction", err)
	}

	return nil
}

func (r *SQLiteTransactionRepository) SaveAll(transaction []domain.Transaction) error {
	tx, err := r.db.Begin()

	if err != nil {
		return errors.DBError("failed to save transactions", err)
	}

	for _, item := range transaction {

		_, err := tx.Exec(`
		INSERT INTO transactions (amount, description, category, date) VALUES (?, ?, ?, ?)`,
			item.Amount, item.Description, item.Category, item.Date,
		)

		if err != nil {
			tx.Rollback()
			return errors.DBError("failed to save transactions", err)
		}

	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.DBError("failed to save transactions", err)
	}

	return nil
}

func (r *SQLiteTransactionRepository) FindAll(filters domain.Filters) ([]domain.Transaction, error) {
	query := "SELECT id, amount, description, category, date FROM transactions WHERE 1=1"
	args := []any{}

	if !filters.From.IsZero() {
		query += " AND date >= ?"
		args = append(args, filters.From)
	}

	if !filters.To.IsZero() {
		query += " AND date <= ?"
		args = append(args, filters.To)
	}

	if filters.Category != "" {
		query += " AND category = ?"
		args = append(args, filters.Category)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, errors.DBError("failed to find transactions", err)
	}
	defer rows.Close()

	result, err := r.scanRows(rows)
	if err != nil {
		return nil, errors.DBError("failed to find transactions", err)
	}
	return result, nil
}

func (r *SQLiteTransactionRepository) FindForForecast(from, to time.Time) ([]domain.Transaction, error) {
	query := "SELECT id, amount, description, category, date FROM transactions WHERE date >= ? AND date <= ?"

	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, errors.DBError("failed to get transactions by date range", err)
	}
	defer rows.Close()

	result, err := r.scanRows(rows)
	if err != nil {
		return nil, errors.DBError("failed to get transactions by date range", err)
	}

	return result, nil
}

func (r *SQLiteTransactionRepository) Update(tx domain.Transaction) error {
	querry := "UPDATE transactions SET amount = ?, description = ?, category = ?, date = ? WHERE id = ?"

	_, err := r.db.Exec(querry, tx.Amount, tx.Description, tx.Category, tx.Date, tx.ID)

	if err != nil {
		return errors.DBError("failed to update transaction", err)
	}

	return nil
}

func (r *SQLiteTransactionRepository) Delete(id int) error {
	querry := "DELETE FROM transactions WHERE id = ?"

	_, err := r.db.Exec(querry, id)

	if err != nil {
		return errors.DBError("failed to delete transaction", err)
	}

	return nil
}
