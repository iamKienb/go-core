package observabilityx

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func InitMeterProvider() (*metric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("observability core: failed to create prometheus exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}

func NewMetricsHandler() http.Handler {
	return promhttp.Handler()
}
