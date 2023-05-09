package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/vstdy/otus-highload/pkg/metrics"
)

const AppName = "otus-project"

const (
	labelApp    = "app"
	labelPath   = "path"
	labelCode   = "code"
	labelMethod = "method"
)

type HTTPMetrics struct {
	httpCounter           *metrics.CounterVec
	httpDurationSummary   *metrics.SummaryVec
	httpDurationHistogram *metrics.HistogramVec
}

func BuildHTTPMetrics() HTTPMetrics {
	return HTTPMetrics{
		httpCounter:           NewHTTPRequestCounter(),
		httpDurationSummary:   NewHTTPRequestDurationSummary(),
		httpDurationHistogram: NewHTTPRequestDurationHistogram(),
	}
}

func (hm *HTTPMetrics) SaveHTTPCount(value float64, path string, code int, method string) {
	hm.httpCounter.Add(
		value,
		prometheus.Labels{
			labelApp:    AppName,
			labelPath:   path,
			labelCode:   strconv.Itoa(code),
			labelMethod: method,
		})
}

func (hm *HTTPMetrics) SaveHTTPDurationSummary(timeSince time.Time, path string, code int, method string) {
	hm.httpDurationSummary.Observe(
		float64(time.Since(timeSince).Milliseconds()),
		prometheus.Labels{
			labelApp:    AppName,
			labelPath:   path,
			labelCode:   strconv.Itoa(code),
			labelMethod: method,
		})
}

func (hm *HTTPMetrics) SaveHTTPDurationHistogram(timeSince time.Time, path string, code int, method string) {
	hm.httpDurationHistogram.Observe(
		float64(time.Since(timeSince).Milliseconds()),
		prometheus.Labels{
			labelApp:    AppName,
			labelPath:   path,
			labelCode:   strconv.Itoa(code),
			labelMethod: method,
		})
}
