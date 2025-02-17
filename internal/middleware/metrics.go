package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	methodLabel     = "method"
	pathLabel       = "path"
	statusCodeLabel = "status_code"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{methodLabel, pathLabel},
	)
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{methodLabel, pathLabel, statusCodeLabel},
	)
	serverErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_errors_total",
			Help: "Total number of server errors (5xx)",
		},
		[]string{methodLabel, pathLabel, statusCodeLabel},
	)
	responseDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Histogram of response duration for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{methodLabel, pathLabel, statusCodeLabel},
	)
)

// wraps http.ResponseWriter to capture the HTTP status code
type statusCodeResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader writes the HTTP status code to the response
func (w *statusCodeResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Metrics records metrics for each HTTP request
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &statusCodeResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		start := time.Now()

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		requestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()

		if rw.StatusCode >= 400 {
			errorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.StatusCode)).Inc()
		}
		if rw.StatusCode >= 500 {
			serverErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.StatusCode)).Inc()
		}

		responseDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.StatusCode)).Observe(duration)
	})
}
