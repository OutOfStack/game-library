package handler

import (
	"log"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/jmoiron/sqlx"
)

// Service constructs that contains all API routes
func Service(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger, middleware.Errors(logger))

	//app.UseMiddleware()

	c := Check{
		DB: db,
	}
	app.Handle(http.MethodGet, "/api/health", c.Health)

	g := Game{
		DB:  db,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/api/games", g.List)
	app.Handle(http.MethodGet, "/api/games/:id", g.Retrieve)
	app.Handle(http.MethodPost, "/api/games", g.Create)
	app.Handle(http.MethodPatch, "/api/games/:id", g.Update)
	app.Handle(http.MethodDelete, "/api/games/:id", g.Delete)

	app.Handle(http.MethodPost, "/api/games/:id/sales", g.AddSale)
	app.Handle(http.MethodGet, "api/games/:id/sales", g.ListSales)

	return app
}
