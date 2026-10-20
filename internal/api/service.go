package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/OutOfStack/game-library/api"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/OutOfStack/game-library/internal/web"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	chicors "github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// Service constructs router with all API routes
func Service(
	log *zap.Logger,
	au *auth.Client,
	pr *Provider,
	conf *appconf.API,
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

	h := NewHandler(pr)

	// games
	r.Route("/api/games", func(r chi.Router) {
		r.Get("/", h.GetGames)

		r.Get("/{id}", h.GetGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Post("/", h.CreateGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Delete("/{id}", h.DeleteGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Patch("/{id}", h.UpdateGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RoleRegisteredUser),
		).Post("/{id}/rate", h.RateGame)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Post("/images", h.UploadGameImages)

		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Get("/{id}/moderations", h.GetGameModerations)
	})

	// user
	r.Route("/api/user", func(r chi.Router) {
		// ratings
		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RoleRegisteredUser),
		).Post("/ratings", h.GetUserRatings)

		// published games
		r.With(
			middleware.Authenticate(log, au),
			middleware.Authorize(log, au, auth.RolePublisher),
		).Get("/games", h.GetUserGames)
	})

	// genres
	r.Route("/api/genres", func(r chi.Router) {
		r.Get("/", h.GetGenres)

		r.Get("/top", h.GetTopGenres)
	})

	// platforms
	r.Get("/api/platforms", h.GetPlatforms)

	// companies
	r.Get("/api/companies/top", h.GetTopCompanies)

	// swagger
	r.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(api.OpenAPISpec)
	})
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/openapi.yaml")))

	return http.Server{
		Addr:              conf.Web.HTTPAddress,
		Handler:           r,
		ReadTimeout:       conf.Web.ReadTimeout,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      conf.Web.WriteTimeout,
	}, tp, nil
}

// NewHandler returns HTTP handlers for API operations generated from OpenAPI spec
func NewHandler(pr *Provider) *ServerInterfaceWrapper {
	return &ServerInterfaceWrapper{
		Handler: pr,
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			web.RespondError(w, paramError(err))
		},
	}
}

// paramError maps request parameter binding error to web error
func paramError(err error) error {
	var name string
	if formatErr, isFormatErr := errors.AsType[*InvalidParamFormatError](err); isFormatErr {
		name = formatErr.ParamName
	} else if requiredErr, isRequiredErr := errors.AsType[*RequiredParamError](err); isRequiredErr {
		name = requiredErr.ParamName
	}

	switch name {
	case "id":
		return errInvalidID
	case "type": //nolint:goconst // param name from OpenAPI spec
		return errInvalidCompanyType
	default:
		return errInvalidQueryParams
	}
}
