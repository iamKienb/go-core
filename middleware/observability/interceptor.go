package observabilityx

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	app_error "github.com/iamKienb/shopify-go-platform/app_error"
)

func ErrorResponseInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				return nil, connectErr
			}

			var appErr *app_error.AppError
			if errors.As(err, &appErr) {
				cErr := connect.NewError(appErr.Kind.ConnectCode(), errors.New(appErr.Message))
				cErr.Meta().Set("x-error-code", appErr.Kind.ConnectCode().String())
				return nil, cErr
			}

			return nil, connect.NewError(connect.CodeInternal, errors.New("INTERNAL_SERVER_ERROR"))
		}
	}
}

type validator interface {
	Validate() error
}

func ValidationRequestInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if v, ok := req.Any().(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}
			}
			return next(ctx, req)
		}
	}
}
