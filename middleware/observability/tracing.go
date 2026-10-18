package observabilityx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetTraceIDFromContext(ctx context.Context) string {
	carrier := propagation.MapCarrier{}

	otel.GetTextMapPropagator().Inject(ctx, carrier)

	if val, ok := carrier["traceparent"]; ok {
		return val
	}

	return ""
}
