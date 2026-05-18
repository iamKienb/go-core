package observabilityx

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"
	app_error "github.com/iamKienb/go-core/app_error"
	authx "github.com/iamKienb/go-core/middleware/auth"
	"github.com/iamKienb/go-core/utils"
)

func RecoveryInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					reqID := authx.GetRequestID(ctx)
					traceID := utils.ExtractTraceID(ctx)

					logger.ErrorContext(ctx, "CRITICAL_PANIC_RECOVERED",
						slog.String("req_id", reqID),
						slog.String("trace_id", traceID),
						slog.String("method", req.Spec().Procedure),
						slog.Any("panic", r),
						slog.String("stack", string(debug.Stack())),
					)

					err = app_error.Internal(errors.New("panic recovered"))
				}
			}()

			return next(ctx, req)
		}
	}
}
