package parser

import (
	"github.com/terracodum/expensemind/backend/internal/errors"
	"github.com/terracodum/expensemind/backend/internal/parser/pdf"
)

type Factory struct{}

func (f Factory) Create(contentType, bank string) (Parser, error) {
	switch contentType {
	case "text/csv":
		return &CSVParser{}, nil
	case "application/pdf":
		switch bank {
		case "tbank":
			return &pdf.TBankParser{}, nil
		case "vtb":
			return &pdf.VTBParser{}, nil
		case "sber":
			return &pdf.SberParser{}, nil
		case "ozon":
			return &pdf.OzonParser{}, nil
		default:
			return nil, errors.ParseError("unknown bank")
		}
	default:
		return nil, errors.ParseError("unsupported file type")
	}
}
