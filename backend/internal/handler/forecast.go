package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	apperrors "github.com/terracodum/expensemind/backend/internal/errors"
)

func (h *Handler) createForecastJob(w http.ResponseWriter, r *http.Request) {
	id, err := h.svc.CreateForecastJob()
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]int{"job_id": id})
}

func (h *Handler) getForecastJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.svc.GetForecastJobs()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

func (h *Handler) getForecastJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, apperrors.ValidationError("invalid id"))
		return
	}

	job, err := h.svc.GetForecastJob(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, job)
}
