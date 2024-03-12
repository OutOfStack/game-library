package taskprocessor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"go.uber.org/zap"
)

const (
	// FetchIGDBGamesTaskName ...
	FetchIGDBGamesTaskName = "fetch_igdb_games"

	fetchGamesMinRating        = 60
	fetchGamesScreenshotsLimit = 7
	fetchGamesRequestsCount    = 5
)

type fetchGamesSettings struct {
	LastReleasedAt time.Time `json:"lastReleasedAt"`
}

func (f fetchGamesSettings) convertToTaskSettings() repo.TaskSettings {
	b, _ := json.Marshal(f)
	return b
}

var (
	startReleasedAtDate = time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC)

	getMinRatingsCountAndLimit = func(date time.Time) (count, limit int64) {
		// [start date, 2001.1.1)
		if date.Year() < 2001 {
			return 200, 10
		}
		// [2001.1.1, 2011.1.1)
		if date.Year() < 2011 {
			return 120, 15
		}
		// [2011.1.1, 2021.1.1)
		if date.Year() < 2021 {
			return 80, 15
		}
		// [2021.1.1, now-6 months)
		if date.Before(time.Now().AddDate(0, -6, 0)) {
			return 50, 10
		}
		// [now-6 months, now-1 month)
		if date.Before(time.Now().AddDate(0, -1, 0)) {
			return 30, 5
		}
		// [now-1 month, now]
		return 15, 3
	}
)

