package observabilityx

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	authx "github.com/iamKienb/go-core/middleware/auth"
	"github.com/iamKienb/go-core/utils"
)

func LoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			resp, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				return resp, err
			}

			if duration < 500*time.Millisecond {
				logger.DebugContext(ctx, "request completed",
					slog.String("req_id", authx.GetRequestID(ctx)),
					slog.String("trace_id", utils.ExtractTraceID(ctx)),
					slog.String("method", req.Spec().Procedure),
					slog.Duration("latency", duration),
				)
				return resp, nil
			}

			logger.WarnContext(ctx, "slow request",
				slog.String("req_id", authx.GetRequestID(ctx)),
				slog.String("trace_id", utils.ExtractTraceID(ctx)),
				slog.String("method", req.Spec().Procedure),
				slog.String("status", "ok"),
				slog.Duration("latency", duration),
				slog.String("ip", req.Peer().Addr),
			)

			return resp, err
		}
	}
}
