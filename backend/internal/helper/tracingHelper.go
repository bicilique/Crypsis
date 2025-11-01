package helper

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingHelper provides utilities for creating and managing OpenTelemetry spans
// This centralizes span creation logic and ensures consistent tracing across the application
type TracingHelper struct {
	tracer trace.Tracer
}

// NewTracingHelper creates a new tracing helper instance
func NewTracingHelper(serviceName string) *TracingHelper {
	return &TracingHelper{
		tracer: otel.Tracer(serviceName),
	}
}

// StartSpan creates a new span with the given name and options
// It returns the new context with the span and the span itself
//
// Usage:
//
//	ctx, span := helper.StartSpan(ctx, "operation.name")
//	defer span.End()
func (h *TracingHelper) StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return h.tracer.Start(ctx, spanName, opts...)
}

// StartDBSpan creates a span for database operations
// Automatically adds db.system attribute and formats the span name
//
// Usage:
//
//	ctx, span := helper.StartDBSpan(ctx, "SELECT", "files")
//	defer span.End()
func (h *TracingHelper) StartDBSpan(ctx context.Context, operation, table string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("DB: %s %s", operation, table)
	ctx, span := h.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
	)

	// Add database-specific attributes
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", operation),
		attribute.String("db.table", table),
	)

	return ctx, span
}

// StartKMSSpan creates a span for KMS operations
// Automatically adds kms-specific attributes
//
// Usage:
//
//	ctx, span := helper.StartKMSSpan(ctx, "EncryptKey", keyID)
//	defer span.End()
func (h *TracingHelper) StartKMSSpan(ctx context.Context, operation, keyID string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("KMS: %s", operation)
	ctx, span := h.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
	)

	// Add KMS-specific attributes
	span.SetAttributes(
		attribute.String("kms.operation", operation),
		attribute.String("kms.key_id", keyID),
	)

	return ctx, span
}

// StartStorageSpan creates a span for storage operations (MinIO/S3)
// Automatically adds storage-specific attributes
//
// Usage:
//
//	ctx, span := helper.StartStorageSpan(ctx, "PutObject", bucketName, objectName)
//	defer span.End()
func (h *TracingHelper) StartStorageSpan(ctx context.Context, operation, bucket, object string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("Storage: %s", operation)
	ctx, span := h.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
	)

	// Add storage-specific attributes
	span.SetAttributes(
		attribute.String("storage.operation", operation),
		attribute.String("storage.bucket", bucket),
		attribute.String("storage.object", object),
	)

	return ctx, span
}

// StartHTTPClientSpan creates a span for HTTP client calls
// Automatically adds http-specific attributes
//
// Usage:
//
//	ctx, span := helper.StartHTTPClientSpan(ctx, "GET", "https://api.example.com/users")
//	defer span.End()
func (h *TracingHelper) StartHTTPClientSpan(ctx context.Context, method, url string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("HTTP: %s %s", method, url)
	ctx, span := h.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
	)

	// Add HTTP-specific attributes
	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("http.url", url),
	)

	return ctx, span
}

// StartCryptoSpan creates a span for cryptographic operations
// Automatically adds crypto-specific attributes
//
// Usage:
//
//	ctx, span := helper.StartCryptoSpan(ctx, "Encrypt", "AES-256-GCM")
//	defer span.End()
func (h *TracingHelper) StartCryptoSpan(ctx context.Context, operation, algorithm string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("Crypto: %s", operation)
	ctx, span := h.tracer.Start(ctx, spanName)

	// Add crypto-specific attributes
	span.SetAttributes(
		attribute.String("crypto.operation", operation),
		attribute.String("crypto.algorithm", algorithm),
	)

	return ctx, span
}

// RecordError records an error in the span and sets the span status to error
// This is a convenience method to avoid boilerplate code
//
// Usage:
//
//	if err != nil {
//	    helper.RecordError(span, err)
//	    return err
//	}
func RecordError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// RecordSuccess sets the span status to OK with an optional message
//
// Usage:
//
//	helper.RecordSuccess(span, "File uploaded successfully")
func RecordSuccess(span trace.Span, message string) {
	span.SetStatus(codes.Ok, message)
}

// AddAttributes adds multiple attributes to a span at once
//
// Usage:
//
//	helper.AddAttributes(span, map[string]interface{}{
//	    "file.id": fileID,
//	    "file.size": fileSize,
//	    "file.type": mimeType,
//	})
func AddAttributes(span trace.Span, attrs map[string]interface{}) {
	for key, value := range attrs {
		switch v := value.(type) {
		case string:
			span.SetAttributes(attribute.String(key, v))
		case int:
			span.SetAttributes(attribute.Int(key, v))
		case int64:
			span.SetAttributes(attribute.Int64(key, v))
		case float64:
			span.SetAttributes(attribute.Float64(key, v))
		case bool:
			span.SetAttributes(attribute.Bool(key, v))
		default:
			span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
}

// Global tracing helper instance
var globalTracingHelper *TracingHelper

// InitGlobalTracingHelper initializes the global tracing helper
// Call this once during application startup
func InitGlobalTracingHelper(serviceName string) {
	globalTracingHelper = NewTracingHelper(serviceName)
}

// GetTracingHelper returns the global tracing helper instance
func GetTracingHelper() *TracingHelper {
	if globalTracingHelper == nil {
		// Fallback: create a default helper if not initialized
		globalTracingHelper = NewTracingHelper("crypsis-backend")
	}
	return globalTracingHelper
}
