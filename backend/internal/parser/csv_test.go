package parser_test

import (
	"os"
	"testing"
	"time"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/parser"
)

func TestCSVParse_Valid(t *testing.T) {
	f, err := os.Open("testdata/valid.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	txs, err := (&parser.CSVParser{}).Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	expected := domain.Transaction{
		Date:        time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC),
		Amount:      -100.00,
		Description: "Яндекс",
		Category:    "Подписки",
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
	got := txs[0]
	if got.Date != expected.Date {
		t.Errorf("date: expected %v, got %v", expected.Date, got.Date)
	}
	if got.Amount != expected.Amount {
		t.Errorf("amount: expected %v, got %v", expected.Amount, got.Amount)
	}
	if got.Description != expected.Description {
		t.Errorf("description: expected %q, got %q", expected.Description, got.Description)
	}
	if got.Category != expected.Category {
		t.Errorf("category: expected %q, got %q", expected.Category, got.Category)
	}
}

func TestCSVParse_NoCategory(t *testing.T) {
	f, err := os.Open("testdata/no_category.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	txs, err := (&parser.CSVParser{}).Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
	if txs[0].Category != "unknown" {
		t.Errorf("expected unknown, got %q", txs[0].Category)
	}
}

func TestCSVParse_EmptyCategory(t *testing.T) {
	f, err := os.Open("testdata/empty_category.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	txs, err := (&parser.CSVParser{}).Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
	if txs[0].Category != "unknown" {
		t.Errorf("expected unknown, got %q", txs[0].Category)
	}
}

func TestCSVParse_NoDescription(t *testing.T) {
	f, err := os.Open("testdata/no_description.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	txs, err := (&parser.CSVParser{}).Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
	if txs[0].Description != "" {
		t.Errorf("expected empty description, got %q", txs[0].Description)
	}
}
