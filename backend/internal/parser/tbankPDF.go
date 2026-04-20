package parser

import (
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type TBankParser struct{}

func (p *TBankParser) Parse(file io.Reader) ([]domain.Transaction, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.ParseError("cannot read pdf")
	}

	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, errors.ParseError("cannot open pdf")
	}

	var sb strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return nil, errors.ParseError("cannot read pdf page")
		}
		sb.WriteString(text)
		sb.WriteByte('\n')
	}

	return parseTBankText(sb.String())
}

var txRe = regexp.MustCompile(
	`(\d{2}\.\d{2}\.\d{4})\n\d{2}:\d{2}\n` + // дата1 + время1
		`\d{2}\.\d{2}\.\d{4}\n\d{2}:\d{2}\n` + // дата2 + время2
		`[+\-][\d ]+\.\d{2} ₽\n` + // сумма1 (пропускаем)
		`([+\-][\d ]+\.\d{2}) ₽\n` + // сумма2 (захватываем)
		`(.+)\n` + // описание
		`(\d{4}|—)\n`, // карта
)

func parseTBankText(text string) ([]domain.Transaction, error) {
	matches := txRe.FindAllStringSubmatch(text, -1)
	if matches == nil {
		return nil, errors.ParseError("no transactions found in pdf")
	}

	result := make([]domain.Transaction, 0, len(matches))
	for _, m := range matches {
		date, err := time.Parse("02.01.2006", m[1])
		if err != nil {
			return nil, errors.ParseError("invalid date in pdf")
		}

		amount, err := parseAmount(m[2])
		if err != nil {
			return nil, errors.ParseError("invalid amount in pdf")
		}

		desc := strings.TrimSpace(strings.ReplaceAll(m[3], "\n", " "))

		result = append(result, domain.Transaction{
			Date:        date,
			Amount:      amount,
			Description: desc,
		})
	}
	return result, nil
}

func parseAmount(s string) (float64, error) {
	s = strings.ReplaceAll(s, " ", "")
	return strconv.ParseFloat(s, 64)
}
