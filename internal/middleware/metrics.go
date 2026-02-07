package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpServerRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{observability.MethodLabel, observability.PathLabel, observability.CodeLabel},
	)
	httpServerRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_request_duration_seconds",
			Help:    "Histogram of response duration for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{observability.MethodLabel, observability.PathLabel},
	)
	httpServerInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_server_in_flight_requests",
		Help: "Number of HTTP requests currently being processed",
	})
)

// statusCodeResponseWriter wraps http.ResponseWriter to capture the HTTP status code
type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader writes the HTTP status code to the response
func (w *statusCodeResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Flush implements http.Flusher interface for streaming responses
func (w *statusCodeResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Unwrap returns the underlying ResponseWriter for middleware compatibility
func (w *statusCodeResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// Metrics records metrics for each HTTP request
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &statusCodeResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		httpServerInFlight.Inc()
		defer httpServerInFlight.Dec()

		start := time.Now()

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		// use chi's route pattern to avoid high cardinality
		path := chi.RouteContext(r.Context()).RoutePattern()
		if path == "" {
			path = r.URL.Path
		}

		code := strconv.Itoa(rw.statusCode)

		httpServerRequestsTotal.WithLabelValues(r.Method, path, code).Inc()
		httpServerRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}
