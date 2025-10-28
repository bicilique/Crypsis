package helper

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type OpenTelemetryConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

// InitOpenTelemetry initializes OpenTelemetry with tracing and metrics
func InitOpenTelemetry(config OpenTelemetryConfig) (func(context.Context) error, error) {
	ctx := context.Background()

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Set up trace provider
	traceShutdown, err := initTracer(ctx, res, config.OTLPEndpoint)
	if err != nil {
		return nil, err
	}

	// Set up metric provider
	metricShutdown, err := initMetrics(ctx, res)
	if err != nil {
		traceShutdown(ctx)
		return nil, err
	}

	// Return a combined shutdown function
	shutdown := func(ctx context.Context) error {
		if err := metricShutdown(ctx); err != nil {
			return err
		}
		return traceShutdown(ctx)
	}

	return shutdown, nil
}

func initTracer(ctx context.Context, res *resource.Resource, endpoint string) (func(context.Context) error, error) {
	// Create OTLP trace exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), // Use WithTLSClientConfig for production
	)
	if err != nil {
		return nil, err
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(time.Second*5),
		),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func initMetrics(ctx context.Context, res *resource.Resource) (func(context.Context) error, error) {
	// Create Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	// Create meter provider
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exporter),
	)

	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}
