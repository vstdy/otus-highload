package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/vstdy/otus-highload/pkg/metrics"
)

const (
	httpRequestCounterName = "http_request_count"
	httpRequestCounterHelp = "Count of HTTP requests"
)

// NewHTTPRequestCounter ...
func NewHTTPRequestCounter() *metrics.CounterVec {
	return metrics.MustRegisterCounterVec(
		prometheus.CounterOpts{
			Name: httpRequestCounterName,
			Help: httpRequestCounterHelp,
		},
		labelApp,
		labelPath,
		labelCode,
		labelMethod,
	)
}
