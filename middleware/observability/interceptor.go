package observabilityx

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	app_error "github.com/iamKienb/shopify-go-platform/app_error"
	authx "github.com/iamKienb/shopify-go-platform/middleware/auth"
	"github.com/iamKienb/shopify-go-platform/utils"
)

func ErrorResponseInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				if reqID := authx.GetRequestID(ctx); reqID != "" {
					connectErr.Meta().Set("x-request-id", reqID)
				}
				return nil, connectErr
			}

			appErr := app_error.From(err)
			if logger != nil {
				logErrorResponse(ctx, logger, req, appErr)
			}

			cErr := connect.NewError(appErr.Kind.ConnectCode(), errors.New(appErr.PublicMessage()))
			cErr.Meta().Set("x-error-code", appErr.PublicCode())
			cErr.Meta().Set("x-error-kind", appErr.Kind.ConnectCode().String())
			if reqID := authx.GetRequestID(ctx); reqID != "" {
				cErr.Meta().Set("x-request-id", reqID)
			}
			return nil, cErr
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
					return nil, app_error.New(app_error.KindValidation, "request_invalid", err.Error(), err)
				}
			}
			return next(ctx, req)
		}
	}
}

func logErrorResponse(ctx context.Context, logger *slog.Logger, req connect.AnyRequest, appErr *app_error.AppError) {
	if logger == nil || appErr == nil {
		return
	}

	level := slog.LevelWarn
	if appErr.Kind == app_error.KindInternal || appErr.Kind == app_error.KindUnavailable {
		level = slog.LevelError
	}

	attrs := []any{
		slog.String("req_id", authx.GetRequestID(ctx)),
		slog.String("trace_id", utils.ExtractTraceID(ctx)),
		slog.String("method", req.Spec().Procedure),
		slog.String("public_code", appErr.PublicCode()),
		slog.String("public_message", appErr.PublicMessage()),
		slog.String("connect_code", appErr.Kind.ConnectCode().String()),
	}

	if appErr.Err != nil {
		attrs = append(attrs, slog.String("cause", appErr.Err.Error()))
	}

	logger.Log(ctx, level, "request failed", attrs...)
}
