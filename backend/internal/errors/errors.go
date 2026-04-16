package errors

type AppError interface {
	Code() ErrorCode
	Error() string
	Unwrap() error
}
