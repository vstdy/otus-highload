package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/vstdy/otus-highload/pkg/metrics"
)

const (
	httpRequestDurationSummaryName = "http_request_duration_summary"
	httpRequestDurationSummaryHelp = "Duration of HTTP requests Summary"
)

var httpRequestDurationSummaryObjectives = map[float64]float64{0.5: 0.5, 0.9: 0.9, 0.99: 0.99}
var httpRequestDurationSummaryAgeBuckets uint32 = 3
var httpRequestDurationSummaryMaxAge = 120 * time.Second

// NewHTTPRequestDurationSummary ...
func NewHTTPRequestDurationSummary() *metrics.SummaryVec {
	return metrics.MustRegisterSummaryVec(
		prometheus.SummaryOpts{
			Name:       httpRequestDurationSummaryName,
			Help:       httpRequestDurationSummaryHelp,
			Objectives: httpRequestDurationSummaryObjectives,
			AgeBuckets: httpRequestDurationSummaryAgeBuckets,
			MaxAge:     httpRequestDurationSummaryMaxAge,
		},
		labelApp,
		labelPath,
		labelCode,
		labelMethod,
	)
}
