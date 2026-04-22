package parser_test

import (
	"testing"

	"github.com/terracodum/expensemind/backend/internal/parser"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		contentType string
		expectErr   bool
		expectType  parser.Parser
	}{
		{"text/csv", false, &parser.CSVParser{}},
		{"application/pdf", false, &parser.TBankParser{}},
		{"text/plain", true, nil},
	}

	for _, tc := range tests {
		f := parser.Factory{}
		p, err := f.Create(tc.contentType)
		if tc.expectErr {
			if err == nil {
				t.Errorf("%s: expected error, got nil", tc.contentType)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.contentType, err)
		}
		switch tc.expectType.(type) {
		case *parser.CSVParser:
			if _, ok := p.(*parser.CSVParser); !ok {
				t.Errorf("%s: expected *CSVParser, got %T", tc.contentType, p)
			}
		case *parser.TBankParser:
			if _, ok := p.(*parser.TBankParser); !ok {
				t.Errorf("%s: expected *TBankParser, got %T", tc.contentType, p)
			}
		}
	}
}
