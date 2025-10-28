package middlewere

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// OpenTelemetryMiddleware creates a middleware that traces HTTP requests
func OpenTelemetryMiddleware(tracer trace.Tracer, meter metric.Meter) gin.HandlerFunc {
	// Create metrics
	requestCounter, _ := meter.Int64Counter(
		"http.server.request_count",
		metric.WithDescription("Total number of HTTP requests"),
	)

	requestDuration, _ := meter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("ms"),
	)

	return func(c *gin.Context) {
		// Start a new span
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+" "+c.FullPath(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.route", c.FullPath()),
				attribute.String("http.client_ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Store context in gin.Context
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		// Add response attributes to span
		span.SetAttributes(
			attribute.Int("http.status_code", status),
			attribute.Int64("http.response_size", int64(c.Writer.Size())),
		)

		// Record error if any
		if len(c.Errors) > 0 {
			span.SetStatus(codes.Error, c.Errors.String())
			span.RecordError(c.Errors.Last().Err)
		} else if status >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
		}

		// Record metrics
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		attrs := metric.WithAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", path),
			attribute.Int("http.status_code", status),
		)

		requestCounter.Add(context.Background(), 1, attrs)
		requestDuration.Record(context.Background(), float64(duration), attrs)
	}
}

// GetTraceID returns the trace ID from the context
func GetTraceID(c *gin.Context) string {
	spanContext := trace.SpanFromContext(c.Request.Context()).SpanContext()
	if spanContext.IsValid() {
		return spanContext.TraceID().String()
	}
	return ""
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(c *gin.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(c.Request.Context())
	span.AddEvent(name, trace.WithAttributes(attrs...))
}
