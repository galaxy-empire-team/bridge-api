package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	subsystem = "gp"
)

type PrometheusClient struct {
	requests *prometheus.CounterVec
	latency  *prometheus.HistogramVec
}

var (
	instance *PrometheusClient
	once     sync.Once
)

func NewPrometheusClient() *PrometheusClient {
	once.Do(func() {
		requests := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      "requests_total",
				Help:      "The total number of requests",
			},
			[]string{"endpoint", "method", "status"},
		)

		latency := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      "latency_seconds",
				Help:      "The latency of requests",
			},
			[]string{"endpoint", "method"},
		)

		prometheus.MustRegister(
			requests,
			latency,
		)

		instance = &PrometheusClient{
			requests: requests,
			latency:  latency,
		}
	})

	return instance
}

func (c *PrometheusClient) IncrementRequests(endpoint, method, status string) {
	c.requests.WithLabelValues(endpoint, method, status).Inc()
}

func (c *PrometheusClient) SetLatency(endpoint, method string, latency float64) {
	c.latency.WithLabelValues(endpoint, method).Observe(latency)
}
