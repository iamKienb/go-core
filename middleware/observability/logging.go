package observability

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/iamKienb/shopify-go-platform/middleware/auth"
	"github.com/iamKienb/shopify-go-platform/utils"
)

func LoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			resp, err := next(ctx, req)

			duration := time.Since(start)
			code := connect.CodeOf(err)

			level := slog.LevelInfo
			if err != nil {
				if code == connect.CodeInternal || code == connect.CodeUnknown || code == connect.CodeUnavailable {
					level = slog.LevelError
				} else {
					level = slog.LevelWarn
				}
			}

			attrs := []any{
				slog.String("req_id", auth.GetRequestID(ctx)),
				slog.String("trace_id", utils.ExtractTraceID(ctx)),
				slog.String("method", req.Spec().Procedure),
				slog.String("status", code.String()),
				slog.Duration("latency", duration),
				slog.String("ip", req.Peer().Addr),
			}

			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			logger.Log(ctx, level, "gRPC Request Completed", attrs...)

			return resp, err
		}
	}
}
