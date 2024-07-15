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
	for _, gID := range game.Genres {
		if genre, err := p.gameFacade.GetGenreByID(ctx, gID); err != nil {
			return api.GameResponse{}, fmt.Errorf("get genre %d by id: %v", gID, err)
		} else {
			resp.Genres = append(resp.Genres, api.Genre{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}
	}
	for _, pID := range game.Platforms {
		if pl, err := p.gameFacade.GetPlatformByID(ctx, pID); err != nil {
			return api.GameResponse{}, fmt.Errorf("get platform %d by id: %v", pID, err)
		} else {
			resp.Platforms = append(resp.Platforms, api.Platform{
				ID:           pl.ID,
				Name:         pl.Name,
				Abbreviation: pl.Abbreviation,
			})
		}
	}
	for _, cID := range game.Developers {
		if c, err := p.gameFacade.GetCompanyByID(ctx, cID); err != nil {
			return api.GameResponse{}, fmt.Errorf("get developer %d by id: %v", cID, err)
		} else {
			resp.Developers = append(resp.Developers, api.Company{
				ID:   c.ID,
				Name: c.Name,
			})
		}
	}
	for _, cID := range game.Publishers {
		if c, err := p.gameFacade.GetCompanyByID(ctx, cID); err != nil {
			return api.GameResponse{}, fmt.Errorf("get publisher %d by id: %v", cID, err)
		} else {
			resp.Publishers = append(resp.Publishers, api.Company{
				ID:   c.ID,
				Name: c.Name,
			})
		}
	}

	return resp, nil
}
