package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/OutOfStack/game-library/docs" // swagger docs
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Service constructs router with all API routes
func Service(logger *log.Logger, db *sqlx.DB, a *auth.Auth, igdb *igdb.Client, conf appconf.Web, zipkinConf appconf.Zipkin) (http.Handler, error) {
	err := initTracer(zipkinConf.ReporterURL)
	if err != nil {
		return nil, fmt.Errorf("initializing exporter: %w", err)
	}
	r := gin.Default()
	r.Use(otelgin.Middleware(appconf.ServiceName))
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
		Log: logger,
		Storage: &repo.Storage{
			DB: db,
		},
		IGDB: igdb,
	}

	// health
	r.GET("/api/readiness", c.Readiness)
	r.GET("/api/liveness", c.Liveness)

	// games
	r.GET("/api/games", g.GetGames)
	r.GET("/api/games/:id", g.GetGame)
	r.GET("/api/games/search", g.SearchGames)
	r.POST("/api/games",
		middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher),
		g.CreateGame)
	r.DELETE("/api/games/:id",
		middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher),
		g.DeleteGame)
	r.PATCH("/api/games/:id",
		middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RolePublisher),
		g.UpdateGame)
	r.POST("/api/games/:id/rate",
		middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleRegisteredUser),
		g.RateGame)

	// user
	r.POST("/api/user/ratings",
		middleware.Authenticate(logger, a), middleware.Authorize(logger, a, auth.RoleRegisteredUser),
		g.GetUserRatings)

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r, nil
}

func initTracer(reporterURL string) error {
	exporter, err := zipkin.New(reporterURL)
	if err != nil {
		return errors.Wrap(err, "creating new exporter")
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(appconf.ServiceName),
			)),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return nil
}
