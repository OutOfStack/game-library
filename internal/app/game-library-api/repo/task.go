package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// GetTask returns task status.
// If tx provided, query will be executed on it.
// If task does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetTask(ctx context.Context, tx pgx.Tx, name string) (task model.Task, err error) {
	ctx, span := tracer.Start(ctx, "db.getTask")
	defer span.End()

	q := `
		SELECT name, status, run_count, last_run, settings
		FROM background_tasks
		WHERE name = $1
		FOR NO KEY UPDATE NOWAIT`

	if tx != nil {
		err = pgxscan.Get(ctx, tx, &task, q, name)
	} else {
		err = pgxscan.Get(ctx, s.db, &task, q, name)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, apperr.NewNotFoundError("task", name)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == codeLockNotAvailable {
			return model.Task{}, ErrTransactionLocked
		}
		return model.Task{}, err
	}

	return task, nil
}

// UpdateTask updates task.
// If tx provided, query will be executed on it.
// If task does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateTask(ctx context.Context, tx pgx.Tx, task model.Task) (err error) {
	ctx, span := tracer.Start(ctx, "db.updateTask")
	defer span.End()

	q := `
		UPDATE background_tasks
    	SET status = $2, last_run = $3, run_count = $4, settings = coalesce($5, settings), updated_at = $6
		WHERE name=$1`

	var res pgconn.CommandTag
	if tx != nil {
		res, err = tx.Exec(ctx, q, task.Name, string(task.Status), task.LastRun, task.RunCount, task.Settings, time.Now())
	} else {
		res, err = s.db.Exec(ctx, q, task.Name, string(task.Status), task.LastRun, task.RunCount, task.Settings, time.Now())
	}
	if err != nil {
		return fmt.Errorf("updating task %s: %v", task.Name, err)
	}

	return checkRowsAffected(res, "game", task.Name)
}
