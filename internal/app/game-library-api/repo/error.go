package repo

import (
	"database/sql"
	"errors"

	"github.com/OutOfStack/game-library/internal/pkg/apperr"
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

func checkRowsAffected[T apperr.EntityIDType](res sql.Result, entity string, id T) error {
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return apperr.NewNotFoundError(entity, id)
	}
	return nil
}
