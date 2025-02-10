package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/go-playground/form/v4"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetGames godoc
// @Summary Get games
// @Description returns paginated games
// @ID get-games
// @Produce json
// @Param pageSize  query int32  false "page size"
// @Param page      query int32  false "page"
// @Param orderBy   query string false "order by"	Enums(default, name, releaseDate)
// @Param name 	    query string false "name filter"
// @Param genre     query int32  false "genre filter"
// @Param developer query int32  false "developer id filter"
// @Param publisher query int32  false "publisher id filter"
// @Success 200 {object}  api.GamesResponse
// @Failure 500 {object}  web.ErrorResponse
// @Router /games [get]
func (p *Provider) GetGames(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "api.getGames")
	defer span.End()

	for key, values := range r.URL.Query() {
		if len(values) == 1 {
			span.SetAttributes(att.String(key, values[0]))
		} else if len(values) > 1 {
			span.SetAttributes(att.StringSlice(key, values))
		}
	}

	// get query params
	var params api.GetGamesQueryParams
	decoder := form.NewDecoder()
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		web.RespondError(w, web.NewErrorFromMessage("invalid query params", http.StatusBadRequest))
		return
	}

	filter, err := mapToGamesFilter(&params)
	if err != nil {
		web.RespondError(w, web.NewErrorFromMessage(err.Error(), http.StatusBadRequest))
		return
	}

	list, count, err := p.gameFacade.GetGames(ctx, params.Page, params.PageSize, filter)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("get games", zap.Error(err))
		web.Respond500(w)
		return
	}

	games := make([]api.GameResponse, 0, len(list))
	for _, game := range list {
		gr, mErr := p.mapToGameResponse(ctx, game)
		if mErr != nil {
			p.log.Error("map game to response", zap.Error(mErr))
			web.RespondError(w, web.NewErrorFromMessage("error converting response", http.StatusInternalServerError))
			return
		}
		games = append(games, gr)
	}

	response := api.GamesResponse{
		Games: games,
		Count: count,
	}
	web.Respond(w, response, http.StatusOK)
}
