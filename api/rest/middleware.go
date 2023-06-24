package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/vstdy/otus-highload/api/rest/metrics"
)

func addMetrics() func(next http.Handler) http.Handler {
	httpMetrics := metrics.BuildHTTPMetrics()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			path := r.URL.EscapedPath()
			start := time.Now()

			next.ServeHTTP(recorder, r)

			httpMetrics.SaveHTTPDurationSummary(start, path, recorder.Status(), r.Method)
			httpMetrics.SaveHTTPDurationHistogram(start, path, recorder.Status(), r.Method)
			httpMetrics.SaveHTTPCount(1, path, recorder.Status(), r.Method)
		})
	}
}
