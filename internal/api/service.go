package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/OutOfStack/game-library/docs" // swagger docs
	"github.com/OutOfStack/game-library/internal/api/tools"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	chicors "github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	swag "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// Service constructs router with all API routes
func Service(
	log *zap.Logger,
	db *pgxpool.Pool,
	au *auth.Client,
	pr *Provider,
	conf *appconf.Cfg,
) (http.Server, *trace.TracerProvider, error) {
	tp, err := observability.InitTracer(log, conf.Jaeger.OTLPEndpoint, appconf.ServiceName)
	if err != nil {
		return http.Server{}, nil, fmt.Errorf("initializing tracer: %w", err) //nolint:gosec
	}

	r := chi.NewRouter()
	r.Use(mw.RequestID)
	r.Use(middleware.Metrics)
	r.Use(middleware.Logger(log))
	r.Use(mw.Recoverer)
	r.Use(otelchi.Middleware(appconf.ServiceName))
	r.Use(chicors.Handler(chicors.Options{
		AllowedOrigins:   strings.Split(conf.Web.AllowedCORSOrigin, ","),
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-type", "Authorization"},
		AllowCredentials: true,
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

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Get("/{id}/moderations", pr.GetGameModerations)
	})

	// user
	r.Route("/api/user", func(r chi.Router) {
		// ratings
		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RoleRegisteredUser),
		).Post("/ratings", pr.GetUserRatings)

		// published games
		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Get("/games", pr.GetUserGames)
	})

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
		Addr:              conf.Web.HTTPAddress,
		Handler:           r,
		ReadTimeout:       conf.Web.ReadTimeout,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      conf.Web.WriteTimeout,
	}, tp, nil
}
