package parser

import (
	"io"

	"github.com/terracodum/expensemind/backend/internal/domain"
)

type Parser interface {
	Parse(file io.Reader) ([]domain.Transaction, error)
}
