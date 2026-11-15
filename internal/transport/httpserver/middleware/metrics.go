package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Metrics interface {
	IncrementRequests(endpoint, method, status string)
	SetLatency(endpoint, method string, latency float64)
}

func MetricsMiddleware(metrics Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		endpoint := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		metrics.SetLatency(endpoint, method, latency)
		metrics.IncrementRequests(endpoint, method, status)
	}
}
