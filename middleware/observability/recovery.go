package observability

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"
	"github.com/iamKienb/shopify-go-platform/middleware/auth"
	"github.com/iamKienb/shopify-go-platform/utils"
)

func RecoveryInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					reqID := auth.GetRequestID(ctx)
					traceID := utils.ExtractTraceID(ctx)

					logger.ErrorContext(ctx, "CRITICAL_PANIC_RECOVERED",
						slog.String("req_id", reqID),
						slog.String("trace_id", traceID),
						slog.String("method", req.Spec().Procedure),
						slog.Any("panic", r),
						slog.String("stack", string(debug.Stack())),
					)

					err = connect.NewError(connect.CodeInternal, errors.New("INTERNAL_SERVER_ERROR"))
				}
			}()

			return next(ctx, req)
		}
	}
}
