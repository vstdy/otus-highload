package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/vstdy/otus-highload/pkg/metrics"
)

const (
	httpRequestDurationHistogramName = "http_request_duration_histogram"
	httpRequestDurationHistogramHelp = "Duration of HTTP requests Histogram"
)

var httpRequestDurationHistogramBuckets = []float64{10, 50, 100, 200, 500, 1000}

// NewHTTPRequestDurationHistogram ...
func NewHTTPRequestDurationHistogram() *metrics.HistogramVec {
	return metrics.MustRegisterHistogramVec(
		prometheus.HistogramOpts{
			Name:    httpRequestDurationHistogramName,
			Help:    httpRequestDurationHistogramHelp,
			Buckets: httpRequestDurationHistogramBuckets,
		},
		labelApp,
		labelPath,
		labelCode,
		labelMethod,
	)
}
