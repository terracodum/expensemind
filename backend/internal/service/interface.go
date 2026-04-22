package service

import (
	"io"

	"github.com/terracodum/expensemind/backend/internal/domain"
)

type Service interface {
	UploadTransactions(contentType string, file io.Reader) error
	GetTransactions(filters domain.Filters) ([]domain.Transaction, error)
	UpdateTransaction(tx domain.Transaction) error
	DeleteTransaction(id int) error
	CreateForecastJob() (int, error)
	GetForecastJob(id int) (domain.ForecastJob, error)
	GetForecastJobs() ([]domain.ForecastJob, error)
	SaveRecurringRule(rule domain.RecurringRule) error
	GetRecurringRules() ([]domain.RecurringRule, error)
	DeleteRecurringRule(sourceID string) error
}
