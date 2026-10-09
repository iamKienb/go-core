package observabilityx

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
)

func WithServerInterceptors() (connect.Option, error) {
	meter := otel.GetMeterProvider().Meter("connect-server-metrics")

	requestCounter, err := meter.Int64Counter("rpc_server_requests_total",
		metric.WithDescription("Total RPC requests processed, partitioned by method and status"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request counter: %w", err)
	}

	latencyHistogram, err := meter.Float64Histogram("rpc_server_duration_milliseconds",
		metric.WithDescription("RPC request processing latency in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create latency histogram: %w", err)
	}

	interceptor := connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header()))

			tracer := otel.Tracer("connect-server")
			apiMethod := req.Spec().Procedure

			ctx, span := tracer.Start(ctx, apiMethod)

			defer span.End()

			defer func() {
				if r := recover(); r != nil {
					panicErr := fmt.Errorf("panic recovered: %v", r)
					stackTrace := debug.Stack()

					span.RecordError(panicErr)
					span.SetStatus(codes.Error, panicErr.Error())

					slog.ErrorContext(ctx, "CRITICAL_RPC_PANIC_DETECTED",
						slog.String("method", apiMethod),
						slog.Any("error", panicErr),
						slog.String("stack", string(stackTrace)),
					)

					err = connect.NewError(connect.CodeInternal, fmt.Errorf("an internal server error occurred"))
				}
			}()

			start := time.Now()
			resp, err = next(ctx, req)
			duration := time.Since(start)
			ms := float64(duration) / float64(time.Millisecond)

			statusLabel := "ok"

			if err != nil {
				statusLabel = "error"

				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				slog.ErrorContext(ctx, "rpc request failed",
					slog.String("method", apiMethod),
					slog.Duration("latency", duration),
					slog.Any("error", err),
				)

				requestCounter.Add(ctx, 1, metric.WithAttributes(
					attribute.String("method", apiMethod),
					attribute.String("status", statusLabel),
				))
				latencyHistogram.Record(ctx, ms, metric.WithAttributes(
					attribute.String("method", apiMethod),
					attribute.String("status", statusLabel),
				))

				return resp, err
			}

			span.SetStatus(codes.Ok, "success")

			slog.InfoContext(ctx, "rpc request completed",
				slog.String("method", apiMethod),
				slog.Duration("latency", duration),
				slog.String("status", statusLabel),
			)

			requestCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("method", apiMethod),
				attribute.String("status", statusLabel),
			))
			latencyHistogram.Record(ctx, ms, metric.WithAttributes(
				attribute.String("method", apiMethod),
				attribute.String("status", statusLabel),
			))

			return resp, nil
		}
	})

	return connect.WithInterceptors(interceptor), nil
}
