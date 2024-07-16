package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// GetTask returns task status.
// If tx provided, query will be executed on it.
// If task does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetTask(ctx context.Context, tx *sqlx.Tx, name string) (task model.Task, err error) {
	ctx, span := tracer.Start(ctx, "db.getTask")
	defer span.End()

	q := `
		SELECT name, status, run_count, last_run, settings
		FROM background_tasks
		WHERE name = $1
		FOR NO KEY UPDATE NOWAIT`

	if tx != nil {
		err = tx.GetContext(ctx, &task, q, name)
	} else {
		err = s.db.GetContext(ctx, &task, q, name)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, apperr.NewNotFoundError("task", name)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == codeLockNotAvailable {
			return model.Task{}, ErrTransactionLocked
		}
		return model.Task{}, err
	}

	return task, nil
}

// UpdateTask updates task.
// If tx provided, query will be executed on it.
// If task does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateTask(ctx context.Context, tx *sqlx.Tx, task model.Task) (err error) {
	ctx, span := tracer.Start(ctx, "db.updateTask")
	defer span.End()

	q := `
		UPDATE background_tasks
    	SET status = $2, last_run = $3, run_count = $4, settings = coalesce($5, settings), updated_at = $6
		WHERE name=$1`

	var res sql.Result
	if tx != nil {
		res, err = tx.ExecContext(ctx, q, task.Name, string(task.Status), task.LastRun, task.RunCount, task.Settings, time.Now())
	} else {
		res, err = s.db.ExecContext(ctx, q, task.Name, string(task.Status), task.LastRun, task.RunCount, task.Settings, time.Now())
	}
	if err != nil {
		return fmt.Errorf("updating task %s: %v", task.Name, err)
	}

	return checkRowsAffected(res, "game", task.Name)
}
