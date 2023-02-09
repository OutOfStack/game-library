package taskprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/client/igdb"
)

const (
	// FetchIGDBGamesTaskName ...
	FetchIGDBGamesTaskName = "fetch_igdb_games"

	fetchGamesMinRating        = 60
	fetchGamesScreenshotsLimit = 8
	fetchGamesLimit            = 10
	fetchGamesRequestCount     = 10
)

type fetchGamesSettings struct {
	LastReleasedAt time.Time `json:"lastReleasedAt"`
}

func (f fetchGamesSettings) convertToTaskSettings() repo.TaskSettings {
	b, _ := json.Marshal(f)
	return b
}

var (
	startDate = time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC)

	getMinRatingCount = func(date time.Time) int64 {
		if date.Year() < 2001 {
			return 200
		}
		if date.Year() < 2011 {
			return 100
		}
		if date.Year() < 2021 {
			return 70
		}
		if date.Year() < time.Now().Year() {
			return 35
		}
		return 25
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
			s.LastReleasedAt = startDate
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
			igdbCompanies[c.IGDBID] = c
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

		resetLastReleasedAt := false
		for i := 0; i < fetchGamesRequestCount; i++ {
			igdbGames, gErr := tp.igdbProvider.GetTopRatedGames(ctx, getMinRatingCount(s.LastReleasedAt), fetchGamesMinRating,
				s.LastReleasedAt, fetchGamesLimit, allPlatformsIDs)
			if gErr != nil {
				return settings, fmt.Errorf("get games from igdb: %v", gErr)
			}

			for _, g := range igdbGames {
				_, err = tp.storage.GetGameIDByIGDBID(ctx, g.ID)
				if err == nil {
					s.LastReleasedAt = time.Unix(g.FirstReleaseDate, 0)
					settings = s.convertToTaskSettings()
					continue
				}
				if !errors.As(err, &repo.ErrNotFound[int32]{}) {
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
							c.IGDBID = ic.Company.ID
							id, cErr := tp.storage.CreateCompany(ctx, c)
							if cErr != nil {
								return settings, fmt.Errorf("create company %v: %v", c, cErr)
							}
							c.ID = id
							developersIDs = append(developersIDs, id)
							igdbCompanies[c.IGDBID] = c
						}
					}
					if ic.Publisher {
						if c, ok := igdbCompanies[ic.Company.ID]; ok {
							publishersIDs = append(publishersIDs, c.ID)
						} else {
							c.Name = ic.Company.Name
							c.IGDBID = ic.Company.ID
							id, cErr := tp.storage.CreateCompany(ctx, c)
							if cErr != nil {
								return settings, fmt.Errorf("create company %v: %v", c, cErr)
							}
							c.ID = id
							publishersIDs = append(publishersIDs, id)
							igdbCompanies[c.IGDBID] = c
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

				s.LastReleasedAt = time.Unix(g.FirstReleaseDate, 0)
				settings = s.convertToTaskSettings()
			}

			if len(igdbGames) < fetchGamesLimit {
				if resetLastReleasedAt {
					s.LastReleasedAt = startDate
					break
				}
				s.LastReleasedAt = time.Date(s.LastReleasedAt.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
				resetLastReleasedAt = true
			}
		}

		return s.convertToTaskSettings(), nil
	}

	return tp.DoTask(FetchIGDBGamesTaskName, task)
}
