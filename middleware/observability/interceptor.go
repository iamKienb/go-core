package observabilityx

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"connectrpc.com/connect"
	kafkax "github.com/iamKienb/go-core/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
)

func WithServerInterceptors(serviceName string) (connect.Option, error) {
	meterName := fmt.Sprintf("%s.connect-server", serviceName)
	meter := otel.GetMeterProvider().Meter(meterName)

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

			tracer := otel.Tracer(serviceName)
			apiMethod := req.Spec().Procedure

			ctx, span := tracer.Start(ctx, apiMethod)
			defer span.End()

			start := time.Now()

			defer func() {
				duration := time.Since(start)
				ms := float64(duration) / float64(time.Millisecond)
				statusLabel := "ok"

				if r := recover(); r != nil {
					statusLabel = "panic"
					panicErr := fmt.Errorf("rpc panic recovered: %v", r)
					stackTrace := debug.Stack()

					span.RecordError(panicErr)
					span.SetStatus(codes.Error, panicErr.Error())

					slog.ErrorContext(ctx, "CRITICAL_RPC_PANIC_DETECTED",
						slog.String("method", apiMethod),
						slog.Any("error", panicErr),
						slog.String("stack", string(stackTrace)),
					)

					err = connect.NewError(connect.CodeInternal, fmt.Errorf("an internal server error occurred"))
				} else if err != nil {
					statusLabel = "error"

					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())

					slog.ErrorContext(ctx, "rpc request failed",
						slog.String("method", apiMethod),
						slog.Duration("latency", duration),
						slog.Any("error", err),
					)
				} else {
					span.SetStatus(codes.Ok, "success")

					slog.InfoContext(ctx, "rpc request completed",
						slog.String("method", apiMethod),
						slog.Duration("latency", duration),
						slog.String("status", statusLabel),
					)
				}

				requestCounter.Add(ctx, 1, metric.WithAttributes(
					attribute.String("method", apiMethod),
					attribute.String("status", statusLabel),
				))
				latencyHistogram.Record(ctx, ms, metric.WithAttributes(
					attribute.String("method", apiMethod),
					attribute.String("status", statusLabel),
				))
			}()

			resp, err = next(ctx, req)
			return resp, err
		}
	})

	return connect.WithInterceptors(interceptor), nil
}

func WithWorkerInterceptors(serviceName string, next kafkax.WorkerHandler) (kafkax.WorkerHandler, error) {
	meterName := fmt.Sprintf("%s.kafka-worker", serviceName)
	meter := otel.GetMeterProvider().Meter(meterName)

	requestCounter, err := meter.Int64Counter("kafka_worker_messages_total",
		metric.WithDescription("Total Kafka messages processed by worker, partitioned by topic and status"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create worker message request counter: %w", err)
	}

	latencyHistogram, err := meter.Float64Histogram("kafka_worker_duration_milliseconds",
		metric.WithDescription("Kafka message processing duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create worker duration histogram: %w", err)
	}

	interceptor := WorkerHandlerFunc(func(ctx context.Context, msg kafkax.Message) error {
		topic := msg.Topic
		tracer := otel.Tracer(serviceName)

		ctx, span := tracer.Start(ctx, fmt.Sprintf("kafka.process.%s", topic))
		defer span.End()

		start := time.Now()
		var handlerErr error

		defer func() {
			duration := time.Since(start)
			ms := float64(duration) / float64(time.Millisecond)
			statusLabel := "ok"

			if r := recover(); r != nil {
				statusLabel = "panic"
				panicErr := fmt.Errorf("worker panic recovered: %v", r)
				stackTrace := debug.Stack()

				span.RecordError(panicErr)
				span.SetStatus(codes.Error, panicErr.Error())

				slog.ErrorContext(ctx, "CRITICAL_WORKER_PANIC_DETECTED",
					slog.String("topic", topic),
					slog.Any("error", panicErr),
					slog.String("stack", string(stackTrace)),
				)

				handlerErr = panicErr
			} else if handlerErr != nil {
				statusLabel = "error"

				span.RecordError(handlerErr)
				span.SetStatus(codes.Error, handlerErr.Error())

				slog.ErrorContext(ctx, "kafka worker message processing failed",
					slog.String("topic", topic),
					slog.Duration("latency", duration),
					slog.Any("error", handlerErr),
				)
			} else {
				span.SetStatus(codes.Ok, "success")

				slog.InfoContext(ctx, "kafka worker message processing completed",
					slog.String("topic", topic),
					slog.Duration("latency", duration),
					slog.String("status", statusLabel),
				)
			}

			requestCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("status", statusLabel),
			))
			latencyHistogram.Record(ctx, ms, metric.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("status", statusLabel),
			))
		}()

		handlerErr = next.Handle(ctx, msg)
		return handlerErr
	})

	return interceptor, nil
}

type WorkerHandlerFunc func(ctx context.Context, msg kafkax.Message) error

func (f WorkerHandlerFunc) Handle(ctx context.Context, msg kafkax.Message) error {
	return f(ctx, msg)
}
