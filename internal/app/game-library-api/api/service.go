package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/OutOfStack/game-library/docs" // swagger docs
	"github.com/OutOfStack/game-library/internal/app/game-library-api/api/tools"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	chicors "github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	swag "github.com/swaggo/http-swagger/v2"
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
	db *pgxpool.Pool,
	au *auth.Client,
	pr *Provider,
	conf appconf.Cfg,
) (http.Server, error) {
	err := initTracer(log, conf.Zipkin.ReporterURL)
	if err != nil {
		return http.Server{}, fmt.Errorf("initializing exporter: %w", err) //nolint:gosec
	}

	r := chi.NewRouter()
	r.Use(mw.RequestID)
	r.Use(middleware.Metrics)
	r.Use(middleware.Logger(log))
	r.Use(mw.Recoverer)
	r.Use(otelchi.Middleware(appconf.ServiceName))
	r.Use(chicors.Handler(chicors.Options{
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-type", "Authorization"},
		AllowCredentials: true,
		AllowedOrigins:   strings.Split(conf.Web.AllowedCORSOrigin, ","),
	}))

	hc := tools.NewHealthCheck(db)

	// metrics
	r.Handle("/metrics", promhttp.Handler())

	// tools
	r.Get("/api/readiness", hc.Readiness)

	r.Get("/api/liveness", hc.Liveness)

	// games
	r.Route("/api/games", func(r chi.Router) {
		r.Get("/", pr.GetGames)

		r.Get("/{id}", pr.GetGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Post("/", pr.CreateGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Delete("/{id}", pr.DeleteGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Patch("/{id}", pr.UpdateGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RoleRegisteredUser),
		).Post("/{id}/rate", pr.RateGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Post("/images", pr.UploadGameImages)
	})

	// user
	r.With(
		middleware.Authenticate(log, au),
		middleware.Authorize(log, au, auth.RoleRegisteredUser),
	).Post("/api/user/ratings", pr.GetUserRatings)

	// genres
	r.Route("/api/genres", func(r chi.Router) {
		r.Get("/", pr.GetGenres)

		r.Get("/top", pr.GetTopGenres)
	})

	// platforms
	r.Get("/api/platforms", pr.GetPlatforms)

	// companies
	r.Get("/api/companies/top", pr.GetTopCompanies)

	// swagger
	r.Get("/swagger/*", swag.Handler())

	return http.Server{
		Addr:              conf.Web.Address,
		Handler:           r,
		ReadTimeout:       conf.Web.ReadTimeout,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      conf.Web.WriteTimeout,
	}, nil
}

func initTracer(log *zap.Logger, reporterURL string) error {
	exporter, err := zipkin.New(reporterURL)
	if err != nil {
		return fmt.Errorf("create new exporter: %v", err)
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
		log.Error("zipkin error", zap.Error(err))
	}))
	otel.SetTracerProvider(tp)

	return nil
}
