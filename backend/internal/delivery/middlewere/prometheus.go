package middlewere

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(httpRequests *prometheus.CounterVec, httpDuration *prometheus.HistogramVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path // fallback
		}

		status := fmt.Sprint(c.Writer.Status())
		method := c.Request.Method

		httpRequests.WithLabelValues(path, method, status).Inc()
		httpDuration.WithLabelValues(path, method).Observe(duration)
	}
}
