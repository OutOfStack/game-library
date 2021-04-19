package handler

import (
	"context"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Check has methods for health checking
type Check struct {
	DB *sqlx.DB
}

// Health determines whether service is healthy
func (che *Check) Health(ctx context.Context, c *gin.Context) error {
	var health struct {
		Status string `json:"status"`
	}

	err := database.StatusCheck(che.DB)
	if err != nil {
		health.Status = "database not ready"
		return web.Respond(ctx, c, health, http.StatusInternalServerError)
	}
	health.Status = "OK"
	return web.Respond(ctx, c, health, http.StatusOK)
}
