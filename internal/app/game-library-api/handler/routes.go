package handler

import (
	"log"
	"net/http"

	_ "github.com/OutOfStack/game-library/docs"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Service constructs that contains all API routes
func Service(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger, middleware.Errors(logger), middleware.Metrics())

	c := Check{
		DB: db,
	}
	app.Handle(http.MethodGet, "/api/health", c.Health)

	g := Game{
		DB:  db,
		Log: logger,
	}

	// games
	app.Handle(http.MethodGet, "/api/games", g.List)
	app.Handle(http.MethodGet, "/api/games/:id", g.Retrieve)
	app.Handle(http.MethodPost, "/api/games", g.Create)
	app.Handle(http.MethodPatch, "/api/games/:id", g.Update)
	app.Handle(http.MethodDelete, "/api/games/:id", g.Delete)
	app.Handle(http.MethodPost, "/api/games/:id/sales", g.AddGameOnSale)
	app.Handle(http.MethodGet, "/api/games/:id/sales", g.ListGameSales)

	// sales
	app.Handle(http.MethodPost, "/api/sales", g.AddSale)
	app.Handle(http.MethodGet, "/api/sales", g.ListSales)

	app.Handle(http.MethodGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return app
}
