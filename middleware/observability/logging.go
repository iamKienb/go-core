package observability

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/iamKienb/shopify-go-platform/utils"
)

func LoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			resp, err := next(ctx, req)

			level := slog.LevelInfo
			code := connect.CodeOf(err)

			if err != nil {
				if code >= connect.CodeUnknown {
					level = slog.LevelError
				} else {
					level = slog.LevelWarn
				}
			}

			traceID := utils.ExtractTraceID(ctx)
			attrs := []slog.Attr{
				slog.String("trace_id", traceID),
				slog.String("procedure", req.Spec().Procedure),
				slog.String("status", code.String()),
				slog.Duration("duration", time.Since(start)),
				slog.String("protocol", req.Peer().Protocol),
			}

			if err != nil {
				attrs = append(attrs, slog.Any("error.detail", err))
			}

			logger.LogAttrs(ctx, level, "GRPC_REQUEST_COMPLETED", attrs...)

			return resp, err
		}
	}
}
