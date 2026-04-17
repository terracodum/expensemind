package errors

type ErrorCode string

const (
	INTERNAL_ERROR         ErrorCode = "INTERNAL_ERROR"
	VALIDATION_ERROR       ErrorCode = "VALIDATION_ERROR"
	NOT_FOUND              ErrorCode = "NOT_FOUND"
	PARSE_ERROR            ErrorCode = "PARSE_ERROR"
	INVALID_PDF_FORMAT     ErrorCode = "INVALID_PDF_FORMAT"
	INVALID_CSV_FORMAT     ErrorCode = "INVALID_CSV_FORMAT"
	ML_SERVICE_UNAVAILABLE ErrorCode = "ML_SERVICE_UNAVAILABLE"
	ML_RESPONSE_INVALID    ErrorCode = "ML_RESPONSE_INVALID"
	DB_ERROR               ErrorCode = "DB_ERROR"
)
