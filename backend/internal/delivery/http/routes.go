package http

import (
	"crypsis-backend/internal/delivery/middlewere"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// RouterConfig holds the configuration for HTTP routing
//
// This struct contains all dependencies needed to set up HTTP routes:
// - Gin router instance
// - HTTP handlers for client and admin operations
// - OAuth2 configuration (Hydra)
// - OpenTelemetry tracer and meter for observability
//
// The tracer and meter are used to create spans and record metrics
// for all HTTP requests, following OpenTelemetry best practices
type RouterConfig struct {
	Router          *gin.Engine
	ClientHandler   *ClientHandler
	AdminHandler    *AdminHandler
	HydraAdminURL   string
	TokenMiddlewere middlewere.TokenMiddlewareConfig
	Tracer          trace.Tracer // OpenTelemetry tracer for distributed tracing
	Meter           metric.Meter // OpenTelemetry meter for metrics collection
}

// Setup configures all HTTP routes and middleware
//
// This method:
// 1. Initializes OpenTelemetry instrumentation (if not provided)
// 2. Applies global middleware for tracing and metrics
// 3. Sets up route groups (public, client, admin, debug)
//
// Observability Flow:
//
//	All Requests → OpenTelemetry Middleware → OTLP → Collector → Jaeger & Prometheus → Grafana
//
// No direct Prometheus instrumentation is used - everything goes through OpenTelemetry
func (c *RouterConfig) Setup() {
	// Initialize OpenTelemetry tracer and meter if not provided
	// This uses the global TraceProvider and MeterProvider set up in main.go
	if c.Tracer == nil {
		c.Tracer = otel.Tracer("crypsis-backend")
	}
	if c.Meter == nil {
		c.Meter = otel.Meter("crypsis-backend")
	}

	// Apply global OpenTelemetry middleware to ALL routes
	// This automatically instruments every HTTP request with:
	// - Distributed tracing (spans)
	// - Metrics (request count, duration, size)
	// - Error tracking
	c.Router.Use(middlewere.OpenTelemetryMiddleware(c.Tracer, c.Meter))

	// Set up route groups
	c.setupPublic()
	c.setupClient()
	c.setupAdmin()
	c.setupDebug()
}

func (c *RouterConfig) setupPublic() {
	group := c.Router.Group("/api")
	group.POST("/admin/login", c.AdminHandler.Login)
}

func (c *RouterConfig) setupClient() {
	group := c.Router.Group("/api")
	group.Use(middlewere.TokenMiddleware(c.TokenMiddlewere))

	// group.Use(middlewere.PrometheusMiddleware(httpRequests, httpDuration)) // Apply Prometheus middleware

	group.POST("/files", c.ClientHandler.UploadFile)
	group.GET("/files/:id/download", c.ClientHandler.DownloadFile)
	group.PUT("/files/:id/update", c.ClientHandler.UpdateFile)
	group.DELETE("/files/:id/delete", c.ClientHandler.DeleteFile)

	group.GET("/files/list", c.ClientHandler.ListFiles)
	group.GET("/files/:id/metadata", c.ClientHandler.MetaDataFile)

	group.POST("/files/encrypt", c.ClientHandler.EncryptFile)
	group.POST("/files/decrypt", c.ClientHandler.DecryptFile)

	// Temporarily disabled due to potential security issues
	// group.POST("/files/:id/recover", c.ClientHandler.RecoverFile)
}

func (c *RouterConfig) setupAdmin() {
	group := c.Router.Group("/api")
	group.Use(middlewere.AdminTokenMiddleware(c.TokenMiddlewere))

	// Admin Account Management
	group.GET("/admin/logout", c.AdminHandler.Logout)
	group.GET("/admin/refresh-token", c.AdminHandler.RefreshToken)
	group.GET("/admin/list", c.AdminHandler.ListAdmin)
	group.PATCH("/admin/username", c.AdminHandler.UpdateAdminUsername)
	group.PATCH("/admin/password", c.AdminHandler.UpdateAdminPassword)
	group.DELETE("/admin", c.AdminHandler.DeleteAdmin)
	group.POST("/admin/add", c.AdminHandler.AddAdmin)

	// Application / Client Management
	group.POST("/admin/apps", c.AdminHandler.AddApp)
	group.GET("/admin/apps", c.AdminHandler.ListApps)
	group.GET("/admin/apps/:id", c.AdminHandler.GetApp)
	group.DELETE("/admin/apps/:id", c.AdminHandler.DeleteApp)
	group.POST("/admin/apps/:id/recover", c.AdminHandler.RecoverApp)
	group.PUT("/admin/apps/:id/rotate-secret", c.AdminHandler.RotateSecret)

	// File Management
	group.GET("/admin/files", c.AdminHandler.ListFiles)
	group.GET("/admin/apps/:id/files", c.AdminHandler.ListFilesByAppId)
	group.GET("/admin/logs", c.AdminHandler.ListLogs)
	group.POST("/admin/files/re-key", c.AdminHandler.Rekey)
}

// setupDebug sets up pprof debugging endpoints
func (c *RouterConfig) setupDebug() {
	debug := c.Router.Group("/debug/pprof")
	{
		debug.GET("/", gin.WrapF(pprof.Index))
		debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.GET("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
		debug.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		debug.GET("/block", gin.WrapH(pprof.Handler("block")))
		debug.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		debug.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		debug.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		debug.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
