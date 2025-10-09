package api

import (
	"context"
	"errors"
	"fmt"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
)

const (
	minLengthForSearch = 2
)

// Mappings

func mapToCreateGame(cgr *api.CreateGameRequest, publisher string) model.CreateGame {
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

func (p *Provider) mapToGameResponse(ctx context.Context, game model.Game) (api.GameResponse, error) {
	resp := api.GameResponse{
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
		return api.GameResponse{}, fmt.Errorf("get genres: %v", err)
	}
	companies, err := p.gameFacade.GetCompaniesMap(ctx)
	if err != nil {
		return api.GameResponse{}, fmt.Errorf("get companies: %v", err)
	}
	platforms, err := p.gameFacade.GetPlatformsMap(ctx)
	if err != nil {
		return api.GameResponse{}, fmt.Errorf("get platforms: %v", err)
	}

	for _, gID := range game.GenresIDs {
		genre, ok := genres[gID]
		if !ok {
			return api.GameResponse{}, fmt.Errorf("genre %d not found", gID)
		}
		resp.Genres = append(resp.Genres, api.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}
	for _, plID := range game.PlatformsIDs {
		pl, ok := platforms[plID]
		if !ok {
			return api.GameResponse{}, fmt.Errorf("platform %d not found", plID)
		}
		resp.Platforms = append(resp.Platforms, api.Platform{
			ID:           pl.ID,
			Name:         pl.Name,
			Abbreviation: pl.Abbreviation,
		})
	}
	for _, cID := range game.DevelopersIDs {
		c, ok := companies[cID]
		if !ok {
			return api.GameResponse{}, fmt.Errorf("company %d not found", cID)
		}
		resp.Developers = append(resp.Developers, api.Company{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	for _, cID := range game.PublishersIDs {
		c, ok := companies[cID]
		if !ok {
			return api.GameResponse{}, fmt.Errorf("company %d not found", cID)
		}
		resp.Publishers = append(resp.Publishers, api.Company{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	return resp, nil
}

func mapToUpdateGame(ugr *api.UpdateGameRequest, publisher string) model.UpdateGame {
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

func mapToGamesFilter(p *api.GetGamesQueryParams) (model.GamesFilter, error) {
	if p.Page <= 0 || p.PageSize <= 0 {
		return model.GamesFilter{}, errors.New("invalid page or page size param: should be greater than 0")
	}

	var filter model.GamesFilter
	switch p.OrderBy {
	case "", "default":
		filter.OrderBy = repo.OrderGamesByDefault
	case "name":
		filter.OrderBy = repo.OrderGamesByName
	case "releaseDate":
		filter.OrderBy = repo.OrderGamesByReleaseDate
	case "rating":
		filter.OrderBy = repo.OrderGamesByRating
	default:
		return model.GamesFilter{}, errors.New("invalid orderBy: should be one of [default, releaseDate, name]")
	}
	if len(p.Name) >= minLengthForSearch {
		filter.Name = p.Name
	}
	if p.GenreID != 0 {
		filter.GenreID = p.GenreID
	}
	if p.DeveloperID != 0 {
		filter.DeveloperID = p.DeveloperID
	}
	if p.PublisherID != 0 {
		filter.PublisherID = p.PublisherID
	}

	return filter, nil
}
