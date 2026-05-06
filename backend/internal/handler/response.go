package handler

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/terracodum/expensemind/backend/internal/errors"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var code apperrors.ErrorCode
	var message string

	if appErr, ok := err.(apperrors.AppError); ok {
		code = appErr.Code()
		message = appErr.Error()
	} else {
		code = apperrors.INTERNAL_ERROR
		message = "internal error"
	}

	status := errorStatus(code)
	writeJSON(w, status, map[string]string{"error": string(code), "message": message})
}

func errorStatus(code apperrors.ErrorCode) int {
	switch code {
	case apperrors.NOT_FOUND:
		return http.StatusNotFound
	case apperrors.VALIDATION_ERROR:
		return http.StatusBadRequest
	case apperrors.PARSE_ERROR, apperrors.INVALID_PDF_FORMAT, apperrors.INVALID_CSV_FORMAT:
		return http.StatusUnprocessableEntity
	case apperrors.ML_SERVICE_UNAVAILABLE:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
