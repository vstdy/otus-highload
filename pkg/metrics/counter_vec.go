package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ICounterVec ...
type ICounterVec interface {
	Add(value float64, labelValues ...string)
	Inc(labelValues ...string)
}

// CounterVec ...
type CounterVec struct {
	counter *prometheus.CounterVec
}

// MustRegisterCounterVec ...
func MustRegisterCounterVec(opts prometheus.CounterOpts, labelNames ...string) *CounterVec {
	return &CounterVec{counter: promauto.NewCounterVec(opts, labelNames)}
}

// Add ...
func (c *CounterVec) Add(value float64, labels prometheus.Labels) {
	c.counter.With(labels).Add(value)
}

// Inc ...
func (c *CounterVec) Inc(labels prometheus.Labels) {
	c.counter.With(labels).Inc()
}
