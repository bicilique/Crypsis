package helper

import (
	"context"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// OpenTelemetryConfig contains configuration for initializing OpenTelemetry
// This struct defines all necessary parameters to set up distributed tracing and metrics collection
type OpenTelemetryConfig struct {
	ServiceName    string // Name of the service (e.g., "crypsis-backend")
	ServiceVersion string // Version of the service (e.g., "1.0.0")
	Environment    string // Deployment environment (e.g., "production", "staging", "development")
	OTLPEndpoint   string // OTLP collector endpoint (e.g., "otel-collector:4318")
}

// InitOpenTelemetry initializes the complete OpenTelemetry stack for the application
//
// This function sets up:
// 1. Distributed tracing - Track requests across services
// 2. Metrics collection - Collect performance and business metrics
// 3. Resource attributes - Add metadata to all telemetry
// 4. Runtime metrics - Monitor Go runtime (memory, goroutines, GC)
//
// Flow: Application -> OTLP Exporter -> OpenTelemetry Collector -> Jaeger & Prometheus -> Grafana
//
// Usage:
//
//	shutdown, err := InitOpenTelemetry(OpenTelemetryConfig{
//	    ServiceName:    "crypsis-backend",
//	    ServiceVersion: "1.0.0",
//	    Environment:    "production",
//	    OTLPEndpoint:   "otel-collector:4318",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer shutdown(context.Background())
func InitOpenTelemetry(config OpenTelemetryConfig) (func(context.Context) error, error) {
	ctx := context.Background()

	// Step 1: Create a resource that describes this service
	// Resources are attached to all telemetry and provide context about where it came from
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// Service identification
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),

			// Additional useful attributes
			attribute.String("service.instance.id", getHostname()),
			attribute.String("runtime.name", "go"),
			attribute.String("runtime.version", runtime.Version()),
		),
	)
	if err != nil {
		return nil, err
	}

	// Step 2: Initialize distributed tracing
	// Traces show the journey of a request through the system
	traceShutdown, err := initTracer(ctx, res, config.OTLPEndpoint)
	if err != nil {
		return nil, err
	}

	// Step 3: Initialize metrics collection
	// Metrics provide aggregated data about system performance
	metricShutdown, err := initMetrics(ctx, res, config.OTLPEndpoint)
	if err != nil {
		if shutdownErr := traceShutdown(ctx); shutdownErr != nil {
			// Log the shutdown error but return the original error
			// In production, use a proper logger
			println("Failed to shutdown tracer:", shutdownErr.Error())
		}
		return nil, err
	}

	// Step 4: Start runtime metrics collection
	// This monitors Go runtime statistics (memory, goroutines, GC)
	if err := startRuntimeMetrics(); err != nil {
		// Runtime metrics are optional, log error but continue
		// In production, use a proper logger
		println("Failed to start runtime metrics:", err.Error())
	}

	// Return a combined shutdown function
	// This ensures clean shutdown of all telemetry components
	shutdown := func(ctx context.Context) error {
		// Shutdown metrics first to stop collecting data
		if err := metricShutdown(ctx); err != nil {
			return err
		}
		// Then shutdown tracer to flush any remaining traces
		return traceShutdown(ctx)
	}

	return shutdown, nil
}

// initTracer sets up the OpenTelemetry trace provider
//
// Tracing workflow:
// 1. Application creates spans (units of work)
// 2. OTLP exporter sends spans to collector via HTTP
// 3. Collector processes and forwards to Jaeger
// 4. Jaeger stores and visualizes traces
//
// Configuration:
// - BatchSpanProcessor: Batches spans before export for efficiency
// - AlwaysSample: Captures 100% of traces (adjust for production)
func initTracer(ctx context.Context, res *resource.Resource, endpoint string) (func(context.Context) error, error) {
	// Create OTLP HTTP exporter for traces
	// This sends trace data to the OpenTelemetry Collector
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),                                 // Use TLS in production with WithTLSClientConfig
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression), // Compress data
	)
	if err != nil {
		return nil, err
	}

	// Create and configure the trace provider
	tp := trace.NewTracerProvider(
		// Batch span processor - batches spans before sending for efficiency
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(5*time.Second), // Send batch every 5 seconds
			trace.WithMaxExportBatchSize(512),     // Or when batch reaches 512 spans
		),

		// Resource - metadata about this service
		trace.WithResource(res),

		// Sampler - determines which traces to capture
		// AlwaysSample: Capture 100% of traces
		// For production, consider: trace.ParentBased(trace.TraceIDRatioBased(0.1)) for 10% sampling
		trace.WithSampler(trace.AlwaysSample()),
	)

	// Set as global tracer provider
	// This allows any code to get a tracer with: otel.Tracer("my-component")
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

