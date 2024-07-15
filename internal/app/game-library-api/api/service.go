package api

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/OutOfStack/game-library/docs" // swagger docs
	"github.com/OutOfStack/game-library/internal/app/game-library-api/api/tools"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
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
	log *zap.Logger,
	db *sqlx.DB,
	au *auth.Client,
	pr *Provider,
	conf appconf.Cfg,
) (http.Server, error) {
	err := initTracer(log, conf.Zipkin.ReporterURL)
	if err != nil {
		return http.Server{}, fmt.Errorf("initializing exporter: %w", err)
	}
	r := gin.Default()
	r.Use(otelgin.Middleware(appconf.ServiceName))
	r.Use(middleware.Errors(log), middleware.Metrics(), cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(conf.Web.AllowedCORSOrigin, origin)
		},
	}))

	hc := tools.NewHealthCheck(db)

	// tools
	r.GET("/api/readiness", hc.Readiness)
	r.GET("/api/liveness", hc.Liveness)

	// games
	r.GET("/api/games", pr.GetGames)
	r.GET("/api/games/:id", pr.GetGame)
	r.POST("/api/games",
		middleware.Authenticate(log, au), middleware.Authorize(log, au, auth.RolePublisher),
		pr.CreateGame)
	r.DELETE("/api/games/:id",
		middleware.Authenticate(log, au), middleware.Authorize(log, au, auth.RolePublisher),
		pr.DeleteGame)
	r.PATCH("/api/games/:id",
		middleware.Authenticate(log, au), middleware.Authorize(log, au, auth.RolePublisher),
		pr.UpdateGame)
	r.POST("/api/games/:id/rate",
		middleware.Authenticate(log, au), middleware.Authorize(log, au, auth.RoleRegisteredUser),
		pr.RateGame)

	// user
	r.POST("/api/user/ratings",
		middleware.Authenticate(log, au), middleware.Authorize(log, au, auth.RoleRegisteredUser),
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

	return http.Server{
		Addr:         conf.Web.Address,
		Handler:      r,
		ReadTimeout:  conf.Web.ReadTimeout,
		WriteTimeout: conf.Web.WriteTimeout,
	}, nil
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
