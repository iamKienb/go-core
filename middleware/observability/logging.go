package observabilityx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	configx "github.com/iamKienb/go-core/config"
)

type LokiWriter struct {
	url         string
	serviceName string
	httpClient  *http.Client
}

func NewLokiWriter(serviceName string, lokiURL string) *LokiWriter {
	return &LokiWriter{
		url:         lokiURL,
		serviceName: serviceName,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (w *LokiWriter) Write(p []byte) (n int, err error) {
	if w.url == "" {
		return len(p), nil
	}

	logLine := string(p)

	lokiPayload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"job": w.serviceName,
				},
				"values": [][]string{
					{fmt.Sprintf("%d", time.Now().UnixNano()), logLine},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(lokiPayload)
	if err != nil {
		return 0, err
	}

	go func() {
		resp, err := w.httpClient.Post(w.url, "application/json", bytes.NewBuffer(payloadBytes))
		if err == nil {
			_ = resp.Body.Close()
		}
	}()

	return len(p), nil
}

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

func InitLogger(serverName string, cfg configx.LokiConfig) {
	terminalHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	lokiWriter := NewLokiWriter(serverName, cfg.LokiURL)
	lokiHandler := slog.NewJSONHandler(lokiWriter, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	multiHandler := NewMultiHandler(terminalHandler, lokiHandler)
	otelHandler := &OtelHandler{
		next: multiHandler,
	}

	logger := slog.New(otelHandler)
	slog.SetDefault(logger)
}

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.handlers[0].Enabled(ctx, level)
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, record)
	}
	return nil
}