// initMetrics sets up the OpenTelemetry meter provider
//
// Metrics workflow:
// 1. Application records metrics (counters, histograms, gauges)
// 2. OTLP exporter sends metrics to collector via HTTP
// 3. Collector processes and exports to Prometheus
// 4. Prometheus stores time-series data
// 5. Grafana queries Prometheus for visualization
//
// Metric types:
// - Counter: Monotonically increasing value (e.g., request count)
// - Histogram: Distribution of values (e.g., request duration)
// - Gauge: Point-in-time value (e.g., active connections)
func initMetrics(ctx context.Context, res *resource.Resource, endpoint string) (func(context.Context) error, error) {
	// Create OTLP HTTP exporter for metrics
	// This sends metric data to the OpenTelemetry Collector
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithInsecure(), // Use TLS in production
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	)
	if err != nil {
		return nil, err
	}

	// Create and configure the meter provider
	mp := sdkmetric.NewMeterProvider(
		// Resource - metadata about this service
		sdkmetric.WithResource(res),

		// Periodic reader - collects and exports metrics at regular intervals
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(10*time.Second), // Export metrics every 10 seconds
			),
		),
	)

	// Set as global meter provider
	// This allows any code to get a meter with: otel.Meter("my-component")
	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}

// startRuntimeMetrics registers and starts collection of Go runtime metrics
//
// This monitors:
// - Memory usage (heap, stack, allocations)
// - Goroutines (count, state)
// - GC (garbage collection) metrics
// - CPU usage
//
// These metrics are crucial for:
// - Detecting memory leaks
// - Identifying goroutine leaks
// - Understanding GC pressure
// - Monitoring resource usage
func startRuntimeMetrics() error {
	meter := otel.Meter("crypsis-runtime")

	// Memory metrics
	_, err := meter.Int64ObservableGauge(
		"runtime.go.mem.heap_alloc",
		metric.WithDescription("Bytes of allocated heap objects"),
		metric.WithUnit("bytes"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			observer.Observe(int64(m.HeapAlloc))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	_, err = meter.Int64ObservableGauge(
		"runtime.go.mem.heap_idle",
		metric.WithDescription("Bytes in idle (unused) heap spans"),
		metric.WithUnit("bytes"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			observer.Observe(int64(m.HeapIdle))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	_, err = meter.Int64ObservableGauge(
		"runtime.go.mem.heap_inuse",
		metric.WithDescription("Bytes in in-use heap spans"),
		metric.WithUnit("bytes"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			observer.Observe(int64(m.HeapInuse))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	_, err = meter.Int64ObservableGauge(
		"runtime.go.mem.heap_sys",
		metric.WithDescription("Bytes of heap memory obtained from the OS"),
		metric.WithUnit("bytes"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			observer.Observe(int64(m.HeapSys))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	// Goroutine metrics
	_, err = meter.Int64ObservableGauge(
		"runtime.go.goroutines",
		metric.WithDescription("Number of goroutines that currently exist"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(int64(runtime.NumGoroutine()))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	// GC metrics
	_, err = meter.Int64ObservableCounter(
		"runtime.go.gc.count",
		metric.WithDescription("Number of completed GC cycles"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			observer.Observe(int64(m.NumGC))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	_, err = meter.Float64ObservableGauge(
		"runtime.go.gc.pause_ns",
		metric.WithDescription("Pause time of the last GC cycle in nanoseconds"),
		metric.WithUnit("ns"),
		metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if m.NumGC > 0 {
				observer.Observe(float64(m.PauseNs[(m.NumGC+255)%256]))
			}
			return nil
		}),
	)
	if err != nil {
		return err
	}

	// CPU metrics
	_, err = meter.Int64ObservableGauge(
		"runtime.go.cpu.count",
		metric.WithDescription("Number of logical CPUs"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(int64(runtime.NumCPU()))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

// getHostname returns the hostname or a default value
// This is used to identify which instance of the service generated telemetry
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
