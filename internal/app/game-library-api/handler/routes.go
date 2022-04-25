package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/OutOfStack/game-library/docs"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	gintrace "go.opentelemetry.io/contrib/instrumentation/gin-gonic/gin"
	global "go.opentelemetry.io/otel/api/global"
	zipkin "go.opentelemetry.io/otel/exporters/trace/zipkin"
	trace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	serviceName = "game-library-api"
)

var tracer = global.Tracer(serviceName)

// Service constructs router with all API routes
func Service(logger *log.Logger, db *sqlx.DB, a *auth.Auth, conf appconf.Web, zipkinConf appconf.Zipkin) (http.Handler, error) {
	err := initTracer(zipkinConf.ReporterURL)
	if err != nil {
		return nil, fmt.Errorf("initializing exporter: %w", err)
	}
	r := gin.Default()
	r.Use(gintrace.Middleware("game-library-api"))
	r.Use(middleware.Errors(logger), middleware.Metrics(), cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(conf.AllowedCORSOrigin, origin)
		},
	}))

	c := Check{
		DB: db,
	}

	g := Game{
		DB:  db,
		Log: logger,
	}

	// health
	r.GET("/api/readiness", c.Readiness)
	r.GET("/api/liveness", c.Liveness)

	// games
	r.GET("/api/games", g.GetList)
	r.GET("/api/games/:id", g.Get)
	r.GET("/api/games/search", g.Search)
	// authorization required
	r.POST("/api/games", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.Create)
	r.DELETE("/api/games/:id", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.Delete)
	r.PATCH("/api/games/:id", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.Update)
	r.POST("/api/games/:id/sales", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.AddGameOnSale)
	r.GET("/api/games/:id/sales", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher), g.ListGameSales)
	r.POST("/api/games/:id/rate", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleRegisteredUser), g.RateGame)

	// sales
	r.GET("/api/sales", g.ListSales)
	// authorization required
	r.POST("/api/sales", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleModerator), g.AddSale)

	// user
	// authorization required
	r.POST("/api/user/ratings", middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleRegisteredUser), g.GetUserRatings)

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r, nil
}

func initTracer(reporterURL string) error {
	exporter, err := zipkin.NewExporter(reporterURL, serviceName)
	if err != nil {
		return errors.Wrap(err, "creating new exporter")
	}
	cfg := trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	}
	tp, err := trace.NewProvider(
		trace.WithConfig(cfg),
		trace.WithBatcher(exporter),
	)
	if err != nil {
		return errors.Wrap(err, "creating new trace provider")
	}
	global.SetTraceProvider(tp)
	return nil
}
