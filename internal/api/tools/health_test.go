package tools_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OutOfStack/game-library/internal/api/tools"
	"github.com/OutOfStack/game-library/internal/version"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck_Liveness(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
	}{
		{name: "without Kubernetes metadata", env: map[string]string{
			"KUBERNETES_PODNAME": "", "KUBERNETES_PODIP": "", "KUBERNETES_NODENAME": "", "KUBERNETES_NAMESPACE": "",
		}},
		{name: "with Kubernetes metadata", env: map[string]string{
			"KUBERNETES_PODNAME": "api-1", "KUBERNETES_PODIP": "10.0.0.1", "KUBERNETES_NODENAME": "node-1", "KUBERNETES_NAMESPACE": "games",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			response := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/liveness", nil)
			tools.NewHealthCheck(nil).Liveness(response, req)

			require.Equal(t, http.StatusOK, response.Code)
			body := checkHealthResponse(t, response, "OK")
			for field, env := range map[string]string{
				"pod": "KUBERNETES_PODNAME", "podIP": "KUBERNETES_PODIP", "node": "KUBERNETES_NODENAME", "namespace": "KUBERNETES_NAMESPACE",
			} {
				if tt.env[env] == "" {
					require.NotContains(t, body, field)
				} else {
					require.Equal(t, tt.env[env], body[field])
				}
			}
		})
	}
}

func TestHealthCheck_Readiness_DatabaseUnavailable(t *testing.T) {
	pool, err := pgxpool.New(t.Context(), "postgres://test:test@localhost/test?sslmode=disable&pool_min_conns=0")
	require.NoError(t, err)
	// a closed pool makes Ping fail deterministically without contacting PostgreSQL
	pool.Close()

	response := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readiness", nil)
	hc := tools.NewHealthCheck(pool)
	hc.SetReady(true)
	hc.Readiness(response, req)

	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	checkHealthResponse(t, response, "database not ready")
}

func TestHealthCheck_Readiness_NotReady(t *testing.T) {
	for _, setReady := range []func(*tools.HealthCheck){
		func(*tools.HealthCheck) {},
		func(hc *tools.HealthCheck) { hc.SetReady(true); hc.SetReady(false) },
	} {
		response := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readiness", nil)
		// nil db: not-ready check must short-circuit before the database ping
		hc := tools.NewHealthCheck(nil)
		setReady(hc)
		hc.Readiness(response, req)

		require.Equal(t, http.StatusServiceUnavailable, response.Code)
		checkHealthResponse(t, response, "not ready")
	}
}

func checkHealthResponse(t *testing.T, response *httptest.ResponseRecorder, status string) map[string]string {
	t.Helper()
	require.Contains(t, response.Header().Get("Content-Type"), "application/json")
	var body map[string]string
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, status, body["status"])
	host, err := os.Hostname()
	require.NoError(t, err)
	require.Equal(t, host, body["host"])
	var info version.Info
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &info))
	require.Equal(t, version.Get(), info)
	return body
}