// StartFetchIGDBGames starts fetch igdb games task
func (tp *TaskProvider) StartFetchIGDBGames() error {
	task := func(ctx context.Context, settings repo.TaskSettings) (repo.TaskSettings, error) {
		var s fetchGamesSettings
		if settings != nil {
			err := json.Unmarshal(settings, &s)
			if err != nil {
				return nil, fmt.Errorf("unmarshal settings: %v", err)
			}
		}
		if s.LastReleasedAt.IsZero() {
			s.LastReleasedAt = time.Now()
		}

		// get stored platform
		platforms, err := tp.storage.GetPlatforms(ctx)
		if err != nil {
			return nil, fmt.Errorf("get platforms: %v", err)
		}
		var allPlatformsIDs []int64
		igdbPlatforms := make(map[int64]repo.Platform)
		for _, p := range platforms {
			allPlatformsIDs = append(allPlatformsIDs, p.IGDBID)
			igdbPlatforms[p.IGDBID] = p
		}

		// get stored companies
		companies, err := tp.storage.GetCompanies(ctx)
		if err != nil {
			return nil, fmt.Errorf("get companies: %v", err)
		}
		igdbCompanies := make(map[int64]repo.Company)
		for _, c := range companies {
			if c.IGDBID.Valid {
				igdbCompanies[c.IGDBID.Int64] = c
			}
		}

		// get stored genres
		genres, err := tp.storage.GetGenres(ctx)
		if err != nil {
			return nil, fmt.Errorf("get genres: %v", err)
		}
		igdbGenres := make(map[int64]repo.Genre)
		for _, g := range genres {
			igdbGenres[g.IGDBID] = g
		}

		var gamesAdded int

		for i := 0; i < fetchGamesRequestsCount; i++ {
			ratingsCount, limit := getMinRatingsCountAndLimit(s.LastReleasedAt)
			igdbGames, gErr := tp.igdbProvider.GetTopRatedGames(ctx, allPlatformsIDs, s.LastReleasedAt, ratingsCount, fetchGamesMinRating, limit)
			if gErr != nil {
				return settings, fmt.Errorf("get games from igdb: %v", gErr)
			}

			for _, g := range igdbGames {
				_, err = tp.storage.GetGameIDByIGDBID(ctx, g.ID)
				if err == nil {
					s.LastReleasedAt = time.Unix(g.FirstReleaseDate, 0)
					settings = s.convertToTaskSettings()
					continue
				} else if !errors.As(err, &repo.ErrNotFound[int64]{}) {
					return settings, fmt.Errorf("get game id by igdb id: %v", err)
				}

				// get developers, publishers ids
				var developersIDs, publishersIDs []int32
				for _, ic := range g.InvolvedCompanies {
					if ic.Developer {
						if c, ok := igdbCompanies[ic.Company.ID]; ok {
							developersIDs = append(developersIDs, c.ID)
						} else {
							c.Name = ic.Company.Name
							c.IGDBID = sql.NullInt64{Int64: ic.Company.ID, Valid: true}
							id, cErr := tp.storage.CreateCompany(ctx, c)
							if cErr != nil {
								return settings, fmt.Errorf("create company %v: %v", c, cErr)
							}
							c.ID = id
							developersIDs = append(developersIDs, id)
							igdbCompanies[c.IGDBID.Int64] = c
						}
					}
					if ic.Publisher {
						if c, ok := igdbCompanies[ic.Company.ID]; ok {
							publishersIDs = append(publishersIDs, c.ID)
						} else {
							c.Name = ic.Company.Name
							c.IGDBID = sql.NullInt64{Int64: ic.Company.ID, Valid: true}
							id, cErr := tp.storage.CreateCompany(ctx, c)
							if cErr != nil {
								return settings, fmt.Errorf("create company %v: %v", c, cErr)
							}
							c.ID = id
							publishersIDs = append(publishersIDs, id)
							igdbCompanies[c.IGDBID.Int64] = c
						}
					}
				}

				// get genres ids
				var genresIDs []int32
				for _, ig := range g.Genres {
					if genre, ok := igdbGenres[ig.ID]; ok {
						genresIDs = append(genresIDs, genre.ID)
					} else {
						genre.Name = ig.Name
						genre.IGDBID = ig.ID
						id, cErr := tp.storage.CreateGenre(ctx, genre)
						if cErr != nil {
							return settings, fmt.Errorf("create genre %v: %v", genre, cErr)
						}
						genre.ID = id
						genresIDs = append(genresIDs, id)
						igdbGenres[genre.IGDBID] = genre
					}
				}

				// get platforms ids
				var platformsIDs []int32
				for _, ipID := range g.Platforms {
					if p, ok := igdbPlatforms[ipID]; ok {
						platformsIDs = append(platformsIDs, p.ID)
					}
				}

				// get websites
				var websites []string
				for _, w := range g.Websites {
					if _, ok := igdb.WebsiteCategoryNames[igdb.WebsiteCategory(w.Category)]; ok {
						websites = append(websites, w.URL)
					}
				}

				// get logo url
				igdbLogoURL := igdb.GetImageURL(g.Cover.URL, igdb.ImageCoverBig2xAlias)
				logoURL, uErr := tp.uploadcareProvider.UploadImageFromURL(ctx, igdbLogoURL)
				if uErr != nil {
					return settings, fmt.Errorf("upload logo %s: %v", igdbLogoURL, uErr)
				}

				// get screenshots
				var screenshots []string
				for j, scr := range g.Screenshots {
					if j == fetchGamesScreenshotsLimit {
						break
					}
					igdbScreenshotURL := igdb.GetImageURL(scr.URL, igdb.ImageScreenshotBigAlias)
					screenshotURL, uErr := tp.uploadcareProvider.UploadImageFromURL(ctx, igdbScreenshotURL)
					if uErr != nil {
						return settings, fmt.Errorf("upload screenshot %s: %v", igdbScreenshotURL, uErr)
					}
					screenshots = append(screenshots, screenshotURL)
				}

				cg := repo.CreateGame{
					Name:        g.Name,
					Developers:  developersIDs,
					Publishers:  publishersIDs,
					ReleaseDate: time.Unix(g.FirstReleaseDate, 0).Format("2006-01-02"),
					Genres:      genresIDs,
					LogoURL:     logoURL,
					Summary:     g.Summary,
					Slug:        g.Slug,
					Platforms:   platformsIDs,
					Screenshots: screenshots,
					Websites:    websites,
					IGDBRating:  g.TotalRating,
					IGDBID:      g.ID,
				}

				_, cErr := tp.storage.CreateGame(ctx, cg)
				if cErr != nil {
					return settings, fmt.Errorf("create game %s with igdb id %d: %v", cg.Name, cg.IGDBID, cErr)
				}

				gamesAdded++
				s.LastReleasedAt = time.Unix(g.FirstReleaseDate, 0)
				settings = s.convertToTaskSettings()
			}

			if s.LastReleasedAt.Before(startReleasedAtDate) {
				s.LastReleasedAt = time.Now()
			}
		}

		tp.log.Info("task info", zap.String("name", FetchIGDBGamesTaskName), zap.Int("games_added", gamesAdded), zap.Time("last_released_at", s.LastReleasedAt))

		return s.convertToTaskSettings(), nil
	}

	return tp.DoTask(FetchIGDBGamesTaskName, task)
}
