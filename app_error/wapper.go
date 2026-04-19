package app_error

import "errors"

type ErrorMapping struct {
	Kind Kind
	Msg  string
}

type ServiceErrorMap map[error]ErrorMapping

func WrapError(err error, errMap ServiceErrorMap) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return err
	}

	for sentinel, mapping := range errMap {
		if errors.Is(err, sentinel) {
			return Transform(mapping.Kind, mapping.Msg, err)
		}
	}

	return Internal(err)
}
