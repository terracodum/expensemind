package pdf

import (
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type SberParser struct{}

var (
	// Line A: DD.MM.YYYY   HH:MM   Category...
	sberDateLineRe = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}\s+\d{2}:\d{2}\s`)
	// Amount and balance at end of line A: +?N NNN,DD   N NNN,DD
	sberAmountRe = regexp.MustCompile(`(\+?\d{1,3}(?:\s\d{3})*,\d{2})\s+(\d{1,3}(?:\s\d{3})*,\d{2})\s*$`)
	// Line B: DD.MM.YYYY   AUTHCODE   Description
	sberDescBRe = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}\s+\d+\s+(.+)$`)
	// Card number suffix line e.g. ****0717
	sberCardRe = regexp.MustCompile(`^\*{4}\d+$`)
)

func (p *SberParser) Parse(file io.Reader) ([]domain.Transaction, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.ParseError("cannot read pdf")
	}

	tmp, err := os.CreateTemp("", "sber-*.pdf")
	if err != nil {
		return nil, errors.ParseError("cannot create temp file")
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return nil, errors.ParseError("cannot write temp file")
	}
	tmp.Close()

	out, err := exec.Command("pdftotext", "-layout", tmp.Name(), "-").Output()
	if err != nil {
		return nil, errors.ParseError("pdftotext failed: install poppler-utils")
	}

	return parseSberText(string(out))
}

func parseSberText(text string) ([]domain.Transaction, error) {
	lines := strings.Split(text, "\n")

	var blocks [][]string
	var current []string

	for _, line := range lines {
		if sberDateLineRe.MatchString(line) {
			if current != nil {
				blocks = append(blocks, current)
			}
			current = []string{line}
		} else if current != nil {
			if strings.TrimSpace(line) == "" || strings.ContainsRune(line, '\x0c') {
				blocks = append(blocks, current)
				current = nil
			} else {
				current = append(current, line)
			}
		}
	}
	if current != nil {
		blocks = append(blocks, current)
	}

	var result []domain.Transaction
	for _, block := range blocks {
		if len(block) < 2 {
			continue
		}

		date, err := time.Parse("02.01.2006", block[0][:10])
		if err != nil {
			continue
		}

		amtMatch := sberAmountRe.FindStringSubmatch(block[0])
		if amtMatch == nil {
			continue
		}

		amtStr := strings.ReplaceAll(amtMatch[1], " ", "")
		amtStr = strings.ReplaceAll(amtStr, ",", ".")
		positive := strings.HasPrefix(amtMatch[1], "+")
		if positive {
			amtStr = amtStr[1:]
		}
		amount, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			continue
		}
		if !positive {
			amount = -amount
		}

		var descParts []string
		if m := sberDescBRe.FindStringSubmatch(block[1]); m != nil {
			descParts = append(descParts, strings.TrimSpace(m[1]))
		}
		for _, line := range block[2:] {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !sberCardRe.MatchString(trimmed) {
				descParts = append(descParts, trimmed)
			}
		}

		result = append(result, domain.Transaction{
			Date:        date,
			Amount:      amount,
			Description: strings.Join(descParts, " "),
		})
	}

	if len(result) == 0 {
		return nil, errors.InvalidPDFFormat("no transactions found in Sber pdf")
	}
	return result, nil
}
