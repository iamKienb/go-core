package app_error

import "errors"

type ErrorMapping struct {
	Kind Kind
	Code string
	Msg  string
}

type ServiceErrorMap map[error]ErrorMapping

func WrapError(err error, errMap ServiceErrorMap) error {
	if err == nil {
		return nil
	}

	if appErr := From(err); appErr != nil && errors.As(err, &appErr) {
		return appErr
	}

	for sentinel, mapping := range errMap {
		if errors.Is(err, sentinel) {
			code := mapping.Code
			if code == "" {
				code = sentinel.Error()
			}
			return New(mapping.Kind, code, mapping.Msg, err)
		}
	}

	return Internal(err)
}
