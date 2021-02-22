package web

// ErrorResponse is a response in case of an error
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error adds information to request error
type Error struct {
	Err    error
	Status int
}

// NewErrorRequest is used for creating known error
func NewErrorRequest(err error, status int) error {
	return &Error{
		Err:    err,
		Status: status,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
