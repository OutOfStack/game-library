package taskprocessor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func (f fetchGamesSettings) convertToTaskSettings() model.TaskSettings {
	b, _ := json.Marshal(f)
	return b
}

var (
	startReleasedAtDate = time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC)

	getMinRatingsCountAndLimit = func(date time.Time) (count, limit int64) {
		switch {
		case date.Year() < 2001: // [start date, 2001.1.1)
			return 200, 10
		case date.Year() < 2011: // [2001.1.1, 2011.1.1)
			return 120, 15
		case date.Year() < 2021: // [2011.1.1, 2021.1.1)
			return 80, 15
		case date.Before(time.Now().AddDate(0, -6, 0)): // [2021.1.1, now-6 months)
			return 50, 10
		case date.Before(time.Now().AddDate(0, -1, 0)): // [now-6 months, now-1 month)
			return 30, 5
		default: // [now-1 month, now]
			return 15, 3
		}
	}
)

var (
	fetchGamesAddedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fetch_igdb_games_added_total",
		Help: "Total number of games successfully added by the fetch IGDB games task",
	})

	fetchGamesSkippedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fetch_igdb_games_skipped_total",
		Help: "Total number of games skipped because they already exist",
	})
)

// StartFetchIGDBGames starts fetch igdb games task
func (tp *TaskProvider) StartFetchIGDBGames() error {
	taskFn := func(ctx context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
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
		igdbPlatforms := make(map[int64]model.Platform)
		for _, p := range platforms {
			allPlatformsIDs = append(allPlatformsIDs, p.IGDBID)
			igdbPlatforms[p.IGDBID] = p
		}

		// get stored companies
		companies, err := tp.storage.GetCompanies(ctx)
		if err != nil {
			return nil, fmt.Errorf("get companies: %v", err)
		}
		igdbCompanies := make(map[int64]model.Company)
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
		igdbGenres := make(map[int64]model.Genre)
		for _, g := range genres {
			igdbGenres[g.IGDBID] = g
		}

		var gamesAdded int

		for range fetchGamesRequestsCount {
			ratingsCount, limit := getMinRatingsCountAndLimit(s.LastReleasedAt)
			igdbGames, gErr := tp.igdbAPIClient.GetTopRatedGames(ctx, allPlatformsIDs, s.LastReleasedAt, ratingsCount, fetchGamesMinRating, limit)
			if gErr != nil {
				return settings, fmt.Errorf("get games from igdb: %v", gErr)
			}

			for _, g := range igdbGames {
				_, err = tp.storage.GetGameIDByIGDBID(ctx, g.ID)
				if err == nil {
					fetchGamesSkippedTotal.Inc()
					s.LastReleasedAt = time.Unix(g.FirstReleaseDate, 0)
					settings = s.convertToTaskSettings()
					continue
				} else if !apperr.IsStatusCode(err, apperr.NotFound) {
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
					if _, ok := igdbapi.WebsiteCategoryNames[igdbapi.WebsiteCategory(w.Category)]; ok {
						websites = append(websites, w.URL)
					}
				}

				// reupload logo url
				logoData, lErr := tp.igdbAPIClient.GetImageByURL(ctx, g.Cover.URL, igdbapi.ImageTypeCoverBig2xAlias)
				if lErr != nil {
					return settings, fmt.Errorf("get logo by url %s: %v", g.Cover.URL, lErr)
				}

				logoUploadData, uErr := tp.s3Client.Upload(ctx, logoData.Body, logoData.FileName, logoData.ContentType)
				if uErr != nil {
					return settings, fmt.Errorf("upload logo %s: %v", g.Cover.URL, uErr)
				}

				// reupload screenshots
				var screenshots []string
				for j, scr := range g.Screenshots {
					if j == fetchGamesScreenshotsLimit {
						break
					}
					scrData, sErr := tp.igdbAPIClient.GetImageByURL(ctx, scr.URL, igdbapi.ImageTypeScreenshotBigAlias)
					if sErr != nil {
						return settings, fmt.Errorf("get screenshot by url %s: %v", scr.URL, sErr)
					}
					screnshotUploadData, sErr := tp.s3Client.Upload(ctx, scrData.Body, scrData.FileName, scrData.ContentType)
					if sErr != nil {
						return settings, fmt.Errorf("upload screenshot %s: %v", scr.URL, sErr)
					}
					screenshots = append(screenshots, screnshotUploadData.FileURL)
				}

				cg := model.CreateGame{
					Name:          g.Name,
					DevelopersIDs: developersIDs,
					PublishersIDs: publishersIDs,
					ReleaseDate:   time.Unix(g.FirstReleaseDate, 0).Format("2006-01-02"),
					GenresIDs:     genresIDs,
					LogoURL:       logoUploadData.FileURL,
					Summary:       g.Summary,
					Slug:          g.Slug,
					PlatformsIDs:  platformsIDs,
					Screenshots:   screenshots,
					Websites:      websites,
					IGDBRating:    g.TotalRating,
					IGDBID:        g.ID,
				}

				_, cErr := tp.storage.CreateGame(ctx, cg)
				if cErr != nil {
					return settings, fmt.Errorf("create game %s with igdb id %d: %v", cg.Name, cg.IGDBID, cErr)
				}

				fetchGamesAddedTotal.Inc()
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

	return tp.DoTask(FetchIGDBGamesTaskName, taskFn)
}
