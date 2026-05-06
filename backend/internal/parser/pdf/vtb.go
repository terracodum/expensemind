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

type VTBParser struct{}

var (
	vtbDateLineRe = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}\s`)
	vtbAmountRe   = regexp.MustCompile(`(-?[\d,]+\.\d{2}) RUB`)
	vtbDescARe    = regexp.MustCompile(`\s{3,}0\.00\s+(.+)$`)
	vtbDescBRe    = regexp.MustCompile(`RUB\s{5,}(.+)$`)
	vtbDescCRe    = regexp.MustCompile(`RUB\s+RUB\s+(.+)$`)
)

func (p *VTBParser) Parse(file io.Reader) ([]domain.Transaction, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.ParseError("cannot read pdf")
	}

	tmp, err := os.CreateTemp("", "vtb-*.pdf")
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

	return parseVTBText(string(out))
}

func parseVTBText(text string) ([]domain.Transaction, error) {
	lines := strings.Split(text, "\n")

	var blocks [][]string
	var current []string

	for _, line := range lines {
		if vtbDateLineRe.MatchString(line) {
			if current != nil {
				blocks = append(blocks, current)
			}
			current = []string{line}
		} else if current != nil {
			if strings.ContainsRune(line, '\x0c') {
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

		amtMatch := vtbAmountRe.FindStringSubmatch(block[1])
		if amtMatch == nil {
			continue
		}
		amtStr := strings.ReplaceAll(amtMatch[1], ",", "")
		amount, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			continue
		}

		var descParts []string

		if m := vtbDescARe.FindStringSubmatch(block[0]); m != nil {
			descParts = append(descParts, strings.TrimSpace(m[1]))
		}
		if m := vtbDescBRe.FindStringSubmatch(block[1]); m != nil {
			descParts = append(descParts, strings.TrimSpace(m[1]))
		}
		if len(block) > 2 {
			if m := vtbDescCRe.FindStringSubmatch(block[2]); m != nil {
				descParts = append(descParts, strings.TrimSpace(m[1]))
			}
		}
		for _, line := range block[3:] {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			col := len(line) - len(strings.TrimLeft(line, " "))
			if col >= 75 && col <= 115 {
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
		return nil, errors.InvalidPDFFormat("no transactions found in VTB pdf")
	}
	return result, nil
}
