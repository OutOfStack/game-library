package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// CreateGame godoc
// @Summary Create game
// @Description creates new game
// @ID create-game
// @Accept  json
// @Produce json
// @Param   game body 	 api.CreateGameRequest true "create game"
// @Success 201 {object} api.IDResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [post]
func (p *Provider) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "api.createGame")
	defer span.End()

	var cg api.CreateGameRequest
	if err := web.Decode(p.log, r, &cg); err != nil {
		web.RespondError(w, err)
		return
	}
	span.SetAttributes(att.String("data.name", cg.Name))

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.log.Error("get claims from context", zap.Error(err))
		web.Respond500(w)
		return
	}

	developer, publisher := cg.Developer, claims.Name

	create := mapToCreateGame(&cg, developer, publisher)

	id, err := p.gameFacade.CreateGame(ctx, create)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("create game", zap.String("name", create.Name), zap.String("user_id", claims.UserID()), zap.Error(err))
		web.Respond500(w)
		return
	}

	web.Respond(w, api.IDResponse{ID: id}, http.StatusCreated)
}
