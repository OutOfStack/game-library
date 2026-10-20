package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/web"
	"go.uber.org/zap"
)

const (
	minLengthForSearch = 2
)

var (
	errInvalidID          = web.NewErrorFromMessage("invalid id", http.StatusBadRequest)
	errInvalidCompanyType = web.NewErrorFromMessage("invalid company type: should be one of [dev, pub]", http.StatusBadRequest)
	errInvalidQueryParams = web.NewErrorFromMessage("invalid query params", http.StatusBadRequest)
	errInvalidRequest     = web.NewErrorFromMessage("invalid request", http.StatusBadRequest)
)

var gameModeNames = map[int32]string{
	1: "Single player",
	2: "Multiplayer",
	3: "Co-operative",
	4: "Split screen",
	5: "MMO",
	6: "Battle Royale",
}

// Mappings

func mapToCreateGame(cgr *CreateGameRequest, publisher string) model.CreateGame {
	return model.CreateGame{
		Name:         cgr.Name,
		ReleaseDate:  cgr.ReleaseDate,
		GenresIDs:    cgr.GenresIDs,
		LogoURL:      cgr.LogoURL,
		Summary:      cgr.Summary,
		Slug:         model.GetGameSlug(cgr.Name),
		PlatformsIDs: cgr.PlatformsIDs,
		Screenshots:  cgr.Screenshots,
		Websites:     cgr.Websites,
		Developer:    cgr.Developer,
		Publisher:    publisher,
	}
}

