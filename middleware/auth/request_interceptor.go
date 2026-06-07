package authx

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/iamKienb/go-core/middleware/shared"
)

func RequestContextInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			reqID := strings.TrimSpace(req.Header().Get(HeaderRequestID))
			if reqID == "" {
				reqID = "req-" + uuid.NewString()
			}

			traceID := strings.TrimSpace(req.Header().Get(HeaderTraceID))
			if traceID == "" {
				traceID = shared.ExtractTraceID(ctx)
			}
			if traceID == "" {
				traceID = strings.ReplaceAll(uuid.NewString(), "-", "")
			}

			req.Header().Set(HeaderRequestID, reqID)
			req.Header().Set(HeaderTraceID, traceID)
			ctx = SetRequestIDToCtx(ctx, reqID)
			ctx = SetTraceIDToCtx(ctx, traceID)

			return next(ctx, req)
		}
	}
}
