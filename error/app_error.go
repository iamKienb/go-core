package app_error

type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func New(httpStatus, code int, msg string) AppError {
	return AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    msg,
	}
}

var (
	ErrBadRequest      = New(400, 1000, "Bad request")
	ErrInternal        = New(500, 9999, "Internal server error")
	ErrTooManyRequests = New(429, 1003, "Too many requests")
	ErrUnauthorized    = New(401, 1001, "Unauthorized")
)

func (e AppError) Error() string {
	return e.Message
}
