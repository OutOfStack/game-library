package taskprocessor

import (
	"go.uber.org/zap"
)

const (
	// FetchIGDBGamesTaskName ...
	FetchIGDBGamesTaskName = "fetch_igdb_games"
)

// StartFetchIGDBGames starts fetch igdb games task
func (tp *TaskProvider) StartFetchIGDBGames() error {
	task := func() error {
		tp.log.Info("Dummy task completed!", zap.String("task", FetchIGDBGamesTaskName))
		return nil
	}

	// TODO uncomment when task implemented
	// return tp.DoTask(FetchIGDBGamesTaskName, task)
	return task()
}
