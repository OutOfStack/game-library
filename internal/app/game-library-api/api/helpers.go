package api

import (
	"context"
	"fmt"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
)

// Mappings

func mapToCreateGame(cgr *api.CreateGameRequest, developer, publisher string) model.CreateGame {
	return model.CreateGame{
		Name:        cgr.Name,
		ReleaseDate: cgr.ReleaseDate,
		Developer:   developer,
		Publisher:   publisher,
		Genres:      cgr.GenresIDs,
		LogoURL:     cgr.LogoURL,
		Summary:     cgr.Summary,
		Slug:        model.GetGameSlug(cgr.Name),
		Platforms:   cgr.PlatformsIDs,
		Screenshots: cgr.Screenshots,
		Websites:    cgr.Websites,
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

	for _, gID := range game.Genres {
		genre, ok := genres[gID]
		if !ok {
			return api.GameResponse{}, fmt.Errorf("genre %d not found", gID)
		}
		resp.Genres = append(resp.Genres, api.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}
	for _, plID := range game.Platforms {
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
	for _, cID := range game.Developers {
		c, ok := companies[cID]
		if !ok {
			return api.GameResponse{}, fmt.Errorf("company %d not found", cID)
		}
		resp.Developers = append(resp.Developers, api.Company{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	for _, cID := range game.Publishers {
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
