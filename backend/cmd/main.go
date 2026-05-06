package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terracodum/expensemind/backend/internal/handler"
	"github.com/terracodum/expensemind/backend/internal/ml"
	"github.com/terracodum/expensemind/backend/internal/parser"
	"github.com/terracodum/expensemind/backend/internal/repository/sqlite"
	"github.com/terracodum/expensemind/backend/internal/service"
)

func main() {
	dbPath := envOr("DB_PATH", "expensemind.db")
	mlURL := envOr("ML_URL", "http://localhost:8001")
	addr := envOr("ADDR", ":8080")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("failed to open db", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	txRepo, recurRepo, forecastRepo, err := sqlite.New(db)
	if err != nil {
		slog.Error("failed to init db", "err", err)
		os.Exit(1)
	}

	svc := service.New(parser.Factory{}, &txRepo, &recurRepo, &forecastRepo, ml.New(mlURL))
	h := handler.New(svc)

	slog.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
