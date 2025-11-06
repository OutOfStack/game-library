package web

import (
	"errors"
	"fmt"
	"net/http"
)

// FieldError represents error in a struct field
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse represents response in case of a web error
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error contains information about web error
type Error struct {
	Err        error
	StatusCode int
	Fields     []FieldError
}

// NewError creates error with status code
func NewError(err error, statusCode int) error {
	return &Error{
		Err:        err,
		StatusCode: statusCode,
	}
}

// NewErrorFromMessage creates error with status code from message
func NewErrorFromMessage(message string, statusCode int) error {
	return &Error{
		Err:        errors.New(message),
		StatusCode: statusCode,
	}
}

// NewErrorFromStatusCode creates error from status code
func NewErrorFromStatusCode(statusCode int) error {
	return &Error{
		Err:        errors.New(http.StatusText(statusCode)),
		StatusCode: statusCode,
	}
}

// NewErrorWithFields creates error with status code and error fields
func NewErrorWithFields(err error, statusCode int, fields []FieldError) error {
	return &Error{
		Err:        err,
		StatusCode: statusCode,
		Fields:     fields,
	}
}

// Error - returns error as string
func (e *Error) Error() string {
	var fieldsMsg string
	if len(e.Fields) > 0 {
		fieldsMsg = fmt.Sprintf(" - fields: %v", e.Fields)
	}
	return fmt.Sprintf("%s%s", e.Err.Error(), fieldsMsg)
}
