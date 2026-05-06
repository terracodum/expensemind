package pdf

import (
	"bytes"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type VTBParser struct{}

func (p *VTBParser) Parse(file io.Reader) ([]domain.Transaction, error) {
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

	return parseVTBText(sb.String())
}

func parseVTBText(text string) ([]domain.Transaction, error) {
	// TODO: implement VTB PDF parsing
	_ = text
	return nil, errors.InvalidPDFFormat("VTB parser not implemented")
}
