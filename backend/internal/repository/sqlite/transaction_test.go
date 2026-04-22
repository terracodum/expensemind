package sqlite

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terracodum/expensemind/backend/internal/domain"
)

func setupTransactionDB(t *testing.T) *SQLiteTransactionRepository {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	repo, _, _, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return &repo
}

var testTx = domain.Transaction{
	Date:        time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
	Amount:      -500.0,
	Description: "Кофе",
	Category:    "Еда",
}

func TestTransaction_SaveAndFindAll(t *testing.T) {
	repo := setupTransactionDB(t)

	if err := repo.Save(testTx); err != nil {
		t.Fatal(err)
	}

	txs, err := repo.FindAll(domain.Filters{})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1, got %d", len(txs))
	}
	got := txs[0]
	if got.Amount != testTx.Amount {
		t.Errorf("amount: expected %v, got %v", testTx.Amount, got.Amount)
	}
	if got.Description != testTx.Description {
		t.Errorf("description: expected %q, got %q", testTx.Description, got.Description)
	}
	if got.Category != testTx.Category {
		t.Errorf("category: expected %q, got %q", testTx.Category, got.Category)
	}
	if !got.Date.Equal(testTx.Date) {
		t.Errorf("date: expected %v, got %v", testTx.Date, got.Date)
	}
}

func TestTransaction_SaveAll(t *testing.T) {
	repo := setupTransactionDB(t)

	txs := []domain.Transaction{
		testTx,
		{Date: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC), Amount: 1000.0, Description: "Зарплата", Category: "Доход"},
	}

	if err := repo.SaveAll(txs); err != nil {
		t.Fatal(err)
	}

	result, err := repo.FindAll(domain.Filters{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestTransaction_FindAll_FilterByCategory(t *testing.T) {
	repo := setupTransactionDB(t)

	repo.Save(testTx)
	repo.Save(domain.Transaction{
		Date:     time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC),
		Amount:   1000.0,
		Category: "Доход",
	})

	result, err := repo.FindAll(domain.Filters{Category: "Еда"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Category != "Еда" {
		t.Errorf("expected Еда, got %q", result[0].Category)
	}
}

func TestTransaction_FindAll_FilterByDateRange(t *testing.T) {
	repo := setupTransactionDB(t)

	repo.Save(domain.Transaction{Date: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Amount: -100, Category: "A"})
	repo.Save(domain.Transaction{Date: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), Amount: -200, Category: "B"})
	repo.Save(domain.Transaction{Date: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), Amount: -300, Category: "C"})

	result, err := repo.FindAll(domain.Filters{
		From: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Category != "B" {
		t.Errorf("expected B, got %q", result[0].Category)
	}
}

func TestTransaction_Update(t *testing.T) {
	repo := setupTransactionDB(t)

	repo.Save(testTx)
	txs, _ := repo.FindAll(domain.Filters{})
	saved := txs[0]

	saved.Amount = -999.0
	saved.Category = "Обновлено"
	if err := repo.Update(saved); err != nil {
		t.Fatal(err)
	}

	result, _ := repo.FindAll(domain.Filters{})
	if result[0].Amount != -999.0 {
		t.Errorf("expected -999.0, got %v", result[0].Amount)
	}
	if result[0].Category != "Обновлено" {
		t.Errorf("expected Обновлено, got %q", result[0].Category)
	}
}

func TestTransaction_Delete(t *testing.T) {
	repo := setupTransactionDB(t)

	repo.Save(testTx)
	txs, _ := repo.FindAll(domain.Filters{})
	id := txs[0].ID

	if err := repo.Delete(id); err != nil {
		t.Fatal(err)
	}

	result, _ := repo.FindAll(domain.Filters{})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
