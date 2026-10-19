package observability_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

func TestInitTracer_ExportsAndPropagatesWorkerSpans(t *testing.T) {
	oldProvider := otel.GetTracerProvider()
	oldPropagator := otel.GetTextMapPropagator()
	oldErrorHandler := otel.GetErrorHandler()
	t.Cleanup(func() {
		otel.SetTracerProvider(oldProvider)
		otel.SetTextMapPropagator(oldPropagator)
		otel.SetErrorHandler(oldErrorHandler)
	})

	requests := make(chan []byte, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "/v1/traces", r.URL.Path)
		requests <- body
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	tp, err := observability.InitTracer(zap.NewNop(), strings.TrimPrefix(collector.URL, "http://"), "worker-test")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, tp.Shutdown(context.Background())) })

	ctx, span := otel.Tracer("worker-test").Start(t.Context(), "scheduled-job")
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	assert.NotEmpty(t, carrier.Get("traceparent"))
	span.End()
	require.NoError(t, tp.ForceFlush(t.Context()))
	select {
	case body := <-requests:
		assert.Contains(t, string(body), "worker-test")
		assert.Contains(t, string(body), "scheduled-job")
	default:
		t.Fatal("no trace export received")
	}
}

func TestMetricsServer_ExposesClientMetrics(t *testing.T) {
	transport := observability.NewTransport("worker-metrics-test", observability.WithRoundTripper(
		&mockRoundTripper{resp: newMockResponse(http.StatusOK)},
	))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	server := observability.NewDebugServer(":6060", http.NotFound, http.NotFound)
	recorder := httptest.NewRecorder()
	server.Handler.ServeHTTP(recorder, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/metrics", nil))
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "worker-metrics-test")
}

func TestNewDebugServer_ServesProbes(t *testing.T) {
	probe := func(status int) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(status) }
	}
	server := observability.NewDebugServer(":6060", probe(http.StatusOK), probe(http.StatusServiceUnavailable))

	for path, want := range map[string]int{
		"/liveness": http.StatusOK, "/readiness": http.StatusServiceUnavailable, "/debug/pprof/": http.StatusOK, "/debug/vars": http.StatusOK,
	} {
		recorder := httptest.NewRecorder()
		server.Handler.ServeHTTP(recorder, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))
		assert.Equal(t, want, recorder.Code, path)
	}
}
