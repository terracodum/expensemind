package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/terracodum/expensemind/backend/internal/domain"
	apperrors "github.com/terracodum/expensemind/backend/internal/errors"
)

func (h *Handler) getRecurringRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.svc.GetRecurringRules()
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, rules)
}

func (h *Handler) saveRecurringRule(w http.ResponseWriter, r *http.Request) {
	var rule domain.RecurringRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeError(w, apperrors.ValidationError("invalid request body"))
		return
	}

	if err := h.svc.SaveRecurringRule(rule); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteRecurringRule(w http.ResponseWriter, r *http.Request) {
	sourceID := chi.URLParam(r, "sourceID")
	if sourceID == "" {
		writeError(w, apperrors.ValidationError("missing sourceID"))
		return
	}

	if err := h.svc.DeleteRecurringRule(sourceID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
