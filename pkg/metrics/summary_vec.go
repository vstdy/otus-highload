package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ISummaryVec ...
type ISummaryVec interface {
	Observe(float64)
}

// SummaryVec ...
type SummaryVec struct {
	summary *prometheus.SummaryVec
}

// MustRegisterSummaryVec ...
func MustRegisterSummaryVec(opts prometheus.SummaryOpts, labelNames ...string) *SummaryVec {
	return &SummaryVec{summary: promauto.NewSummaryVec(opts, labelNames)}
}

// Observe ...
func (c *SummaryVec) Observe(value float64, labels prometheus.Labels) {
	c.summary.With(labels).Observe(value)
}
