package repository

import (
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
)

type Filters struct {
	From     time.Time
	To       time.Time
	Category string
}

type TransactionRepository interface {
	Save(tx domain.Transaction) error
	SaveAll(txs []domain.Transaction) error
	FindAll(filters Filters) ([]domain.Transaction, error)
	FindByDateRange(from, to time.Time) ([]domain.Transaction, error)
	Update(tx domain.Transaction) error
	Delete(id int) error
}
