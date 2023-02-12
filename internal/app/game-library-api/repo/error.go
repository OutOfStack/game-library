package repo

import (
	"database/sql"
	"errors"
	"fmt"
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

// EntityID generics type for entity id
type EntityID interface {
	int32 | int64 | string
}

// ErrNotFound is used when a requested entity with id does not exist
type ErrNotFound[T EntityID] struct {
	Entity string
	ID     T
}

func (e ErrNotFound[T]) Error() string {
	return fmt.Sprintf("%v with id %v was not found", e.Entity, e.ID)
}

func checkRowsAffected[T EntityID](res sql.Result, entity string, id T) error {
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound[T]{Entity: entity, ID: id}
	}
	return nil
}
