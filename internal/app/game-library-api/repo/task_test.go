package repo_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

func TestGetTask_TaskExists_ShouldReturnTask(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	taskName := td.String()
	lastRun := time.Now().UTC().Truncate(time.Second)
	settings := `{"foo":"bar"}`

	_, err := db.Exec(ctx, `INSERT INTO background_tasks (name, status, run_count, last_run, settings) VALUES ($1, $2, $3, $4, $5)`,
		taskName,
		model.RunningTaskStatus,
		int64(42),
		lastRun,
		settings,
	)
	require.NoError(t, err)

	task, err := s.GetTask(ctx, taskName)
	require.NoError(t, err)

	require.Equal(t, taskName, task.Name)
	require.Equal(t, model.RunningTaskStatus, task.Status)
	require.Equal(t, int64(42), task.RunCount)
	require.True(t, task.LastRun.Valid)
	require.WithinDuration(t, lastRun, task.LastRun.Time.UTC(), time.Second)

	var gotSettings map[string]any
	require.NoError(t, json.Unmarshal(task.Settings, &gotSettings))
	require.Equal(t, map[string]any{"foo": "bar"}, gotSettings)
}

func TestGetTask_TaskNotFound_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	taskName := td.String()

	task, err := s.GetTask(ctx, taskName)
	require.ErrorIs(t, err, apperr.NewNotFoundError("task", taskName))
	require.Equal(t, model.Task{}, task)
}

func TestGetTask_TaskLocked_ShouldReturnErrTransactionLocked(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	taskName := td.String()

	_, err := db.Exec(ctx, `INSERT INTO background_tasks (name, status, run_count, last_run, settings) VALUES ($1, $2, $3, $4, $5)`,
		taskName,
		model.IdleTaskStatus,
		int64(0),
		nil,
		"{}",
	)
	require.NoError(t, err)

	conn, err := db.Acquire(ctx)
	require.NoError(t, err)
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var lockedName string
	err = tx.QueryRow(ctx, `SELECT name FROM background_tasks WHERE name = $1 FOR UPDATE`, taskName).Scan(&lockedName)
	require.NoError(t, err)

	task, err := s.GetTask(ctx, taskName)
	require.ErrorIs(t, err, repo.ErrTransactionLocked)
	require.Equal(t, model.Task{}, task)
}

func TestUpdateTask_TaskExists_ShouldUpdateRow(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	taskName := td.String()
	initialLastRun := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)

	_, err := db.Exec(ctx,
		`INSERT INTO background_tasks (name, status, run_count, last_run, settings) 
		VALUES ($1, $2, $3, $4, $5)`,
		taskName, model.IdleTaskStatus, int64(1), initialLastRun, `{"enabled":false}`)
	require.NoError(t, err)

	updatedLastRun := time.Now().UTC().Add(time.Minute).Truncate(time.Second)
	update := model.Task{
		Name:     taskName,
		Status:   model.ErrorTaskStatus,
		RunCount: 5,
		LastRun:  sql.NullTime{Time: updatedLastRun, Valid: true},
		Settings: model.TaskSettings(`{"enabled":true}`),
	}

	err = s.UpdateTask(ctx, update)
	require.NoError(t, err)

	var (
		status   string
		runCount int64
		lastRun  sql.NullTime
		settings []byte
	)
	err = db.QueryRow(ctx, `SELECT status, run_count, last_run, settings FROM background_tasks WHERE name = $1`, taskName).
		Scan(&status, &runCount, &lastRun, &settings)
	require.NoError(t, err)

	require.Equal(t, string(update.Status), status)
	require.Equal(t, update.RunCount, runCount)
	require.True(t, lastRun.Valid)
	require.WithinDuration(t, update.LastRun.Time, lastRun.Time.UTC(), time.Second)

	var gotSettings map[string]any
	require.NoError(t, json.Unmarshal(settings, &gotSettings))
	require.Equal(t, map[string]any{"enabled": true}, gotSettings)
}

func TestUpdateTask_TaskNotFound_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	taskName := td.String()

	err := s.UpdateTask(ctx, model.Task{
		Name:     taskName,
		Status:   model.ErrorTaskStatus,
		RunCount: 1,
		LastRun:  sql.NullTime{Time: time.Now().UTC(), Valid: true},
		Settings: []byte(`{"enabled":false}`),
	})
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", taskName))
}
