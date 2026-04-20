package parser

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type CSVParser struct{}

func (p *CSVParser) Parse(file io.Reader) ([]domain.Transaction, error) {
	r := csv.NewReader(file)
	r.Comma = ';'
	headers, err := r.Read()
	if err != nil {
		return nil, errors.ParseError("cannot parse csv")
	}

	colIndex := map[string]int{}
	for i, h := range headers {
		colIndex[h] = i
	}

	var result []domain.Transaction

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.ParseError("cannot parse csv")
		}

		amount, err := strconv.ParseFloat(row[colIndex["amount"]], 64)
		if err != nil {
			return nil, errors.ParseError("invalid amount")
		}

		date, err := time.Parse("02.01.2006", row[colIndex["date"]])
		if err != nil {
			return nil, errors.ParseError("invalid date")
		}

		category := ""
		if idx, ok := colIndex["category"]; ok {
			category = row[idx]
		}
		if category == "" {
			category = "unknown"
		}

		description := ""
		if idx, ok := colIndex["description"]; ok {
			description = row[idx]
		}

		tran := domain.Transaction{ID: 0, Amount: amount, Date: date, Description: description, Category: category}
		result = append(result, tran)
	}

	return result, nil
}
