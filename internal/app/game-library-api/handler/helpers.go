package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
)

// Cache keys
const (
	gamesKey        = "games"
	gameKey         = "game"
	gamesCountKey   = "games-count"
	userRatingsKey  = "user-ratings"
	topCompaniesKey = "top-companies"
	topGenresKey    = "top-genres"
)

func getGamesKey(pageSize, page int64, filter repo.GamesFilter) string {
	return gamesKey + "|" + strconv.FormatInt(pageSize, 10) + "|" + strconv.FormatInt(page, 10) + "|" +
		filter.OrderBy.Field + "|" + filter.Name + "|" + strconv.FormatInt(int64(filter.GenreID), 10) + "|" +
		strconv.FormatInt(int64(filter.DeveloperID), 10) + "|" + strconv.FormatInt(int64(filter.PublisherID), 10)
}

func getGameKey(id int32) string {
	return gameKey + "|" + strconv.FormatInt(int64(id), 10)
}

func getGamesCountKey(filter repo.GamesFilter) string {
	return gamesCountKey + "|" + filter.Name + "|" + strconv.FormatInt(int64(filter.GenreID), 10) + "|" +
		strconv.FormatInt(int64(filter.DeveloperID), 10) + "|" + strconv.FormatInt(int64(filter.PublisherID), 10)
}

func getUserRatingsKey(userID string) string {
	return userRatingsKey + "|" + userID
}

func getTopCompaniesKey(companyType string, limit int64) string {
	return topCompaniesKey + "|" + companyType + "|" + strconv.FormatInt(limit, 10)
}

func getTopGenresKey(limit int64) string {
	return topGenresKey + "|" + strconv.FormatInt(limit, 10)
}

// Cached entities functions

func (p *Provider) getCompanyByID(ctx context.Context, id int32) (Company, bool, error) {
	if companiesMap.Size() == 0 {
		companies, err := p.storage.GetCompanies(ctx)
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

func (p *Provider) getGenreByID(ctx context.Context, id int32) (Genre, bool, error) {
	if genresMap.Size() == 0 {
		genres, err := p.storage.GetGenres(ctx)
		if err != nil {
			return Genre{}, false, err
		}
		for _, genre := range genres {
			genresMap.Set(genre.ID, Genre{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}
	}
	genre, ok := genresMap.Get(id)
	return genre, ok, nil
}

func (p *Provider) getPlatformByID(ctx context.Context, id int32) (Platform, bool, error) {
	if platformsMap.Size() == 0 {
		platforms, err := p.storage.GetPlatforms(ctx)
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

// Mappings

func mapToCreateGame(cgr *CreateGameRequest, developerID, publisherID int32) repo.CreateGame {
	return repo.CreateGame{
		Name:        cgr.Name,
		ReleaseDate: cgr.ReleaseDate,
		Developers:  []int32{developerID},
		Publishers:  []int32{publisherID},
		Genres:      cgr.GenresIDs,
		LogoURL:     cgr.LogoURL,
		Summary:     cgr.Summary,
		Slug:        getGameSlug(cgr.Name),
		Platforms:   cgr.PlatformsIDs,
		Screenshots: cgr.Screenshots,
		Websites:    cgr.Websites,
	}
}

func mapToUpdateGame(g repo.Game, ugr UpdateGameRequest) repo.UpdateGame {
	update := repo.UpdateGame{
		Name:        g.Name,
		Developers:  g.Developers,
		Publishers:  g.Publishers,
		ReleaseDate: g.ReleaseDate.String(),
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
		update.Slug = getGameSlug(*ugr.Name)
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

func (p *Provider) mapToGameResponse(ctx context.Context, game repo.Game) (GameResponse, error) {
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
	for _, gID := range game.Genres {
		if genre, ok, err := p.getGenreByID(ctx, gID); err != nil {
			return GameResponse{}, fmt.Errorf("get genre %d by id: %v", gID, err)
		} else if ok {
			resp.Genres = append(resp.Genres, genre)
		}
	}
	for _, pID := range game.Platforms {
		if p, ok, err := p.getPlatformByID(ctx, pID); err != nil {
			return GameResponse{}, fmt.Errorf("get platform %d by id: %v", pID, err)
		} else if ok {
			resp.Platforms = append(resp.Platforms, p)
		}
	}
	for _, cID := range game.Developers {
		if c, ok, err := p.getCompanyByID(ctx, cID); err != nil {
			return GameResponse{}, fmt.Errorf("get developer %d by id: %v", cID, err)
		} else if ok {
			resp.Developers = append(resp.Developers, c)
		}
	}
	for _, cID := range game.Publishers {
		if c, ok, err := p.getCompanyByID(ctx, cID); err != nil {
			return GameResponse{}, fmt.Errorf("get publisher %d by id: %v", cID, err)
		} else if ok {
			resp.Publishers = append(resp.Publishers, c)
		}
	}

	return resp, nil
}

func getGameSlug(name string) string {
	return strings.ReplaceAll(strings.ToLower(strings.ToValidUTF8(name, "")), " ", "-")
}
