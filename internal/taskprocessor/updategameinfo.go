package taskprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/pkg/slice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

const (
	// UpdateGameInfoTaskName task name for updating trending index
	UpdateGameInfoTaskName = "update_game_info"

	updateGameInfoBatchSize = 200
)

type updateGameInfoSettings struct {
	LastProcessedID int32 `json:"lastProcessedId"`
}

func (u updateGameInfoSettings) convertToTaskSettings() model.TaskSettings {
	b, _ := json.Marshal(u)
	return b
}

var (
	updateGameInfoProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "update_game_info_processed_total",
		Help: "Total number of games processed for info update",
	})

	updateGameInfoUpdatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "update_game_info_updated_total",
		Help: "Total number of games successfully updated info",
	})
)

// StartUpdateGameInfo starts the update game info task
func (tp *TaskProvider) StartUpdateGameInfo() error {
	taskFn := func(ctx context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
		var s updateGameInfoSettings
		if settings != nil {
			err := json.Unmarshal(settings, &s)
			if err != nil {
				return nil, fmt.Errorf("unmarshal settings: %v", err)
			}
		}

		// get games to update
		gameIDs, err := tp.storage.GetGamesIDsAfterID(ctx, s.LastProcessedID, updateGameInfoBatchSize)
		if err != nil {
			return settings, fmt.Errorf("get games in info update task: %v", err)
		}

		if len(gameIDs) == 0 {
			s.LastProcessedID = 0
			return s.convertToTaskSettings(), nil
		}

		// get stored platform
		platforms, err := tp.storage.GetPlatforms(ctx)
		if err != nil {
			return nil, fmt.Errorf("get platforms in update games info task: %v", err)
		}
		igdbIDPlatformMap := make(map[int64]model.Platform)
		for _, p := range platforms {
			igdbIDPlatformMap[p.IGDBID] = p
		}

		var updatedCount int
		for _, gameID := range gameIDs {
			// wait for igdb rate limit
			if err = tp.igdbAPILimiter.Wait(ctx); err != nil {
				return nil, fmt.Errorf("wait for rate limit in update game info task: %w", err)
			}

			// get game
			game, gErr := tp.storage.GetGameByID(ctx, gameID)
			if gErr != nil {
				tp.log.Error("failed to get game", zap.Int32("game_id", gameID), zap.Error(gErr))
				continue
			}

			if game.IGDBID == 0 {
				continue
			}

			updatedInfo, uErr := tp.igdbAPIClient.GetGameInfoForUpdate(ctx, game.IGDBID)
			if uErr != nil {
				tp.log.Error("failed to get game info from igdb", zap.Int32("game_id", gameID), zap.Error(uErr))
				continue
			}

			updatedData, changed := mapGameToUpdateIGDBGameData(game, updatedInfo, igdbIDPlatformMap)
			if !changed {
				s.LastProcessedID = gameID
				continue
			}

			// update game info
			err = tp.storage.UpdateGameIGDBInfo(ctx, gameID, updatedData)
			if err != nil {
				tp.log.Error("failed to update game info", zap.Int32("game_id", gameID), zap.Error(err))
				continue
			}

			updatedCount++
			updateGameInfoUpdatedTotal.Inc()
			s.LastProcessedID = gameID
		}

		updateGameInfoProcessedTotal.Add(float64(len(gameIDs)))

		tp.log.Info("task info",
			zap.String("name", UpdateGameInfoTaskName),
			zap.Int("games_processed", len(gameIDs)),
			zap.Int("games_updated", updatedCount),
			zap.Int32("last_processed_id", s.LastProcessedID))

		return s.convertToTaskSettings(), nil
	}

	return tp.DoTask(UpdateGameInfoTaskName, taskFn)
}

func mapGameToUpdateIGDBGameData(game model.Game, updateInfo igdbapi.GameInfoForUpdate, platformsMap map[int64]model.Platform) (model.UpdateGameIGDBData, bool) {
	var websites []string
	var platformsIDs []int32

	for _, ipID := range updateInfo.Platforms {
		if p, ok := platformsMap[ipID]; ok {
			platformsIDs = append(platformsIDs, p.ID)
		}
	}

	for _, w := range updateInfo.Websites {
		if _, ok := igdbapi.WebsiteTypeNames[w.Type]; ok {
			websites = append(websites, w.URL)
		}
	}

	data := model.UpdateGameIGDBData{
		Name:            updateInfo.Name,
		PlatformsIDs:    platformsIDs,
		Websites:        websites,
		IGDBRating:      updateInfo.TotalRating,
		IGDBRatingCount: updateInfo.TotalRatingCount,
	}

	changed := game.Name != data.Name ||
		math.Abs(game.IGDBRating-data.IGDBRating) >= 0.1 ||
		game.IGDBRatingCount != data.IGDBRatingCount ||
		!slice.SameValues(game.PlatformsIDs, data.PlatformsIDs) ||
		!slice.SameValues(game.Websites, data.Websites)

	return data, changed
}
