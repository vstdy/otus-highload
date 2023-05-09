package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// IGaugeVec ...
type IGaugeVec interface {
	Set(value float64, labelValues ...string)
}

// GaugeVec ...
type GaugeVec struct {
	gauge *prometheus.GaugeVec
}

// MustRegisterGaugeVec ...
func MustRegisterGaugeVec(opts prometheus.GaugeOpts, labelNames ...string) *GaugeVec {
	return &GaugeVec{gauge: promauto.NewGaugeVec(opts, labelNames)}
}

// Set ...
func (c *GaugeVec) Set(value float64, labels prometheus.Labels) {
	c.gauge.With(labels).Set(value)
}
