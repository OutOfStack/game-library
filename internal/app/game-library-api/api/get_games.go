package api

import (
	"errors"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/gin-gonic/gin"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	minLengthForSearch = 2
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
func (p *Provider) GetGames(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.getGames")
	defer span.End()

	for key, values := range c.Request.URL.Query() {
		if len(values) == 1 {
			span.SetAttributes(att.String(key, values[0]))
		} else if len(values) > 1 {
			span.SetAttributes(att.StringSlice(key, values))
		}
	}

	// get query params and form filter
	var queryParams api.GetGamesQueryParams
	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		web.Err(c, web.NewRequestError(errors.New("incorrect query params"), http.StatusBadRequest))
		return
	}

	page, pageSize := queryParams.Page, queryParams.PageSize
	var filter model.GamesFilter
	if len(queryParams.Name) >= minLengthForSearch {
		filter.Name = queryParams.Name
	}
	if queryParams.Genre != 0 {
		filter.GenreID = queryParams.Genre
	}
	if queryParams.Developer != 0 {
		filter.DeveloperID = queryParams.Developer
	}
	if queryParams.Publisher != 0 {
		filter.PublisherID = queryParams.Publisher
	}
	switch queryParams.OrderBy {
	case "", "default":
		filter.OrderBy = repo.OrderGamesByDefault
	case "name":
		filter.OrderBy = repo.OrderGamesByName
	case "releaseDate":
		filter.OrderBy = repo.OrderGamesByReleaseDate
	default:
		web.Err(c, web.NewRequestError(errors.New("incorrect orderBy. Should be one of: default, releaseDate, name"), http.StatusBadRequest))
		return
	}

	list, count, err := p.gameFacade.GetGames(ctx, page, pageSize, filter)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.Err(c, web.NewRequestError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("get games", zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	games := make([]api.GameResponse, 0, len(list))
	for _, game := range list {
		r, mErr := p.mapToGameResponse(c, game)
		if mErr != nil {
			web.Err(c, web.NewRequestError(errors.New("error converting response"), http.StatusInternalServerError))
			return
		}
		games = append(games, r)
	}

	response := api.GamesResponse{
		Games: games,
		Count: count,
	}
	web.Respond(c, response, http.StatusOK)
}