func (p *Provider) mapToGameResponse(ctx context.Context, game model.Game) (GameResponse, error) {
	resp := GameResponse{
		ID:          game.ID,
		Name:        game.Name,
		ReleaseDate: game.ReleaseDate.String(),
		LogoURL:     game.LogoURL,
		Rating:      game.Rating,
		Summary:     game.Summary,
		Slug:        game.Slug,
		Screenshots: game.Screenshots,
		Websites:    game.Websites,
	}

	genres, err := p.gameFacade.GetGenresMap(ctx)
	if err != nil {
		return GameResponse{}, fmt.Errorf("get genres: %v", err)
	}
	companies, err := p.gameFacade.GetCompaniesMap(ctx)
	if err != nil {
		return GameResponse{}, fmt.Errorf("get companies: %v", err)
	}
	platforms, err := p.gameFacade.GetPlatformsMap(ctx)
	if err != nil {
		return GameResponse{}, fmt.Errorf("get platforms: %v", err)
	}

	for _, gID := range game.GenresIDs {
		genre, ok := genres[gID]
		if !ok {
			return GameResponse{}, fmt.Errorf("genre %d not found", gID)
		}
		resp.Genres = append(resp.Genres, Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}
	for _, plID := range game.PlatformsIDs {
		pl, ok := platforms[plID]
		if !ok {
			return GameResponse{}, fmt.Errorf("platform %d not found", plID)
		}
		resp.Platforms = append(resp.Platforms, Platform{
			ID:           pl.ID,
			Name:         pl.Name,
			Abbreviation: pl.Abbreviation,
		})
	}
	for _, cID := range game.DevelopersIDs {
		c, ok := companies[cID]
		if !ok {
			return GameResponse{}, fmt.Errorf("company %d not found", cID)
		}
		resp.Developers = append(resp.Developers, Company{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	for _, cID := range game.PublishersIDs {
		c, ok := companies[cID]
		if !ok {
			return GameResponse{}, fmt.Errorf("company %d not found", cID)
		}
		resp.Publishers = append(resp.Publishers, Company{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	return resp, nil
}

func (p *Provider) mapToGameDetailsResponse(ctx context.Context, details model.GameDetails) (GameDetailsResponse, error) {
	game, err := p.mapToGameResponse(ctx, details.Game)
	if err != nil {
		return GameDetailsResponse{}, err
	}

	resp := GameDetailsResponse{
		ID:             game.ID,
		Name:           game.Name,
		Developers:     game.Developers,
		Publishers:     game.Publishers,
		ReleaseDate:    game.ReleaseDate,
		Genres:         game.Genres,
		LogoURL:        game.LogoURL,
		Rating:         game.Rating,
		Summary:        game.Summary,
		Slug:           game.Slug,
		Platforms:      game.Platforms,
		Screenshots:    game.Screenshots,
		Websites:       game.Websites,
		Trailers:       append([]string{}, details.Trailers...),
		GameModes:      make([]GameMode, 0, len(details.GameModeIDs)),
		SimilarGames:   make([]GameSummary, 0, len(details.SimilarGames)),
		FranchiseGames: make([]GameSummary, 0, len(details.FranchiseGames)),
	}
	for _, id := range details.GameModeIDs {
		if name, ok := gameModeNames[id]; ok {
			resp.GameModes = append(resp.GameModes, GameMode{ID: id, Name: name})
		}
	}
	for _, summary := range details.SimilarGames {
		resp.SimilarGames = append(resp.SimilarGames, mapToGameSummary(summary))
	}
	for _, summary := range details.FranchiseGames {
		resp.FranchiseGames = append(resp.FranchiseGames, mapToGameSummary(summary))
	}
	return resp, nil
}

func mapToGameSummary(game model.GameSummary) GameSummary {
	return GameSummary{
		ID:          game.ID,
		Name:        game.Name,
		ReleaseDate: game.ReleaseDate.String(),
		LogoURL:     game.LogoURL,
		Rating:      game.Rating,
		Slug:        game.Slug,
	}
}

func mapToUpdateGame(ugr *UpdateGameRequest, publisher string) model.UpdateGame {
	return model.UpdateGame{
		Name:         ugr.Name,
		Developer:    ugr.Developer,
		Publisher:    publisher,
		ReleaseDate:  ugr.ReleaseDate,
		GenresIDs:    ugr.GenresIDs,
		LogoURL:      ugr.LogoURL,
		Summary:      ugr.Summary,
		PlatformsIDs: ugr.PlatformsIDs,
		Screenshots:  ugr.Screenshots,
		Websites:     ugr.Websites,
	}
}

func mapToGamesFilter(p GetGamesParams) (model.GamesFilter, error) {
	if p.Page == 0 || p.PageSize == 0 {
		return model.GamesFilter{}, fmt.Errorf("invalid page or page size param: should be greater than 0")
	}

	var filter model.GamesFilter
	switch p.OrderBy {
	case "", GetGamesParamsOrderByDefault:
		filter.OrderBy = model.OrderGamesByDefault
	case GetGamesParamsOrderByName:
		filter.OrderBy = model.OrderGamesByName
	case GetGamesParamsOrderByReleaseDate:
		filter.OrderBy = model.OrderGamesByReleaseDate
	case GetGamesParamsOrderByRating:
		filter.OrderBy = model.OrderGamesByRating
	default:
		return model.GamesFilter{}, fmt.Errorf("invalid orderBy: should be one of [default, releaseDate, name, rating]")
	}
	if len(p.Name) >= minLengthForSearch {
		filter.Name = strings.ToLower(p.Name)
	}
	filter.GenreID = p.Genre
	filter.DeveloperID = p.Developer
	filter.PublisherID = p.Publisher

	return filter, nil
}

// decodeRequest decodes JSON request body, validates and sanitizes it
func (p *Provider) decodeRequest(r *http.Request, req validatable) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		p.log.Error("decode request", zap.Error(err))
		return errInvalidRequest
	}
	return p.validateRequest(req)
}

// respondError responds with app error status code if err is an app error,
// otherwise logs err with msg and fields and responds with 500 Internal Server Error
func (p *Provider) respondError(w http.ResponseWriter, msg string, err error, fields ...zap.Field) {
	if appErr, ok := apperr.IsAppError(err); ok {
		web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
		return
	}
	p.log.Error(msg, append(fields, zap.Error(err))...)
	web.Respond500(w)
}
