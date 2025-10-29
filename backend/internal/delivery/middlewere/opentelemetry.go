package middlewere

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// OpenTelemetryMiddleware creates a comprehensive middleware for HTTP request observability
//
// This middleware instruments all HTTP requests with:
// 1. Distributed Tracing - Creates spans for each request
// 2. Metrics Collection - Records request counts, durations, and sizes
// 3. Error Tracking - Captures and records errors
//
// Flow:
//
//	Request → Middleware → Create Span → Record Metrics → Handle Request → End Span
//
// Usage:
//
//	tracer := otel.Tracer("crypsis-backend")
//	meter := otel.Meter("crypsis-backend")
//	router.Use(OpenTelemetryMiddleware(tracer, meter))
//
// This follows OpenTelemetry best practices:
// - Semantic conventions for HTTP (https://opentelemetry.io/docs/specs/semconv/http/)
// - Proper span lifecycle management
// - Comprehensive attribute tagging
// - Error handling and recording
func OpenTelemetryMiddleware(tracer trace.Tracer, meter metric.Meter) gin.HandlerFunc {
	// Initialize metrics once when middleware is created
	// These metrics follow OpenTelemetry semantic conventions

	// Counter: Total number of HTTP requests
	// Labels: http.method, http.route, http.status_code
	requestCounter, _ := meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)

	// Histogram: Distribution of HTTP request durations
	// This is useful for calculating percentiles (p50, p90, p95, p99)
	// Labels: http.method, http.route, http.status_code
	requestDuration, _ := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("HTTP request duration in milliseconds"),
		metric.WithUnit("ms"),
	)

	// Histogram: Distribution of HTTP request body sizes
	// Labels: http.method, http.route
	requestSize, _ := meter.Int64Histogram(
		"http.server.request.body.size",
		metric.WithDescription("HTTP request body size in bytes"),
		metric.WithUnit("By"),
	)

	// Histogram: Distribution of HTTP response body sizes
	// Labels: http.method, http.route, http.status_code
	responseSize, _ := meter.Int64Histogram(
		"http.server.response.body.size",
		metric.WithDescription("HTTP response body size in bytes"),
		metric.WithUnit("By"),
	)

	// Counter: Total number of active HTTP requests (in-flight requests)
	// This is a gauge-like metric using UpDownCounter
	activeRequests, _ := meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)

	return func(c *gin.Context) {
		// Step 1: Extract or generate trace context
		// This allows distributed tracing across services
		ctx := c.Request.Context()

		// Step 2: Create a new span for this HTTP request
		// Span represents a unit of work in distributed tracing
		spanName := c.Request.Method + " " + c.FullPath()
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer), // This is a server-side span
			trace.WithAttributes(
				// HTTP semantic conventions
				// See: https://opentelemetry.io/docs/specs/semconv/http/http-spans/

				// Request attributes
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.scheme", c.Request.URL.Scheme),
				attribute.String("http.target", c.Request.URL.Path),
				attribute.String("http.route", c.FullPath()),
				attribute.String("http.url", c.Request.URL.String()),

				// Network attributes
				attribute.String("net.host.name", c.Request.Host),
				attribute.String("net.peer.ip", c.ClientIP()),

				// User agent
				attribute.String("http.user_agent", c.Request.UserAgent()),

				// Request size
				attribute.Int64("http.request.content_length", c.Request.ContentLength),
			),
		)
		defer span.End()

		// Step 3: Store the context with span in gin.Context
		// This makes the span available to handlers and other middleware
		c.Request = c.Request.WithContext(ctx)

		// Step 4: Record that a new request started (increase active requests)
		activeRequests.Add(ctx, 1)
		defer activeRequests.Add(ctx, -1) // Decrease when request ends

		// Step 5: Record start time for duration measurement
		startTime := time.Now()

		// Step 6: Process the HTTP request
		// This calls the next middleware/handler in the chain
		c.Next()

		// === Request Completed - Record Metrics and Span Attributes ===

		// Step 7: Calculate request duration
		duration := time.Since(startTime).Milliseconds()

		// Step 8: Get response details
		statusCode := c.Writer.Status()
		responseBodySize := int64(c.Writer.Size())

		// Step 9: Determine the actual route (for better cardinality)
		route := c.FullPath()
		if route == "" {
			// Fallback to request path if FullPath is empty (404 cases)
			route = c.Request.URL.Path
		}

		// Step 10: Create common attributes for metrics
		// These attributes allow filtering and grouping in Grafana
		commonAttrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", route),
			attribute.Int("http.status_code", statusCode),
		}
		metricAttrs := metric.WithAttributes(commonAttrs...)

		// Step 11: Update span with response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int64("http.response.content_length", responseBodySize),
		)

		// Step 12: Record error information if request failed
		if len(c.Errors) > 0 {
			// Request had errors - mark span as error
			lastError := c.Errors.Last()
			span.SetStatus(codes.Error, lastError.Error())
			span.RecordError(lastError.Err,
				trace.WithAttributes(
					attribute.String("error.type", fmt.Sprintf("%d", lastError.Type)),
					attribute.String("error.message", lastError.Error()),
				),
			)

			// Add error event to span timeline
			span.AddEvent("exception",
				trace.WithAttributes(
					attribute.String("exception.type", "GinError"),
					attribute.String("exception.message", lastError.Error()),
				),
			)
		} else if statusCode >= 400 {
			// HTTP error status code (4xx or 5xx)
			if statusCode >= 500 {
				span.SetStatus(codes.Error, "Internal Server Error")
			} else {
				span.SetStatus(codes.Error, "Client Error")
			}
		} else {
			// Success
			span.SetStatus(codes.Ok, "Success")
		}

		// Step 13: Record all metrics

		// Record request count
		requestCounter.Add(ctx, 1, metricAttrs)

		// Record request duration (in milliseconds)
		requestDuration.Record(ctx, float64(duration), metricAttrs)

		// Record request body size (if available)
		if c.Request.ContentLength > 0 {
			requestSize.Record(ctx, c.Request.ContentLength,
				metric.WithAttributes(
					attribute.String("http.method", c.Request.Method),
					attribute.String("http.route", route),
				),
			)
		}

		// Record response body size
		if responseBodySize > 0 {
			responseSize.Record(ctx, responseBodySize, metricAttrs)
		}
	}
}

