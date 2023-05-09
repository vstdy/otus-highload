package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewMetricsServer ...
func NewMetricsServer(config Config) (*http.Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	mux := http.DefaultServeMux
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{Addr: config.MetricsServerAddress, Handler: mux}, nil
}
