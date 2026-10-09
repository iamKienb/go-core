package observabilityx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func GetTraceIDFromContext(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)

	if spanContext.IsValid() {
		return spanContext.TraceID().String()
	}

	return ""
}
