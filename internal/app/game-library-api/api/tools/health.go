package tools

import (
	"net/http"
	"os"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	unavailable = "unavailable"
)

type health struct {
	Status    string `json:"status,omitempty"`
	Host      string `json:"host,omitempty"`
	Pod       string `json:"pod,omitempty"`
	PodIP     string `json:"podIP,omitempty"`
	Node      string `json:"node,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// HealthCheck has methods for readiness and liveness probes
type HealthCheck struct {
	db *sqlx.DB
}

// NewHealthCheck creates new HealthCheck
func NewHealthCheck(db *sqlx.DB) *HealthCheck {
	return &HealthCheck{db: db}
}

// Readiness determines whether service is ready
func (hc *HealthCheck) Readiness(c *gin.Context) {
	var h health
	host, err := os.Hostname()
	if err != nil {
		host = unavailable
	}
	h.Host = host
	if err = hc.db.Ping(); err != nil {
		h.Status = "database not ready"
		web.Respond(c, h, http.StatusInternalServerError)
		return
	}
	h.Status = "OK"
	web.Respond(c, h, http.StatusOK)
}

// Liveness determines whether service is up
func (hc *HealthCheck) Liveness(c *gin.Context) {
	host, err := os.Hostname()
	if err != nil {
		host = unavailable
	}
	h := health{
		Host:      host,
		Status:    "OK",
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_PODIP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	web.Respond(c, h, http.StatusOK)
}
