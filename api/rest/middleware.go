package rest

import (
	"net/http"
	"time"

	"github.com/vstdy/otus-highload/api/rest/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	Status int
}

func (r *responseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func addMetrics() func(next http.Handler) http.Handler {
	httpMetrics := metrics.BuildHTTPMetrics()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &responseWriter{
				ResponseWriter: w,
				Status:         http.StatusOK,
			}
			path := r.URL.EscapedPath()
			start := time.Now()

			next.ServeHTTP(recorder, r)

			httpMetrics.SaveHTTPDurationSummary(start, path, recorder.Status, r.Method)
			httpMetrics.SaveHTTPDurationHistogram(start, path, recorder.Status, r.Method)
			httpMetrics.SaveHTTPCount(1, path, recorder.Status, r.Method)
		})
	}
}
