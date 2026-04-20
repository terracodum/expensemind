package parser_test

import (
	"os"
	"testing"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/parser"
)

func TestTBankParse_Structure(t *testing.T) {
	f, err := os.Open("testdata/statement.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	txs, err := (&parser.TBankParser{}).Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	expected := []domain.Transaction{
		{Date: time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC), Amount: -100.00, Description: "Оплата в Сервисы Яндекса"},
		{Date: time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC), Amount: -40.00, Description: "Оплата вETK55_OMSK_TPP OmskRUS"},
		{Date: time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC), Amount: 300.00, Description: "Пополнение. Системабыстрых платежей"},
		{Date: time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC), Amount: -20.00, Description: "Внутренний перевод надоговор 8352982596"},
		{Date: time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC), Amount: -580.00, Description: "Оплата в tulen.storeCHeboksary' RU"},
	}

	if len(txs) < len(expected) {
		t.Fatalf("expected at least %d transactions, got %d", len(expected), len(txs))
	}

	for i, exp := range expected {
		got := txs[i]
		if got.Date != exp.Date {
			t.Errorf("[%d] date: expected %v, got %v", i, exp.Date, got.Date)
		}
		if got.Amount != exp.Amount {
			t.Errorf("[%d] amount: expected %v, got %v", i, exp.Amount, got.Amount)
		}
		if got.Description != exp.Description {
			t.Errorf("[%d] description: expected %q, got %q", i, exp.Description, got.Description)
		}
	}
}
