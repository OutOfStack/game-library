package handler

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/OutOfStack/game-library/docs" // swagger docs
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/appconf"
	auth_ "github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.uber.org/zap"
)

// Service constructs router with all API routes
func Service(
	logger *zap.Logger,
	db *sqlx.DB,
	auth *auth_.Auth,
	storage *repo.Storage,
	cache *cache.RedisStore,
	conf appconf.Web,
	zipkinConf appconf.Zipkin,
) (http.Handler, error) {
	err := initTracer(logger, zipkinConf.ReporterURL)
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

	c := NewCheck(db)

	pr := NewProvider(logger, storage, cache)

	// health
	r.GET("/api/readiness", c.Readiness)
	r.GET("/api/liveness", c.Liveness)

	// games
	r.GET("/api/games", pr.GetGames)
	r.GET("/api/games/:id", pr.GetGame)
	r.POST("/api/games",
		middleware.Authenticate(logger, auth), middleware.Authorize(logger, auth, auth_.RolePublisher),
		pr.CreateGame)
	r.DELETE("/api/games/:id",
		middleware.Authenticate(logger, auth), middleware.Authorize(logger, auth, auth_.RolePublisher),
		pr.DeleteGame)
	r.PATCH("/api/games/:id",
		middleware.Authenticate(logger, auth), middleware.Authorize(logger, auth, auth_.RolePublisher),
		pr.UpdateGame)
	r.POST("/api/games/:id/rate",
		middleware.Authenticate(logger, auth), middleware.Authorize(logger, auth, auth_.RoleRegisteredUser),
		pr.RateGame)

	// user
	r.POST("/api/user/ratings",
		middleware.Authenticate(logger, auth), middleware.Authorize(logger, auth, auth_.RoleRegisteredUser),
		pr.GetUserRatings)

	// genres
	r.GET("/api/genres", pr.GetGenres)
	r.GET("/api/genres/top", pr.GetTopGenres)

	// platforms
	r.GET("/api/platforms", pr.GetPlatforms)

	// companies
	r.GET("/api/companies/top", pr.GetTopCompanies)

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r, nil
}

func initTracer(logger *zap.Logger, reporterURL string) error {
	exporter, err := zipkin.New(reporterURL)
	if err != nil {
		return fmt.Errorf("creating new exporter: %v", err)
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
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error("zipkin error", zap.Error(err))
	}))
	otel.SetTracerProvider(tp)

	return nil
}
