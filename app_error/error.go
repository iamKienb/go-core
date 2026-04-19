package app_error

import (
	"errors"

	"connectrpc.com/connect"
)

type Kind int

const (
	KindValidation   Kind = iota // 400
	KindUnauthorized             // 401
	KindForbidden                // 403
	KindNotFound                 // 404
	KindConflict                 // 409
	KindInternal                 // 500
)

type AppError struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func Transform(kind Kind, msg string, err error) *AppError {
	return &AppError{
		Kind:    kind,
		Message: msg,
		Err:     err,
	}
}

func (k Kind) ConnectCode() connect.Code {
	return map[Kind]connect.Code{
		KindValidation:   connect.CodeInvalidArgument,
		KindUnauthorized: connect.CodeUnauthenticated,
		KindForbidden:    connect.CodePermissionDenied,
		KindNotFound:     connect.CodeNotFound,
		KindConflict:     connect.CodeAlreadyExists,
	}[k]
}

func Validation(msg string) *AppError {
	return &AppError{Kind: KindValidation, Message: msg}
}

func Conflict(msg string, err error) *AppError {
	return &AppError{Kind: KindConflict, Message: msg, Err: err}
}

func NotFound(msg string, err error) *AppError {
	return &AppError{Kind: KindNotFound, Message: msg, Err: err}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Kind: KindUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Kind: KindForbidden, Message: msg}
}

func Internal(err error) *AppError {
	return &AppError{Kind: KindInternal, Message: "internal server error", Err: err}
}

func IsKind(err error, kind Kind) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Kind == kind
	}
	return false
}
