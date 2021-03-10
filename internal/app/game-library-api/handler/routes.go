package handler

import (
	"log"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/jmoiron/sqlx"
)

// Service constructs that contains all API routes
func Service(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)

	g := Game{
		DB:  db,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/api/games", g.List)
	app.Handle(http.MethodGet, "/api/games/:id", g.Retrieve)
	app.Handle(http.MethodPost, "/api/games", g.Create)

	app.Handle(http.MethodPost, "/api/games/:id/sales", g.AddSale)
	app.Handle(http.MethodGet, "api/games/:id/sales", g.ListSales)

	return app
}
