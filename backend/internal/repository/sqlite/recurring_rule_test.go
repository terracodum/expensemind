package sqlite

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terracodum/expensemind/backend/internal/domain"
)

func setupRecurringRuleDB(t *testing.T) *SQLiteRecurringRuleRepository {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:?_loc=UTC")
	if err != nil {
		t.Fatal(err)
	}
	_, repo, _, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return &repo
}

func TestFindActive_NoRules(t *testing.T) {
	repo := setupRecurringRuleDB(t)
	today := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)

	result, err := repo.FindActive(today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(result))
	}
}

func TestFindActive_OneRule(t *testing.T) {
	repo := setupRecurringRuleDB(t)
	today := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)

	rule := domain.RecurringRule{
		SourceID:  "src-1",
		Type:      "expense",
		Amount:    100.0,
		Day:       15,
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Label:     "rent",
	}
	if err := repo.Save(rule); err != nil {
		t.Fatal(err)
	}

	result, err := repo.FindActive(today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(result))
	}
	if result[0].SourceID != "src-1" {
		t.Errorf("expected source_id src-1, got %s", result[0].SourceID)
	}
	if result[0].Amount != 100.0 {
		t.Errorf("expected amount 100.0, got %f", result[0].Amount)
	}
}

func TestFindActive_TwoRulesSameSourceID(t *testing.T) {
	repo := setupRecurringRuleDB(t)
	today := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)

	older := domain.RecurringRule{
		SourceID:  "src-1",
		Type:      "expense",
		Amount:    100.0,
		Day:       15,
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Label:     "rent",
	}
	newer := domain.RecurringRule{
		SourceID:  "src-1",
		Type:      "expense",
		Amount:    200.0,
		Day:       15,
		StartDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		Label:     "rent updated",
	}

	if err := repo.Save(older); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(newer); err != nil {
		t.Fatal(err)
	}

	result, err := repo.FindActive(today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 rule (latest by start_date), got %d", len(result))
	}
	if result[0].Amount != 200.0 {
		t.Errorf("expected amount 200.0 (newer rule), got %f", result[0].Amount)
	}
}

func TestFindActive_FutureRuleExcluded(t *testing.T) {
	repo := setupRecurringRuleDB(t)
	today := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)

	future := domain.RecurringRule{
		SourceID:  "src-1",
		Type:      "expense",
		Amount:    300.0,
		Day:       1,
		StartDate: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		Label:     "future",
	}
	if err := repo.Save(future); err != nil {
		t.Fatal(err)
	}

	result, err := repo.FindActive(today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 rules (future excluded), got %d", len(result))
	}
}
