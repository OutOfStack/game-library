package api

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/gin-gonic/gin"
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
func (p *Provider) CreateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.createGame")
	defer span.End()

	var cg api.CreateGameRequest
	err := web.Decode(c, &cg)
	if err != nil {
		web.Err(c, fmt.Errorf("decoding new game: %w", err))
		return
	}
	span.SetAttributes(att.String("data.name", cg.Name))

	claims, err := web.GetClaims(c)
	if err != nil {
		web.Err(c, fmt.Errorf("getting claims from context: %w", err))
		return
	}

	developer, publisher := cg.Developer, claims.Name

	create := mapToCreateGame(&cg, developer, publisher)

	id, err := p.gameFacade.CreateGame(ctx, create)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.Err(c, web.NewRequestError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("create game", zap.String("name", create.Name), zap.String("user_id", claims.UserID()), zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	web.Respond(c, api.IDResponse{ID: id}, http.StatusCreated)
}
