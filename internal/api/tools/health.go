package tools

import (
	"net/http"
	"os"

	"github.com/OutOfStack/game-library/internal/version"
	"github.com/OutOfStack/game-library/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	unavailable = "unavailable"
)

type health struct {
	version.Info
	Status    string `json:"status,omitempty"`
	Host      string `json:"host,omitempty"`
	Pod       string `json:"pod,omitempty"`
	PodIP     string `json:"podIP,omitempty"`
	Node      string `json:"node,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// HealthCheck has methods for readiness and liveness probes
type HealthCheck struct {
	db *pgxpool.Pool
}

// NewHealthCheck creates new HealthCheck
func NewHealthCheck(db *pgxpool.Pool) *HealthCheck {
	return &HealthCheck{db: db}
}

// Readiness determines whether service is ready
func (hc *HealthCheck) Readiness(w http.ResponseWriter, r *http.Request) {
	h := newHealth()
	if err := hc.db.Ping(r.Context()); err != nil {
		h.Status = "database not ready"
		web.Respond(w, h, http.StatusServiceUnavailable)
		return
	}
	h.Status = "OK"
	web.Respond(w, h, http.StatusOK)
}

// Liveness determines whether service is up
func (hc *HealthCheck) Liveness(w http.ResponseWriter, _ *http.Request) {
	h := newHealth()
	h.Status = "OK"
	h.Pod = os.Getenv("KUBERNETES_PODNAME")
	h.PodIP = os.Getenv("KUBERNETES_PODIP")
	h.Node = os.Getenv("KUBERNETES_NODENAME")
	h.Namespace = os.Getenv("KUBERNETES_NAMESPACE")

	web.Respond(w, h, http.StatusOK)
}

func newHealth() health {
	host, err := os.Hostname()
	if err != nil {
		host = unavailable
	}

	return health{
		Info: version.Get(),
		Host: host,
	}
}
