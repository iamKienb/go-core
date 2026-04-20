package observability

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"
	"github.com/iamKienb/shopify-go-platform/utils"
)

func RecoveryInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					traceID := utils.ExtractTraceID(ctx)

					logger.ErrorContext(ctx, "PANIC_RECOVERED",
						slog.String("trace_id", traceID),
						slog.Any("panic_info", r),
						slog.String("procedure", req.Spec().Procedure),
						slog.String("stack", string(stack)),
					)

					err = connect.NewError(
						connect.CodeInternal,
						errors.New("INTERNAL_SERVER_ERROR"),
					)

					// err.(*connect.Error).Meta().Set("x-error-type", "panic")
				}
			}()

			return next(ctx, req)
		}
	}
}
