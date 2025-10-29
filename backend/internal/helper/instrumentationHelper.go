package helper

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentationHelper provides helper methods for instrumenting external service calls
//
// This helps track:
// - Database operations (queries, connections)
// - MinIO operations (upload, download, delete)
// - KMS operations (key operations)
// - External HTTP calls (Hydra OAuth2)
//
// Usage:
//
//	helper := NewInstrumentationHelper("crypsis-backend")
//	ctx, span := helper.StartDatabaseSpan(ctx, "SELECT", "files")
//	defer span.End()
type InstrumentationHelper struct {
	tracer trace.Tracer
	meter  metric.Meter

	// Metrics for external services
	dbQueryDuration     metric.Float64Histogram
	dbQueryCounter      metric.Int64Counter
	dbConnectionCounter metric.Int64UpDownCounter

	minioOpDuration metric.Float64Histogram
	minioOpCounter  metric.Int64Counter
	minioOpSize     metric.Int64Histogram

	kmsOpDuration metric.Float64Histogram
	kmsOpCounter  metric.Int64Counter

	httpClientDuration metric.Float64Histogram
	httpClientCounter  metric.Int64Counter
}

// NewInstrumentationHelper creates a new instrumentation helper
//
// Usage:
//
//	helper := NewInstrumentationHelper("crypsis-backend")
func NewInstrumentationHelper(serviceName string) *InstrumentationHelper {
	tracer := otel.Tracer(serviceName)
	meter := otel.Meter(serviceName)

	// Initialize database metrics
	dbQueryDuration, _ := meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Duration of database queries"),
		metric.WithUnit("ms"),
	)

	dbQueryCounter, _ := meter.Int64Counter(
		"db.query.count",
		metric.WithDescription("Number of database queries"),
		metric.WithUnit("{query}"),
	)

	dbConnectionCounter, _ := meter.Int64UpDownCounter(
		"db.connections.active",
		metric.WithDescription("Number of active database connections"),
		metric.WithUnit("{connection}"),
	)

	// Initialize MinIO metrics
	minioOpDuration, _ := meter.Float64Histogram(
		"minio.operation.duration",
		metric.WithDescription("Duration of MinIO operations"),
		metric.WithUnit("ms"),
	)

	minioOpCounter, _ := meter.Int64Counter(
		"minio.operation.count",
		metric.WithDescription("Number of MinIO operations"),
		metric.WithUnit("{operation}"),
	)

	minioOpSize, _ := meter.Int64Histogram(
		"minio.operation.size",
		metric.WithDescription("Size of data transferred in MinIO operations"),
		metric.WithUnit("By"),
	)

	// Initialize KMS metrics
	kmsOpDuration, _ := meter.Float64Histogram(
		"kms.operation.duration",
		metric.WithDescription("Duration of KMS operations"),
		metric.WithUnit("ms"),
	)

	kmsOpCounter, _ := meter.Int64Counter(
		"kms.operation.count",
		metric.WithDescription("Number of KMS operations"),
		metric.WithUnit("{operation}"),
	)

	// Initialize HTTP client metrics
	httpClientDuration, _ := meter.Float64Histogram(
		"http.client.duration",
		metric.WithDescription("Duration of HTTP client requests"),
		metric.WithUnit("ms"),
	)

	httpClientCounter, _ := meter.Int64Counter(
		"http.client.request.count",
		metric.WithDescription("Number of HTTP client requests"),
		metric.WithUnit("{request}"),
	)

	return &InstrumentationHelper{
		tracer:              tracer,
		meter:               meter,
		dbQueryDuration:     dbQueryDuration,
		dbQueryCounter:      dbQueryCounter,
		dbConnectionCounter: dbConnectionCounter,
		minioOpDuration:     minioOpDuration,
		minioOpCounter:      minioOpCounter,
		minioOpSize:         minioOpSize,
		kmsOpDuration:       kmsOpDuration,
		kmsOpCounter:        kmsOpCounter,
		httpClientDuration:  httpClientDuration,
		httpClientCounter:   httpClientCounter,
	}
}

// ================== DATABASE INSTRUMENTATION ==================

// StartDatabaseSpan creates a span for a database operation
//
// Usage:
//
//	ctx, span := helper.StartDatabaseSpan(ctx, "SELECT", "files", "id = ?")
//	defer span.End()
//	// Execute database query
//	if err != nil {
//	    helper.RecordDatabaseError(span, err)
//	}
//
// This creates a span with:
// - Operation type (SELECT, INSERT, UPDATE, DELETE)
// - Table name
// - Optional query statement
func (h *InstrumentationHelper) StartDatabaseSpan(
	ctx context.Context,
	operation string,
	table string,
	statement ...string,
) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", operation),
		attribute.String("db.table", table),
	}

	if len(statement) > 0 && statement[0] != "" {
		attrs = append(attrs, attribute.String("db.statement", statement[0]))
	}

	ctx, span := h.tracer.Start(ctx, "db."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)

	return ctx, span
}

