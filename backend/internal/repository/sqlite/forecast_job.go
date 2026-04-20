package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type SQLiteForecastJobRepository struct {
	db *sql.DB
}

func (r *SQLiteForecastJobRepository) unpackJob(id int, status sql.NullString, result sql.NullString, createdAt time.Time) (domain.ForecastJob, error) {
	var forecast domain.Forecast
	if result.Valid {
		err := json.Unmarshal([]byte(result.String), &forecast)
		if err != nil {
			return domain.ForecastJob{}, errors.DBError("failed to unpack result", err)
		}
		return domain.ForecastJob{ID: id, Status: status.String, Result: &forecast, CreatedAt: createdAt}, nil
	}
	return domain.ForecastJob{ID: id, Status: status.String, Result: nil, CreatedAt: createdAt}, nil
}

func (r *SQLiteForecastJobRepository) Create() (int, error) {
	res, err := r.db.Exec(`
        INSERT INTO forecast_jobs (status, created_at) VALUES ('pending', ?)`,
		time.Now(),
	)
	if err != nil {
		return 0, errors.DBError("failed to create forecast job", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.DBError("failed to get forecast job id", err)
	}

	return int(id), nil

}

func (r *SQLiteForecastJobRepository) FindByID(id int) (domain.ForecastJob, error) {
	query := `SELECT id, status, result, created_at FROM forecast_jobs WHERE id=? `
	row := r.db.QueryRow(query, id)
	var (
		baseId     int
		status     sql.NullString
		result     sql.NullString
		created_at time.Time
	)
	err := row.Scan(&baseId, &status, &result, &created_at)
	if err != nil {
		return domain.ForecastJob{}, errors.DBError("failed to get forecast job", err)
	}

	return r.unpackJob(baseId, status, result, created_at)
}

func (r *SQLiteForecastJobRepository) FindAll() ([]domain.ForecastJob, error) {
	query := `SELECT id, status, result, created_at FROM forecast_jobs`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.DBError("failed to get forecast job", err)
	}
	var jobs []domain.ForecastJob

	for rows.Next() {
		var (
			baseId     int
			status     sql.NullString
			result     sql.NullString
			created_at time.Time
		)

		err := rows.Scan(&baseId, &status, &result, &created_at)
		if err != nil {
			return []domain.ForecastJob{}, errors.DBError("failed to get forecast job", err)
		}

		job, err := r.unpackJob(baseId, status, result, created_at)
		if err != nil {
			return []domain.ForecastJob{}, errors.DBError("failed to get forecast job", err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *SQLiteForecastJobRepository) Update(job domain.ForecastJob) error {
	query := `UPDATE forecast_jobs SET status=?, result=? WHERE id=?`

	if job.Result != nil {
		res, err := json.Marshal(job.Result)
		if err != nil {
			return errors.DBError("failed to pack ForecastJob", err)
		}
		_, err = r.db.Exec(query, job.Status, res, job.ID)
		if err != nil {
			return errors.DBError("failed to update ForecastJob", err)
		}
	} else {
		_, err := r.db.Exec(query, job.Status, nil, job.ID)
		if err != nil {
			return errors.DBError("failed to update ForecastJob", err)
		}
	}
	return nil
}
