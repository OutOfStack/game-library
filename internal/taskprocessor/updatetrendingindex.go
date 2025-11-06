package taskprocessor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

const (
	// UpdateTrendingIndexTaskName task name for updating trending index
	UpdateTrendingIndexTaskName = "update_trending_index"

	updateTrendingIndexBatchSize = 300
)

type updateTrendingIndexSettings struct {
	LastProcessedID int32 `json:"lastProcessedId"`
}

func (u updateTrendingIndexSettings) convertToTaskSettings() model.TaskSettings {
	b, _ := json.Marshal(u)
	return b
}

var (
	updateTrendingIndexProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "update_trending_index_processed_total",
		Help: "Total number of games processed for trending index updates",
	})

	updateTrendingIndexUpdatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "update_trending_index_updated_total",
		Help: "Total number of games successfully updated with new trending index",
	})
)

// StartUpdateTrendingIndex starts the update trending index task
func (tp *TaskProvider) StartUpdateTrendingIndex() error {
	taskFn := func(ctx context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
		var s updateTrendingIndexSettings
		if settings != nil {
			err := json.Unmarshal(settings, &s)
			if err != nil {
				return nil, fmt.Errorf("unmarshal settings: %v", err)
			}
		}

		// get games to update
		gameIDs, err := tp.storage.GetGamesIDsAfterID(ctx, s.LastProcessedID, updateTrendingIndexBatchSize)
		if err != nil {
			return settings, fmt.Errorf("get games for trending index update: %v", err)
		}

		if len(gameIDs) == 0 {
			s.LastProcessedID = 0
			return s.convertToTaskSettings(), nil
		}

		var updatedCount int
		for _, gameID := range gameIDs {
			// update trending index
			err = tp.gameFacade.UpdateGameTrendingIndex(ctx, gameID)
			if err != nil {
				tp.log.Error("failed to update trending index", zap.Int32("game_id", gameID), zap.Error(err))
				continue
			}

			updatedCount++
			updateTrendingIndexUpdatedTotal.Inc()
			s.LastProcessedID = gameID
		}

		updateTrendingIndexProcessedTotal.Add(float64(len(gameIDs)))

		tp.log.Info("task info",
			zap.String("name", UpdateTrendingIndexTaskName),
			zap.Int("games_processed", len(gameIDs)),
			zap.Int("games_updated", updatedCount),
			zap.Int32("last_processed_id", s.LastProcessedID))

		return s.convertToTaskSettings(), nil
	}

	return tp.DoTask(UpdateTrendingIndexTaskName, taskFn)
}
