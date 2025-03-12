package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// Validator interface that request structs must implement
type Validator interface {
	Validate() (bool, []FieldError)
}

// Sanitizer interface that request structs must implement
type Sanitizer interface {
	Sanitize()
}

// Decode unmarshalls JSON request body
func Decode(log *zap.Logger, r *http.Request, val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(val)
	if err != nil {
		log.Error("Decode request", zap.Error(err))
		return NewErrorFromMessage("invalid request", http.StatusBadRequest)
	}

	// validate
	v, ok := val.(Validator)
	// if it doesn't implement Validator, we assume no validation is needed
	if ok {
		if valid, validationErrors := v.Validate(); !valid {
			return NewErrorWithFields(errors.New("validation error"), http.StatusBadRequest, validationErrors)
		}
	}

	// sanitize
	s, ok := val.(Sanitizer)
	if ok {
		s.Sanitize()
	}

	return nil
}
