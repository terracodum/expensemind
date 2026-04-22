package repository

import (
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
)

type TransactionRepository interface {
	Save(tx domain.Transaction) error
	SaveAll(txs []domain.Transaction) error
	FindAll(filters domain.Filters) ([]domain.Transaction, error)
	FindForForecast(from, to time.Time) ([]domain.Transaction, error)
	Update(tx domain.Transaction) error
	Delete(id int) error
}

type RecurringRuleRepository interface {
	Save(rule domain.RecurringRule) error
	FindAll() ([]domain.RecurringRule, error)
	FindActive(today time.Time) ([]domain.RecurringRule, error)
	Delete(sourceID string) error
}

type ForecastJobRepository interface {
	Create() (int, error)
	FindByID(id int) (domain.ForecastJob, error)
	FindAll() ([]domain.ForecastJob, error)
	Update(job domain.ForecastJob) error
}
