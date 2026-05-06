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

type OzonParser struct{}

var (
	ozonAmountRe   = regexp.MustCompile(`([+-])\s*(\d[\d ]*\.\d{2})\s*₽`)
	ozonDateRe     = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4})`)
	ozonPureTimeRe = regexp.MustCompile(`^\s*\d{2}:\d{2}:\d{2}\s*$`)
	ozonTimePfxRe  = regexp.MustCompile(`^\s*\d{2}:\d{2}:\d{2}\s+(.+)$`)
	ozonDocIDRe    = regexp.MustCompile(`^\d{7,12}\s*`)
	ozonDatePfxRe  = regexp.MustCompile(`^\s*\d{2}\.\d{2}\.\d{4}(\s+\d{2}:\d{2}:\d{2})?\s*`)
)

func (p *OzonParser) Parse(file io.Reader) ([]domain.Transaction, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.ParseError("cannot read pdf")
	}

	tmp, err := os.CreateTemp("", "ozon-*.pdf")
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

	return parseOzonText(string(out))
}

func parseOzonText(text string) ([]domain.Transaction, error) {
	lines := strings.Split(text, "\n")

	// Find indices of all amount lines (contain ₽ with a sign)
	var rubIdxs []int
	for i, l := range lines {
		if ozonAmountRe.MatchString(l) {
			rubIdxs = append(rubIdxs, i)
		}
	}

	var result []domain.Transaction
	for k, rubIdx := range rubIdxs {
		line := lines[rubIdx]

		m := ozonAmountRe.FindStringSubmatch(line)
		mIdx := ozonAmountRe.FindStringIndex(line)
		if m == nil || mIdx == nil {
			continue
		}
		sign := 1.0
		if m[1] == "-" {
			sign = -1.0
		}
		amtStr := strings.ReplaceAll(m[2], " ", "")
		amount, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			continue
		}
		amount *= sign

		// Find date: inline first, then look back
		var date time.Time
		dateInline := false
		if dm := ozonDateRe.FindStringSubmatch(line); dm != nil {
			d, err := time.Parse("02.01.2006", dm[1])
			if err == nil {
				date = d
				dateInline = true
			}
		}
		if !dateInline {
			for back := 1; back <= 5; back++ {
				if rubIdx-back < 0 {
					break
				}
				prev := lines[rubIdx-back]
				if ozonAmountRe.MatchString(prev) {
					break
				}
				if dm := ozonDateRe.FindStringSubmatch(prev); dm != nil {
					if d, err := time.Parse("02.01.2006", dm[1]); err == nil {
						date = d
						break
					}
				}
			}
		}
		if date.IsZero() {
			continue
		}

		var descParts []string

		// Description from amount line (strip date+time prefix and docId)
		before := strings.TrimSpace(line[:mIdx[0]])
		before = ozonDatePfxRe.ReplaceAllString(before, "")
		before = strings.TrimSpace(before)
		before = ozonDocIDRe.ReplaceAllString(before, "")
		before = strings.TrimSpace(before)

		// Description ABOVE: look back, continue past date lines
		aboveFound := 0
		for back := 1; back <= 7; back++ {
			if rubIdx-back < 0 {
				break
			}
			pl := lines[rubIdx-back]
			ps := strings.TrimSpace(pl)
			if ps == "" {
				break
			}
			if ozonAmountRe.MatchString(pl) {
				break
			}
			if ozonPureTimeRe.MatchString(pl) {
				continue
			}
			if dm := ozonDateRe.FindStringSubmatch(pl); dm != nil {
				afterDate := strings.TrimSpace(ozonDatePfxRe.ReplaceAllString(pl, ""))
				if afterDate != "" {
					descParts = append([]string{afterDate}, descParts...)
					aboveFound++
				}
				continue // keep looking above the date line
			}
			descParts = append([]string{ps}, descParts...)
			aboveFound++
			if aboveFound >= 2 {
				break
			}
			break
		}

		if before != "" {
			descParts = append(descParts, before)
		}

		// Description BELOW: only for non-inline date; stop near next ₽
		if !dateInline {
			nextRub := len(lines)
			if k+1 < len(rubIdxs) {
				nextRub = rubIdxs[k+1]
			}
			for fwd := 1; fwd <= 4; fwd++ {
				if rubIdx+fwd >= len(lines) {
					break
				}
				nl := lines[rubIdx+fwd]
				ns := strings.TrimSpace(nl)
				if ns == "" {
					break
				}
				if ozonAmountRe.MatchString(nl) {
					break
				}
				if ozonPureTimeRe.MatchString(nl) {
					continue
				}
				// Stop if too close to the next ₽ (those lines belong to next txn)
				if rubIdx+fwd+2 >= nextRub {
					break
				}
				if m2 := ozonTimePfxRe.FindStringSubmatch(nl); m2 != nil {
					if t := strings.TrimSpace(m2[1]); t != "" {
						descParts = append(descParts, t)
					}
				} else if ozonDateRe.MatchString(ns) && len(ns) <= 10 {
					break
				} else {
					descParts = append(descParts, ns)
				}
			}
		}

		result = append(result, domain.Transaction{
			Date:        date,
			Amount:      amount,
			Description: strings.Join(descParts, " "),
		})
	}

	if len(result) == 0 {
		return nil, errors.InvalidPDFFormat("no transactions found in Ozon pdf")
	}
	return result, nil
}
