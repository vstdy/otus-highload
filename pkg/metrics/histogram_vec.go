package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// IHistogramVec ...
type IHistogramVec interface {
	Observe(value float64, labels prometheus.Labels)
}

// HistogramVec ...
type HistogramVec struct {
	histogram *prometheus.HistogramVec
}

// MustRegisterHistogramVec ...
func MustRegisterHistogramVec(opts prometheus.HistogramOpts, labelNames ...string) *HistogramVec {
	return &HistogramVec{histogram: promauto.NewHistogramVec(opts, labelNames)}
}

// Observe ...
func (c *HistogramVec) Observe(value float64, labels prometheus.Labels) {
	c.histogram.With(labels).Observe(value)
}