// RecordDatabaseMetrics records metrics for a database operation
//
// Usage:
//
//	start := time.Now()
//	// Execute query
//	helper.RecordDatabaseMetrics(ctx, "SELECT", "files", time.Since(start), err)
func (h *InstrumentationHelper) RecordDatabaseMetrics(
	ctx context.Context,
	operation string,
	table string,
	duration time.Duration,
	err error,
) {
	durationMs := float64(duration.Milliseconds())

	attrs := metric.WithAttributes(
		attribute.String("db.operation", operation),
		attribute.String("db.table", table),
		attribute.Bool("error", err != nil),
	)

	h.dbQueryDuration.Record(ctx, durationMs, attrs)
	h.dbQueryCounter.Add(ctx, 1, attrs)
}

// RecordDatabaseError records an error in the current database span
func (h *InstrumentationHelper) RecordDatabaseError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ================== MinIO INSTRUMENTATION ==================

// StartMinIOSpan creates a span for a MinIO operation
//
// Usage:
//
//	ctx, span := helper.StartMinIOSpan(ctx, "PutObject", "my-bucket", "file.txt")
//	defer span.End()
//	// Execute MinIO operation
//	if err != nil {
//	    helper.RecordMinIOError(span, err)
//	}
//
// Operations: PutObject, GetObject, DeleteObject, StatObject, ListObjects
func (h *InstrumentationHelper) StartMinIOSpan(
	ctx context.Context,
	operation string,
	bucket string,
	objectKey string,
) (context.Context, trace.Span) {
	ctx, span := h.tracer.Start(ctx, "minio."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("minio.operation", operation),
			attribute.String("minio.bucket", bucket),
			attribute.String("minio.object_key", objectKey),
		),
	)

	return ctx, span
}

// RecordMinIOMetrics records metrics for a MinIO operation
//
// Usage:
//
//	start := time.Now()
//	// Execute MinIO operation
//	helper.RecordMinIOMetrics(ctx, "PutObject", "my-bucket", time.Since(start), fileSize, err)
func (h *InstrumentationHelper) RecordMinIOMetrics(
	ctx context.Context,
	operation string,
	bucket string,
	duration time.Duration,
	size int64,
	err error,
) {
	durationMs := float64(duration.Milliseconds())

	attrs := metric.WithAttributes(
		attribute.String("minio.operation", operation),
		attribute.String("minio.bucket", bucket),
		attribute.Bool("error", err != nil),
	)

	h.minioOpDuration.Record(ctx, durationMs, attrs)
	h.minioOpCounter.Add(ctx, 1, attrs)

	if size > 0 {
		h.minioOpSize.Record(ctx, size, attrs)
	}
}

// RecordMinIOError records an error in the current MinIO span
func (h *InstrumentationHelper) RecordMinIOError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ================== KMS INSTRUMENTATION ==================

// StartKMSSpan creates a span for a KMS operation
//
// Usage:
//
//	ctx, span := helper.StartKMSSpan(ctx, "ExportKey", "key-uid-123")
//	defer span.End()
//	// Execute KMS operation
//	if err != nil {
//	    helper.RecordKMSError(span, err)
//	}
//
// Operations: ExportKey, CreateKey, EncryptData, DecryptData
func (h *InstrumentationHelper) StartKMSSpan(
	ctx context.Context,
	operation string,
	keyID string,
) (context.Context, trace.Span) {
	ctx, span := h.tracer.Start(ctx, "kms."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("kms.operation", operation),
			attribute.String("kms.key_id", keyID),
		),
	)

	return ctx, span
}

// RecordKMSMetrics records metrics for a KMS operation
//
// Usage:
//
//	start := time.Now()
//	// Execute KMS operation
//	helper.RecordKMSMetrics(ctx, "ExportKey", time.Since(start), err)
func (h *InstrumentationHelper) RecordKMSMetrics(
	ctx context.Context,
	operation string,
	duration time.Duration,
	err error,
) {
	durationMs := float64(duration.Milliseconds())

	attrs := metric.WithAttributes(
		attribute.String("kms.operation", operation),
		attribute.Bool("error", err != nil),
	)

	h.kmsOpDuration.Record(ctx, durationMs, attrs)
	h.kmsOpCounter.Add(ctx, 1, attrs)
}