// GetTraceID extracts the trace ID from the current request context
//
// Usage:
//
//	traceID := middlewere.GetTraceID(c)
//	log.Printf("Processing request with trace ID: %s", traceID)
//
// This is useful for:
// - Correlating logs with traces
// - Adding trace ID to response headers
// - Debugging specific requests
func GetTraceID(c *gin.Context) string {
	spanContext := trace.SpanFromContext(c.Request.Context()).SpanContext()
	if spanContext.IsValid() {
		return spanContext.TraceID().String()
	}
	return ""
}

// GetSpanID extracts the span ID from the current request context
//
// Usage:
//
//	spanID := middlewere.GetSpanID(c)
//
// Useful for detailed debugging and log correlation
func GetSpanID(c *gin.Context) string {
	spanContext := trace.SpanFromContext(c.Request.Context()).SpanContext()
	if spanContext.IsValid() {
		return spanContext.SpanID().String()
	}
	return ""
}

// AddSpanEvent adds a custom event to the current span
//
// Events represent significant points in the span's timeline
// Examples: "cache_hit", "validation_failed", "file_encrypted"
//
// Usage:
//
//	AddSpanEvent(c, "file_encrypted",
//	    attribute.String("file.id", fileID),
//	    attribute.Int64("file.size", fileSize),
//	)
//
// These events appear in Jaeger's trace timeline, helping you understand
// what happened during request processing
func AddSpanEvent(c *gin.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(c.Request.Context())
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// AddSpanAttributes adds custom attributes to the current span
//
// # Use this to add business-specific metadata to traces
//
// Usage:
//
//	AddSpanAttributes(c,
//	    attribute.String("user.id", userID),
//	    attribute.String("tenant.id", tenantID),
//	    attribute.Bool("premium.user", isPremium),
//	)
//
// These attributes can be used to filter traces in Jaeger
func AddSpanAttributes(c *gin.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attrs...)
}

// RecordError records an error in the current span
//
// Usage:
//
//	if err != nil {
//	    RecordError(c, err, "Failed to process file")
//	    c.JSON(500, gin.H{"error": err.Error()})
//	    return
//	}
//
// This marks the span as failed and adds error details
func RecordError(c *gin.Context, err error, message string) {
	span := trace.SpanFromContext(c.Request.Context())
	span.RecordError(err, trace.WithAttributes(
		attribute.String("error.message", message),
	))
	span.SetStatus(codes.Error, message)
}
