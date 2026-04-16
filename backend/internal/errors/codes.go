package errors

type ErrorCode string

const (
	INTERNAL_ERROR         ErrorCode = "INTERNAL_ERROR"
	VALIDATION_ERROR       ErrorCode = "VALIDATION_ERROR"
	NOT_FOUND              ErrorCode = "NOT_FOUND"
	CSV_PARSE_ERROR        ErrorCode = "CSV_PARSE_ERROR"
	CSV_INVALID_FORMAT     ErrorCode = "CSV_INVALID_FORMAT"
	ML_SERVICE_UNAVAILABLE ErrorCode = "ML_SERVICE_UNAVAILABLE"
	ML_RESPONSE_INVALID    ErrorCode = "ML_RESPONSE_INVALID"
	DB_ERROR               ErrorCode = "DB_ERROR"
)
