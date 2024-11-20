package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

// EntityIDType generics type for entity id
type EntityIDType interface {
	int32 | int64 | string
}

// StatusCode - error status code
type StatusCode int

// supported status codes
const (
	NotFound  StatusCode = http.StatusNotFound
	Invalid   StatusCode = http.StatusBadRequest
	Forbidden StatusCode = http.StatusForbidden
)

// Error - custom error wrapper
type Error[T EntityIDType] struct {
	Entity   string
	ID       T
	StatCode StatusCode
	Msg      string
}

// AppError - app error interface
type AppError interface {
	error
	StatusCode() StatusCode
	HTTPStatusCode() int
}

// NewNotFoundError - return new custom not found error
func NewNotFoundError[T EntityIDType](entity string, id T) Error[T] {
	return Error[T]{
		Entity:   entity,
		ID:       id,
		StatCode: NotFound,
	}
}

// NewInvalidError - return new custom invalid error
func NewInvalidError[T EntityIDType](entity string, id T, msg string) Error[T] {
	return Error[T]{
		Entity:   entity,
		ID:       id,
		StatCode: Invalid,
		Msg:      msg,
	}
}

// NewForbiddenError - return new custom forbidden error
func NewForbiddenError[T EntityIDType](entity string, id T) Error[T] {
	return Error[T]{
		Entity:   entity,
		ID:       id,
		StatCode: Forbidden,
	}
}

// StatusCode returns status code
func (e Error[T]) StatusCode() StatusCode {
	return e.StatCode
}

// HTTPStatusCode returns status code
func (e Error[T]) HTTPStatusCode() int {
	return int(e.StatCode)
}

// Error returns error description
func (e Error[T]) Error() string {
	switch e.StatCode {
	case NotFound:
		return fmt.Sprintf("%s with id %v not found", e.Entity, e.ID)
	case Invalid:
		msg := fmt.Sprintf("invalid %s with id %v", e.Entity, e.ID)
		if e.Msg != "" {
			msg += ": " + e.Msg
		}
		return msg
	case Forbidden:
		return fmt.Sprintf("forbidden access to %s with id %v", e.Entity, e.ID)
	}

	return fmt.Sprintf("error %s with id %v: %s", e.Entity, e.ID, e.Msg)
}

// IsAppError checks if error is AppError
func IsAppError(err error) (AppError, bool) {
	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// IsStatusCode checks if error is AppError and status code is statusCode
func IsStatusCode(err error, statusCode StatusCode) bool {
	appErr, ok := IsAppError(err)
	if !ok {
		return false
	}
	return appErr.StatusCode() == statusCode
}
