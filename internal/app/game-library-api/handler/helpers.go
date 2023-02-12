package handler

import (
	"context"
	"fmt"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
)

func (g *Game) getCompanyByID(ctx context.Context, id int32) (Company, bool, error) {
	if companiesMap.Size() == 0 {
		companies, err := g.storage.GetCompanies(ctx)
		if err != nil {
			return Company{}, false, err
		}
		for _, c := range companies {
			companiesMap.Set(c.ID, Company{
				ID:   c.ID,
				Name: c.Name,
			})
		}
	}
	company, ok := companiesMap.Get(id)
	return company, ok, nil
}

func (g *Game) getGenreByID(ctx context.Context, id int32) (Genre, bool, error) {
	if genresMap.Size() == 0 {
		genres, err := g.storage.GetGenres(ctx)
		if err != nil {
			return Genre{}, false, err
		}
		for _, g := range genres {
			genresMap.Set(g.ID, Genre{
				ID:   g.ID,
				Name: g.Name,
			})
		}
	}
	genre, ok := genresMap.Get(id)
	return genre, ok, nil
}

func (g *Game) getPlatformByID(ctx context.Context, id int32) (Platform, bool, error) {
	if platformsMap.Size() == 0 {
		platforms, err := g.storage.GetPlatforms(ctx)
		if err != nil {
			return Platform{}, false, err
		}
		for _, p := range platforms {
			platformsMap.Set(p.ID, Platform{
				ID:           p.ID,
				Name:         p.Name,
				Abbreviation: p.Abbreviation,
			})
		}
	}
	platform, ok := platformsMap.Get(id)
	return platform, ok, nil
}

func mapToCreateRating(crr *CreateRatingRequest, gameID int32, userID string) repo.CreateRating {
	return repo.CreateRating{
		Rating: crr.Rating,
		UserID: userID,
		GameID: gameID,
	}
}

func mapToRatingResponse(cr repo.CreateRating) *RatingResponse {
	return &RatingResponse{
		GameID: cr.GameID,
		UserID: cr.UserID,
		Rating: cr.Rating,
	}
}

func mapToCreateGame(cgr *CreateGameRequest) repo.CreateGame {
	return repo.CreateGame{
		Name:        cgr.Name,
		ReleaseDate: cgr.ReleaseDate,
		Developer:   cgr.Developer,
		Genre:       cgr.Genre,
		Genres:      cgr.GenresIDs,
		LogoURL:     cgr.LogoURL,
		Summary:     cgr.Summary,
		Slug:        cgr.Slug,
		Platforms:   cgr.PlatformsIDs,
		Screenshots: cgr.Screenshots,
		Websites:    cgr.Websites,
	}
}

func mapToUpdateGame(g repo.Game, ugr UpdateGameRequest) repo.UpdateGame {
	update := repo.UpdateGame{
		Name:        g.Name,
		Developer:   g.Developer,
		Developers:  g.Developers,
		Publisher:   g.Publisher,
		Publishers:  g.Publishers,
		ReleaseDate: g.ReleaseDate.String(),
		Genre:       g.Genre,
		Genres:      g.Genres,
		LogoURL:     g.LogoURL,
		Summary:     g.Summary,
		Slug:        g.Slug,
		Platforms:   g.Platforms,
		Screenshots: g.Screenshots,
		Websites:    g.Websites,
		IGDBRating:  g.IGDBRating,
		IGDBID:      g.IGDBID,
	}

	if ugr.Name != nil {
		update.Name = *ugr.Name
	}
	if ugr.Developer != nil {
		update.Developer = *ugr.Developer
	}
	if ugr.ReleaseDate != nil {
		update.ReleaseDate = *ugr.ReleaseDate
	}
	if ugr.GenresIDs != nil {
		update.Genres = *ugr.GenresIDs
	}
	if ugr.LogoURL != nil && *ugr.LogoURL != "" {
		update.LogoURL = *ugr.LogoURL
	}
	if ugr.Summary != nil {
		update.Summary = *ugr.Summary
	}
	if ugr.Slug != nil {
		update.Slug = *ugr.Slug
	}
	if ugr.Platforms != nil {
		update.Platforms = *ugr.Platforms
	}
	if ugr.Screenshots != nil {
		update.Screenshots = *ugr.Screenshots
	}
	if ugr.Websites != nil {
		update.Websites = *ugr.Websites
	}

	return update
}

func (g *Game) mapToGameResponse(ctx context.Context, game repo.Game) (GameResponse, error) {
	resp := GameResponse{
		ID:          game.ID,
		Name:        game.Name,
		Developer:   game.Developer,
		Publisher:   game.Publisher,
		ReleaseDate: game.ReleaseDate.String(),
		Genre:       game.Genre,
		LogoURL:     game.LogoURL,
		Rating:      game.Rating,
		Summary:     game.Summary,
		Slug:        game.Slug,
		Screenshots: game.Screenshots,
		Websites:    game.Websites,
	}
	for _, gID := range game.Genres {
		if g, ok, err := g.getGenreByID(ctx, gID); err != nil {
			return GameResponse{}, fmt.Errorf("get genre %d by id: %v", gID, err)
		} else if ok {
			resp.Genres = append(resp.Genres, g)
		}
	}
	for _, pID := range game.Platforms {
		if p, ok, err := g.getPlatformByID(ctx, pID); err != nil {
			return GameResponse{}, fmt.Errorf("get platform %d by id: %v", pID, err)
		} else if ok {
			resp.Platforms = append(resp.Platforms, p)
		}
	}
	for _, cID := range game.Developers {
		if c, ok, err := g.getCompanyByID(ctx, cID); err != nil {
			return GameResponse{}, fmt.Errorf("get developer %d by id: %v", cID, err)
		} else if ok {
			resp.Developers = append(resp.Developers, c)
		}
	}
	for _, cID := range game.Publishers {
		if c, ok, err := g.getCompanyByID(ctx, cID); err != nil {
			return GameResponse{}, fmt.Errorf("get publisher %d by id: %v", cID, err)
		} else if ok {
			resp.Publishers = append(resp.Publishers, c)
		}
	}

	// fill deprecated fields
	if len(resp.Genres) == 0 {
		for _, g := range resp.Genres {
			resp.Genre = append(resp.Genre, g.Name)
		}
	}
	if resp.Developer == "" && len(resp.Developers) > 0 {
		resp.Developer = resp.Developers[0].Name
	}
	if resp.Publisher == "" && len(resp.Publishers) > 0 {
		resp.Publisher = resp.Publishers[0].Name
	}

	return resp, nil
}
