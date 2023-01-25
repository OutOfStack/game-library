package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// GetTask returns task status.
// If tx provided, query will be executed on it.
// If task does not exist returns ErrNotFound
func (s *Storage) GetTask(ctx context.Context, tx *sqlx.Tx, name string) (task Task, err error) {
	ctx, span := tracer.Start(ctx, "db.task.get")
	defer span.End()

	q := `SELECT name, status, run_count, last_run
	FROM background_tasks
	WHERE name=$1
	FOR NO KEY UPDATE NOWAIT`

	if tx != nil {
		err = tx.GetContext(ctx, &task, q, name)
	} else {
		err = s.db.GetContext(ctx, &task, q, name)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNotFound[string]{Entity: "task", ID: name}
		}
		pqErr := err.(*pq.Error)
		if pqErr.Code == codeLockNotAvailable {
			return Task{}, ErrTransactionLocked
		}
		return Task{}, err
	}

	return task, nil
}

// UpdateTask updates task.
// If tx provided, query will be executed on it.
// If task does not exist returns ErrNotFound
func (s *Storage) UpdateTask(ctx context.Context, tx *sqlx.Tx, task Task) error {
	ctx, span := tracer.Start(ctx, "db.task.update")
	defer span.End()

	s.db.Beginx()

	q := `UPDATE background_tasks
    SET status = $2, last_run = $3, run_count = $4, updated_at = now()
	WHERE name=$1`

	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.ExecContext(ctx, q, task.Name, string(task.Status), task.LastRun, task.RunCount)
	} else {
		res, err = s.db.ExecContext(ctx, q, task.Name, string(task.Status), task.LastRun, task.RunCount)
	}
	if err != nil {
		return fmt.Errorf("updating task %s: %v", task.Name, err)
	}

	return checkRowsAffected(res, "game", task.Name)
}
