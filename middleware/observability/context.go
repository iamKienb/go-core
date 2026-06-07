package observabilityx

import (
	"context"
	"net/http"

	authx "github.com/iamKienb/go-core/middleware/auth"
	"github.com/iamKienb/go-core/middleware/shared"
)

func traceIDFromContext(ctx context.Context) string {
	if traceID := authx.GetTraceID(ctx); traceID != "" {
		return traceID
	}

	return shared.ExtractTraceID(ctx)
}

func setResponseMeta(ctx context.Context, meta http.Header) {
	if reqID := authx.GetRequestID(ctx); reqID != "" {
		meta.Set("x-request-id", reqID)
	}
	if traceID := traceIDFromContext(ctx); traceID != "" {
		meta.Set("x-trace-id", traceID)
	}
}
