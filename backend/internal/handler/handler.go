package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/terracodum/expensemind/backend/internal/service"
)

type Handler struct {
	svc service.Service
}

func New(svc service.Service) http.Handler {
	h := &Handler{svc: svc}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", h.getTransactions)
		r.Post("/upload", h.uploadTransactions)
		r.Patch("/{id}", h.updateTransaction)
		r.Delete("/{id}", h.deleteTransaction)
	})

	r.Route("/analytics", func(r chi.Router) {
		r.Post("/forecast", h.createForecastJob)
		r.Get("/forecast/{id}", h.getForecastJob)
	})

	r.Route("/recurring", func(r chi.Router) {
		r.Get("/", h.getRecurringRules)
		r.Post("/", h.saveRecurringRule)
		r.Delete("/{sourceID}", h.deleteRecurringRule)
	})

	return r
}
