package parser

import (
	"github.com/terracodum/expensemind/backend/internal/errors"
)

func NewParser(contentType string) (Parser, error) {
	switch contentType {
	case "text/csv":
		return &CSVParser{}, nil
	case "application/pdf":
		return &TBankParser{}, nil
	default:
		return nil, errors.ParseError("unsupported file type")
	}
}
