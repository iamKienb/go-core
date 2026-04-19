package middleware

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	app_error "github.com/iamKienb/shopify-go-platform/error"
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
				if appErr.Kind == app_error.KindInternal {
					fmt.Printf("[SERVER-ERROR] Detail: %v\n", appErr)
				}
				return nil, connect.NewError(appErr.Kind.ConnectCode(), errors.New(appErr.Message))
			}

			fmt.Printf("[UNHANDLED-ERROR] %v\n", err)
			return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
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
					return nil, app_error.Validation(err.Error())
				}
			}

			return next(ctx, req)
		}
	}
}
