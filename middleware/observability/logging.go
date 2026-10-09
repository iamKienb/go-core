package observabilityx

import (
	"context"
	"log/slog"
	"os"
)

type OtelHandler struct {
	next slog.Handler
}

func (h *OtelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &OtelHandler{h.next.WithAttrs(attrs)}
}

func (h *OtelHandler) WithGroup(name string) slog.Handler {
	return &OtelHandler{h.next.WithGroup(name)}
}

func (h *OtelHandler) Handle(ctx context.Context, record slog.Record) error {
	traceID := GetTraceIDFromContext(ctx)
	if traceID != "" {
		record.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.next.Handle(ctx, record)
}

func InitLogger() {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	otelHandler := &OtelHandler{
		next: baseHandler,
	}

	logger := slog.New(otelHandler)
	slog.SetDefault(logger)
}
