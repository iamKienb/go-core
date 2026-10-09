package observability_v1

import (
	"context"
	"net/http"

	"github.com/iamKienb/go-core/shared"
	"go.opentelemetry.io/otel/trace"
)

func traceIDFromContext(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		return spanContext.TraceID().String()
	}

	if traceID := shared.GetTraceID(ctx); traceID != "" {
		return traceID
	}
	return ""
}

func setResponseMeta(ctx context.Context, meta http.Header) {
	if reqID := shared.GetRequestID(ctx); reqID != "" {
		meta.Set("x-request-id", reqID)
	}
	if traceID := traceIDFromContext(ctx); traceID != "" {
		meta.Set("x-trace-id", traceID)
	}
}
