package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/terracodum/expensemind/backend/internal/domain"
	apperrors "github.com/terracodum/expensemind/backend/internal/errors"
)

func (h *Handler) uploadTransactions(w http.ResponseWriter, r *http.Request) {
	err := h.svc.UploadTransactions(r.Header.Get("Content-Type"), r.Body)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getTransactions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filters := domain.Filters{
		Category: q.Get("category"),
	}

	if v := q.Get("from"); v != "" {
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			writeError(w, apperrors.ValidationError("invalid 'from' date, expected YYYY-MM-DD"))
			return
		}
		filters.From = t
	}

	if v := q.Get("to"); v != "" {
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			writeError(w, apperrors.ValidationError("invalid 'to' date, expected YYYY-MM-DD"))
			return
		}
		filters.To = t
	}

	txs, err := h.svc.GetTransactions(filters)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, txs)
}

func (h *Handler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, apperrors.ValidationError("invalid id"))
		return
	}

	var tx domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		writeError(w, apperrors.ValidationError("invalid request body"))
		return
	}
	tx.ID = id

	if err := h.svc.UpdateTransaction(tx); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, apperrors.ValidationError("invalid id"))
		return
	}

	if err := h.svc.DeleteTransaction(id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
