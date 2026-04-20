package sqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terracodum/expensemind/backend/internal/domain"
)

func setupForecastJobDB(t *testing.T) *SQLiteForecastJobRepository {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, _, repo, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return &repo
}

func TestForecastJob_Create(t *testing.T) {
	repo := setupForecastJobDB(t)

	id, err := repo.Create()
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
}

func TestForecastJob_FindByID_Pending(t *testing.T) {
	repo := setupForecastJobDB(t)

	id, _ := repo.Create()
	job, err := repo.FindByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "pending" {
		t.Errorf("expected pending, got %q", job.Status)
	}
	if job.Result != nil {
		t.Errorf("expected nil result, got %v", job.Result)
	}
	if job.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
}

func TestForecastJob_Update_WithResult(t *testing.T) {
	repo := setupForecastJobDB(t)

	id, _ := repo.Create()

	forecast := &domain.Forecast{
		Points:           []domain.Point{{T: 1, Balance: 100.0}, {T: 2, Balance: 200.0}},
		PredictedBalance: 300.0,
		Confidence:       0.85,
	}

	err := repo.Update(domain.ForecastJob{ID: id, Status: "done", Result: forecast})
	if err != nil {
		t.Fatal(err)
	}

	job, err := repo.FindByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "done" {
		t.Errorf("expected done, got %q", job.Status)
	}
	if job.Result == nil {
		t.Fatal("expected non-nil result")
	}
	if job.Result.PredictedBalance != 300.0 {
		t.Errorf("expected 300.0, got %v", job.Result.PredictedBalance)
	}
	if job.Result.Confidence != 0.85 {
		t.Errorf("expected 0.85, got %v", job.Result.Confidence)
	}
	if len(job.Result.Points) != 2 {
		t.Errorf("expected 2 points, got %d", len(job.Result.Points))
	}
}

func TestForecastJob_Update_Failed(t *testing.T) {
	repo := setupForecastJobDB(t)

	id, _ := repo.Create()

	err := repo.Update(domain.ForecastJob{ID: id, Status: "failed", Result: nil})
	if err != nil {
		t.Fatal(err)
	}

	job, err := repo.FindByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "failed" {
		t.Errorf("expected failed, got %q", job.Status)
	}
	if job.Result != nil {
		t.Errorf("expected nil result, got %v", job.Result)
	}
}

func TestForecastJob_FindAll(t *testing.T) {
	repo := setupForecastJobDB(t)

	repo.Create()
	repo.Create()

	jobs, err := repo.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 2 {
		t.Fatalf("expected 2, got %d", len(jobs))
	}
}
