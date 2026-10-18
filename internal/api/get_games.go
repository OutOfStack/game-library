package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetGames returns paginated games
func (p *Provider) GetGames(w http.ResponseWriter, r *http.Request, params GetGamesParams) {
	ctx, span := tracer.Start(r.Context(), "getGames", trace.WithAttributes(gamesParamsAttributes(params)...))
	defer span.End()

	filter, err := mapToGamesFilter(params)
	if err != nil {
		web.RespondError(w, web.NewErrorFromMessage(err.Error(), http.StatusBadRequest))
		return
	}

	list, count, err := p.gameFacade.GetGames(ctx, params.Page, params.PageSize, filter)
	if err != nil {
		p.respondError(w, "get games", err)
		return
	}

	games := make([]GameResponse, 0, len(list))
	for _, game := range list {
		gr, mErr := p.mapToGameResponse(ctx, game)
		if mErr != nil {
			p.respondError(w, "map game to response", mErr)
			return
		}
		games = append(games, gr)
	}

	web.Respond(w, GamesResponse{
		Games: games,
		Count: count,
	}, http.StatusOK)
}

func gamesParamsAttributes(params GetGamesParams) []att.KeyValue {
	return []att.KeyValue{
		att.Int64("pageSize", int64(params.PageSize)),
		att.Int64("page", int64(params.Page)),
		att.String("orderBy", string(params.OrderBy)),
		att.String("name", params.Name),
		att.Int("genre", int(params.Genre)),
		att.Int("developer", int(params.Developer)),
		att.Int("publisher", int(params.Publisher)),
	}
}
