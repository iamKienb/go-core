package app_error

import (
	"errors"
	"strings"

	"connectrpc.com/connect"
)

type Kind int

const (
	KindValidation   Kind = iota // 400
	KindUnauthorized             // 401
	KindForbidden                // 403
	KindNotFound                 // 404
	KindConflict                 // 409
	KindUnavailable              // 503
	KindInternal                 // 500
)

type AppError struct {
	Kind    Kind
	Code    string
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
		Code:    kind.DefaultCode(),
		Message: msg,
		Err:     err,
	}
}

func New(kind Kind, code, msg string, err error) *AppError {
	if code == "" {
		code = kind.DefaultCode()
	}
	if msg == "" {
		msg = kind.DefaultMessage()
	}

	return &AppError{
		Kind:    kind,
		Code:    normalizeCode(code),
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
		KindUnavailable:  connect.CodeUnavailable,
		KindInternal:     connect.CodeInternal,
	}[k]
}

func (k Kind) DefaultCode() string {
	switch k {
	case KindValidation:
		return "validation_error"
	case KindUnauthorized:
		return "unauthorized"
	case KindForbidden:
		return "forbidden"
	case KindNotFound:
		return "not_found"
	case KindConflict:
		return "conflict"
	case KindUnavailable:
		return "service_unavailable"
	default:
		return "internal_server_error"
	}
}

func (k Kind) DefaultMessage() string {
	switch k {
	case KindValidation:
		return "validation failed"
	case KindUnauthorized:
		return "authentication required"
	case KindForbidden:
		return "you do not have permission to perform this action"
	case KindNotFound:
		return "resource not found"
	case KindConflict:
		return "resource already exists"
	case KindUnavailable:
		return "service temporarily unavailable"
	default:
		return "internal server error"
	}
}

func (e *AppError) PublicCode() string {
	if e.Code != "" {
		return normalizeCode(e.Code)
	}
	return e.Kind.DefaultCode()
}

func (e *AppError) PublicMessage() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Kind.DefaultMessage()
}

func Validation(msg string) *AppError {
	return New(KindValidation, "", msg, nil)
}

func Conflict(msg string, err error) *AppError {
	return New(KindConflict, "", msg, err)
}

func NotFound(msg string, err error) *AppError {
	return New(KindNotFound, "", msg, err)
}

func Unauthorized(msg string) *AppError {
	return New(KindUnauthorized, "", msg, nil)
}

func Forbidden(msg string) *AppError {
	return New(KindForbidden, "", msg, nil)
}

func Unavailable(msg string, err error) *AppError {
	return New(KindUnavailable, "", msg, err)
}

func Internal(err error) *AppError {
	return New(KindInternal, "", "", err)
}

func From(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal(err)
}

func IsKind(err error, kind Kind) bool {
	if appErr := From(err); appErr != nil {
		return appErr.Kind == kind
	}
	return false
}

func normalizeCode(code string) string {
	normalized := strings.TrimSpace(strings.ToLower(code))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	return normalized
}
