package helper

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GormTracingPlugin is a GORM plugin that automatically traces all database operations
// This provides automatic instrumentation without modifying repository code
type GormTracingPlugin struct {
	tracer trace.Tracer
}

// NewGormTracingPlugin creates a new GORM tracing plugin
func NewGormTracingPlugin(serviceName string) *GormTracingPlugin {
	return &GormTracingPlugin{
		tracer: otel.Tracer(serviceName),
	}
}

// Name returns the plugin name
func (p *GormTracingPlugin) Name() string {
	return "gorm:tracing"
}

// Initialize sets up the plugin callbacks
func (p *GormTracingPlugin) Initialize(db *gorm.DB) error {
	// Register callbacks for each GORM operation type

	// CREATE operations
	if err := db.Callback().Create().Before("gorm:create").Register("tracing:before_create", p.before("CREATE")); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:create").Register("tracing:after_create", p.after()); err != nil {
		return err
	}

	// QUERY operations
	if err := db.Callback().Query().Before("gorm:query").Register("tracing:before_query", p.before("SELECT")); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("tracing:after_query", p.after()); err != nil {
		return err
	}

	// UPDATE operations
	if err := db.Callback().Update().Before("gorm:update").Register("tracing:before_update", p.before("UPDATE")); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("tracing:after_update", p.after()); err != nil {
		return err
	}

	// DELETE operations
	if err := db.Callback().Delete().Before("gorm:delete").Register("tracing:before_delete", p.before("DELETE")); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("tracing:after_delete", p.after()); err != nil {
		return err
	}

	// RAW operations
	if err := db.Callback().Raw().Before("gorm:raw").Register("tracing:before_raw", p.before("RAW")); err != nil {
		return err
	}
	if err := db.Callback().Raw().After("gorm:raw").Register("tracing:after_raw", p.after()); err != nil {
		return err
	}

	// ROW operations
	if err := db.Callback().Row().Before("gorm:row").Register("tracing:before_row", p.before("ROW")); err != nil {
		return err
	}
	if err := db.Callback().Row().After("gorm:row").Register("tracing:after_row", p.after()); err != nil {
		return err
	}

	return nil
}

// before creates a callback that starts a span before the operation
func (p *GormTracingPlugin) before(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		// Get context from GORM
		ctx := db.Statement.Context
		if ctx == nil {
			return
		}

		// Extract table name if available
		tableName := db.Statement.Table
		if tableName == "" && db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		}

		// Create span name
		spanName := fmt.Sprintf("DB: %s", operation)
		if tableName != "" {
			spanName = fmt.Sprintf("DB: %s %s", operation, tableName)
		}

		// Start the span
		ctx, span := p.tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)

		// Add database attributes
		span.SetAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", operation),
		)

		if tableName != "" {
			span.SetAttributes(attribute.String("db.table", tableName))
		}

		// Add SQL query if available in development mode
		// This is commented out by default to avoid logging sensitive data in production
		// Uncomment the following block if you need to debug SQL queries:
		/*
			if db.Statement.SQL.String() != "" {
				span.SetAttributes(attribute.String("db.query", db.Statement.SQL.String()))
			}
		*/

		// Store span in context for the "after" callback
		db.Statement.Context = ctx
		db.InstanceSet("tracing:span", span)
	}
}

// after creates a callback that ends the span after the operation
func (p *GormTracingPlugin) after() func(*gorm.DB) {
	return func(db *gorm.DB) {
		// Retrieve the span from the instance
		spanInterface, ok := db.InstanceGet("tracing:span")
		if !ok {
			return
		}

		span, ok := spanInterface.(trace.Span)
		if !ok {
			return
		}
		defer span.End()

		// Add rows affected
		if db.Statement.RowsAffected > 0 {
			span.SetAttributes(attribute.Int64("db.rows_affected", db.Statement.RowsAffected))
		}

		// Record error if any
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			span.RecordError(db.Error)
			span.SetStatus(codes.Error, db.Error.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
}

// RegisterGormTracing registers the GORM tracing plugin with the database
// Call this once when initializing the database connection
//
// Usage:
//
//	db, err := gorm.Open(...)
//	helper.RegisterGormTracing(db, "crypsis-backend")
func RegisterGormTracing(db *gorm.DB, serviceName string) error {
	plugin := NewGormTracingPlugin(serviceName)
	return db.Use(plugin)
}
