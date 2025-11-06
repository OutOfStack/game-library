package repo

import (
	"errors"

	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/jackc/pgx/v5/pgconn"
)

// pg error codes
const (
	// Object Not In Prerequisite State
	codeLockNotAvailable = "55P03"
)

var (
	// ErrTransactionLocked - error representing transaction lock
	ErrTransactionLocked = errors.New("transaction locked")
)

func checkRowsAffected[T apperr.EntityIDType](res pgconn.CommandTag, entity string, id T) error {
	count := res.RowsAffected()
	if count == 0 {
		return apperr.NewNotFoundError(entity, id)
	}
	return nil
}
