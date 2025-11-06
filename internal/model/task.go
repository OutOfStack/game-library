package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// TaskStatus represents task status type
type TaskStatus string

// Task status values
const (
	IdleTaskStatus    TaskStatus = "idle"
	RunningTaskStatus TaskStatus = "running"
	ErrorTaskStatus   TaskStatus = "error"
)

// Task represents task entity
type Task struct {
	Name     string       `db:"name"`
	Status   TaskStatus   `db:"status"`
	RunCount int64        `db:"run_count"`
	LastRun  sql.NullTime `db:"last_run"`
	Settings TaskSettings `db:"settings"`
}

// TaskSettings task settings value
type TaskSettings []byte

// Value implements driver.Valuer interface
func (ts TaskSettings) Value() (driver.Value, error) {
	return string(ts), nil
}

// Scan implements sql.Scanner interface
func (ts *TaskSettings) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		*ts = []byte(v)
	case []byte:
		*ts = v
	default:
		return fmt.Errorf("scan TaskSettings: unsupported type %T", src)
	}

	return nil
}

// TaskInfo - task info
type TaskInfo struct {
	Schedule string
	Fn       func() error
}
