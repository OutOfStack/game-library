package handler

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

// Check has methods for readiness and liveness probes
type Check struct {
	db *sqlx.DB
}

// NewCheck creates new Check
func NewCheck(db *sqlx.DB) *Check {
	return &Check{db: db}
}

type health struct {
	Status    string `json:"status,omitempty"`
	Host      string `json:"host,omitempty"`
	Pod       string `json:"pod,omitempty"`
	PodIP     string `json:"podIP,omitempty"`
	Node      string `json:"node,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// Readiness determines whether service is ready
func (ch *Check) Readiness(c *gin.Context) {
	var h health
	host, err := os.Hostname()
	if err != nil {
		host = unavailable
	}
	h.Host = host
	if err = ch.db.Ping(); err != nil {
		h.Status = "database not ready"
		web.Respond(c, h, http.StatusInternalServerError)
		return
	}
	h.Status = "OK"
	web.Respond(c, h, http.StatusOK)
}

// Liveness determines whether service is up
func (ch *Check) Liveness(c *gin.Context) {
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
