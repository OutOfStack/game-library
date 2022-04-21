package handler

import (
	"net/http"
	"os"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Check has methods for readiness and liveness probes
type Check struct {
	DB *sqlx.DB
}

type health struct {
	Status string `json:"status"`
	Host   string `json:"host"`
}

// Readiness determines whether service is ready
func (ch *Check) Readiness(c *gin.Context) {
	var h health
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	h.Host = host
	err = database.StatusCheck(ch.DB)
	if err != nil {
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
		host = "unavailable"
	}
	h := health{
		Host:   host,
		Status: "OK",
	}

	web.Respond(c, h, http.StatusOK)
}
