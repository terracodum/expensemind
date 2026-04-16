package errors

func New(code ErrorCode, message string) AppError {
	return &appError{code: code, message: message}
}

func Wrap(code ErrorCode, message string, err error) AppError {
	return &appError{code, message, err}
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

func PdfParseError(message string) AppError {
	return New(PDF_PARSE_ERROR, message)
}

func PdfInvalidError(message string) AppError {
	return New(PDF_INVALID_FORMAT, message)
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
