package observability

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	clientLabel = "client"
	urlLabel    = "url"

	maxURLSegments = 4
)

var (
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_client_duration_seconds",
		Help:    "Duration of HTTP client calls",
		Buckets: prometheus.DefBuckets,
	}, []string{clientLabel, urlLabel})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_client_requests_total",
		Help: "Total number of HTTP client requests",
	}, []string{clientLabel, urlLabel})

	httpRequestErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_client_errors_total",
		Help: "Total number of HTTP client requests that resulted in an error",
	}, []string{clientLabel, urlLabel})
)

// MonitoredTransport wraps an http.RoundTripper to add metrics
type MonitoredTransport struct {
	rt         http.RoundTripper
	clientName string
}

// NewMonitoredTransport creates a new instance of MonitoredTransport
func NewMonitoredTransport(rt http.RoundTripper, clientName string) *MonitoredTransport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &MonitoredTransport{
		rt:         rt,
		clientName: clientName,
	}
}

// RoundTrip implements the RoundTripper interface
func (t *MonitoredTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	labelURL := normalizeURL(*req.URL, maxURLSegments)

	// Increment the total requests counter
	httpRequestsTotal.WithLabelValues(t.clientName, labelURL).Inc()

	// Start timer
	start := time.Now()
	resp, err := t.rt.RoundTrip(req)
	// Measure duration
	duration := time.Since(start).Seconds()

	// Set the duration metric
	httpRequestDuration.WithLabelValues(t.clientName, labelURL).Observe(duration)

	if err != nil {
		// Increment the error counter
		httpRequestErrorsTotal.WithLabelValues(t.clientName, labelURL).Inc()
		return nil, err
	}

	return resp, err
}

func normalizeURL(url url.URL, maxSegments int) string {
	parts := strings.Split(url.Path, "/")
	var segments []string
	for _, part := range parts {
		if part != "" {
			segments = append(segments, part)
		}
	}

	if len(segments) > maxSegments {
		segments = segments[:maxSegments]
	}

	// Reconstruct the normalized path
	url.Path = "/" + strings.Join(segments, "/")
	// Remove query parameters if any
	url.RawQuery = ""

	return url.String()
}
