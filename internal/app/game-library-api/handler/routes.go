package handler

import (
	"log"
	"net/http"

	_ "github.com/OutOfStack/game-library/docs"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Service constructs router with all API routes
func Service(logger *log.Logger, db *sqlx.DB, a *auth.Auth) http.Handler {
	r := gin.Default()
	r.Use(middleware.Errors(logger), middleware.Metrics())

	c := Check{
		DB: db,
	}

	g := Game{
		DB:  db,
		Log: logger,
	}

	// readiness
	r.GET("/api/health", c.Health)

	// games
	r.GET("/api/games", g.List)
	r.GET("/api/games/:id", g.Retrieve)
	// authorization required
	r.POST("/api/games", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.Create)
	r.DELETE("/api/games/:id", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.Delete)
	r.PATCH("/api/games/:id", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.Update)
	r.POST("/api/games/:id/sales", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.AddGameOnSale)
	r.GET("/api/games/:id/sales", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.ListGameSales)
	r.POST("/api/games/rate", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleRegisteredUser), g.RateGame)

	// sales
	r.GET("/api/sales", g.ListSales)
	// authorization required
	r.POST("/api/sales", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleModerator), g.AddSale)

	// user
	// authorization required
	r.POST("/api/user/ratings", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleRegisteredUser), g.GetUserRatings)

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
