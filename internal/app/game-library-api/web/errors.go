package web

import "fmt"

// FieldError represents error in a struct field
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is a response in case of an error
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error adds information to request error
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

// NewRequestError is used for creating known error
func NewRequestError(err error, status int) error {
	return &Error{
		Err:    err,
		Status: status,
	}
}

func (e *Error) Error() string {
	var fieldsMsg string
	if len(e.Fields) > 0 {
		fieldsMsg = fmt.Sprintf(" - fields: %v", e.Fields)
	}
	return fmt.Sprintf("%s%s", e.Err.Error(), fieldsMsg)
}
