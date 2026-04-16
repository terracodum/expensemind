package errors

type appError struct {
	code    ErrorCode
	message string
	wrapped error
}

func (e *appError) Code() ErrorCode {
	return e.code
}

func (e *appError) Error() string {
	return e.message
}

func (e *appError) Unwrap() error {
	return e.wrapped
}