// RecordKMSError records an error in the current KMS span
func (h *InstrumentationHelper) RecordKMSError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ================== HTTP CLIENT INSTRUMENTATION ==================

// StartHTTPClientSpan creates a span for an HTTP client request
//
// Usage:
//
//	ctx, span := helper.StartHTTPClientSpan(ctx, "POST", "https://hydra:4445/oauth2/token")
//	defer span.End()
//	// Execute HTTP request
//	helper.SetHTTPClientSpanStatus(span, response.StatusCode)
//	if err != nil {
//	    helper.RecordHTTPClientError(span, err)
//	}
//
// This is useful for tracking calls to:
// - Hydra OAuth2 server
// - External APIs
// - Microservices
func (h *InstrumentationHelper) StartHTTPClientSpan(
	ctx context.Context,
	method string,
	url string,
) (context.Context, trace.Span) {
	ctx, span := h.tracer.Start(ctx, method+" "+url,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("http.method", method),
			attribute.String("http.url", url),
			attribute.String("http.client", "true"),
		),
	)

	return ctx, span
}

// SetHTTPClientSpanStatus sets the status code on an HTTP client span
func (h *InstrumentationHelper) SetHTTPClientSpanStatus(span trace.Span, statusCode int) {
	span.SetAttributes(attribute.Int("http.status_code", statusCode))

	if statusCode >= 400 {
		if statusCode >= 500 {
			span.SetStatus(codes.Error, "Server Error")
		} else {
			span.SetStatus(codes.Error, "Client Error")
		}
	} else {
		span.SetStatus(codes.Ok, "Success")
	}
}

// RecordHTTPClientMetrics records metrics for an HTTP client request
//
// Usage:
//
//	start := time.Now()
//	// Execute HTTP request
//	helper.RecordHTTPClientMetrics(ctx, "POST", "https://hydra:4445/oauth2/token", time.Since(start), statusCode)
func (h *InstrumentationHelper) RecordHTTPClientMetrics(
	ctx context.Context,
	method string,
	url string,
	duration time.Duration,
	statusCode int,
) {
	durationMs := float64(duration.Milliseconds())

	attrs := metric.WithAttributes(
		attribute.String("http.method", method),
		attribute.String("http.url", url),
		attribute.Int("http.status_code", statusCode),
		attribute.Bool("error", statusCode >= 400),
	)

	h.httpClientDuration.Record(ctx, durationMs, attrs)
	h.httpClientCounter.Add(ctx, 1, attrs)
}

// RecordHTTPClientError records an error in the current HTTP client span
func (h *InstrumentationHelper) RecordHTTPClientError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ================== CRYPTO INSTRUMENTATION ==================

// RecordCryptoOperation records metrics for cryptographic operations
//
// This is useful for tracking:
// - Encryption/decryption performance
// - File sizes being processed
// - Algorithm performance
//
// Usage:
//
//	start := time.Now()
//	// Perform encryption
//	helper.RecordCryptoOperation(ctx, "encrypt", "AES-GCM", fileSize, time.Since(start), err)
func (h *InstrumentationHelper) RecordCryptoOperation(
	ctx context.Context,
	operation string, // "encrypt" or "decrypt"
	algorithm string, // "AES-GCM", "ChaCha20-Poly1305", etc.
	dataSize int64, // Size of data processed
	duration time.Duration,
	err error,
) {
	meter := otel.Meter("crypsis-crypto")

	// Duration metric
	cryptoDuration, _ := meter.Float64Histogram(
		"crypto.operation.duration",
		metric.WithDescription("Duration of cryptographic operations"),
		metric.WithUnit("ms"),
	)

	// Counter metric
	cryptoCounter, _ := meter.Int64Counter(
		"crypto.operation.count",
		metric.WithDescription("Number of cryptographic operations"),
		metric.WithUnit("{operation}"),
	)

	// Data size metric
	cryptoSize, _ := meter.Int64Histogram(
		"crypto.operation.data_size",
		metric.WithDescription("Size of data in cryptographic operations"),
		metric.WithUnit("By"),
	)

	durationMs := float64(duration.Milliseconds())
	attrs := metric.WithAttributes(
		attribute.String("crypto.operation", operation),
		attribute.String("crypto.algorithm", algorithm),
		attribute.Bool("error", err != nil),
	)

	cryptoDuration.Record(ctx, durationMs, attrs)
	cryptoCounter.Add(ctx, 1, attrs)
	cryptoSize.Record(ctx, dataSize, attrs)
}
