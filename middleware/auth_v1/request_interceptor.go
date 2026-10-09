package auth_v1

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/iamKienb/go-core/shared"
	"go.opentelemetry.io/otel/trace"
)

func RequestContextInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			traceID := strings.TrimSpace(req.Header().Get(HeaderTraceID))

			if traceID == "" {
				spanContext := trace.SpanContextFromContext(ctx)
				if spanContext.IsValid() {
					traceID = spanContext.TraceID().String()
				}
			}

			if traceID == "" {
				traceID = strings.ReplaceAll(uuid.NewString(), "-", "")
			}

			req.Header().Set(HeaderTraceID, traceID)
			ctx = shared.SetTraceIDToCtx(ctx, traceID)

			return next(ctx, req)
		}
	}
}
