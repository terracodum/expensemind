package errors

func New(code ErrorCode, message string) AppError {
	return &appError{code: code, message: message}
}

func Wrap(code ErrorCode, message string, err error) AppError {
	return &appError{code: code, message: message, wrapped: err}
}

func NotFound(message string) AppError {
	return New(NOT_FOUND, message)
}

func InternalError(message string) AppError {
	return New(INTERNAL_ERROR, message)
}

func ValidationError(message string) AppError {
	return New(VALIDATION_ERROR, message)
}

func ParseError(message string) AppError {
	return New(PARSE_ERROR, message)
}

func InvalidPDFFormat(message string) AppError {
	return New(INVALID_PDF_FORMAT, message)
}

func InvalidCSVFormat(message string) AppError {
	return New(INVALID_CSV_FORMAT, message)
}

func MLServiceUnavailable(message string, err error) AppError {
	return Wrap(ML_SERVICE_UNAVAILABLE, message, err)
}

func MLResponseInvalid(message string, err error) AppError {
	return Wrap(ML_RESPONSE_INVALID, message, err)
}

func DBError(message string, err error) AppError {
	return Wrap(DB_ERROR, message, err)
}
