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
	// ProcessModerationTaskName task name for processing game moderation
	ProcessModerationTaskName = "process_moderation"

	processModerationBatchSize = 10
)

type processModerationSettings struct {
	LastProcessedGameID int32 `json:"lastProcessedGameId"`
}

func (p processModerationSettings) convertToTaskSettings() model.TaskSettings {
	b, _ := json.Marshal(p)
	return b
}

var (
	processModerationProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "process_moderation_processed_total",
		Help: "Total number of games processed for moderation",
	})

	processModerationErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "process_moderation_errors_total",
		Help: "Total number of moderation processing errors",
	})
)

// StartProcessModeration starts the process moderation task
func (tp *TaskProvider) StartProcessModeration() error {
	taskFn := func(ctx context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
		var s processModerationSettings
		if settings != nil {
			err := json.Unmarshal(settings, &s)
			if err != nil {
				return nil, fmt.Errorf("unmarshal settings: %v", err)
			}
		}

		var modGameRecords []model.ModerationIDGameID
		txErr := tp.storage.RunWithTx(ctx, func(ctx context.Context) error {
			var err error
			// get games and moderation records ids with pending moderation status
			modGameRecords, err = tp.storage.GetPendingModerationGameIDs(ctx, processModerationBatchSize)
			if err != nil {
				return fmt.Errorf("get pending moderation games: %v", err)
			}

			if len(modGameRecords) == 0 {
				return nil
			}

			var moderationIDs []int32
			for _, r := range modGameRecords {
				moderationIDs = append(moderationIDs, r.ModerationID)
			}

			// update status to in_progress to prevent other workers from processing
			err = tp.storage.SetModerationRecordsStatus(ctx, moderationIDs, model.ModerationStatusInProgress)
			if err != nil {
				return fmt.Errorf("update moderation status to in_progress: %v", err)
			}

			return nil
		})
		if txErr != nil {
			return settings, txErr
		}

		if len(modGameRecords) == 0 {
			tp.log.Info("no pending moderation games found")
			return s.convertToTaskSettings(), nil
		}

		tp.log.Info("found games for moderation processing", zap.Int("count", len(modGameRecords)))

		var processedCount, errorCount int
		var failedModerationIDs []int32

		// process each game
		for _, record := range modGameRecords {
			tp.log.Info("processing moderation for game", zap.Int32("game_id", record.GameID))

			err := tp.moderationFacade.ProcessModeration(ctx, record.GameID)
			if err != nil {
				failedModerationIDs = append(failedModerationIDs, record.ModerationID)
				tp.log.Error("failed to process moderation", zap.Int32("game_id", record.GameID), zap.Error(err))
				errorCount++
				processModerationErrorsTotal.Inc()
				continue
			}

			processedCount++
			processModerationProcessedTotal.Inc()

			// update last processed game ID
			if record.GameID > s.LastProcessedGameID {
				s.LastProcessedGameID = record.GameID
			}
		}

		// set status to pending to failed moderation attempts
		err := tp.storage.SetModerationRecordsStatus(ctx, failedModerationIDs, model.ModerationStatusPending)
		if err != nil {
			return nil, fmt.Errorf("update moderation status to pending: %v", err)
		}

		tp.log.Info("moderation processing completed",
			zap.String("task", ProcessModerationTaskName),
			zap.Int("games_found", len(modGameRecords)),
			zap.Int("games_processed", processedCount),
			zap.Int("errors", errorCount),
			zap.Int32("last_processed_id", s.LastProcessedGameID))

		return s.convertToTaskSettings(), nil
	}

	return tp.DoTask(ProcessModerationTaskName, taskFn)
}
