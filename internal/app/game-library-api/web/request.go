package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api/validation"
	"github.com/OutOfStack/game-library/internal/appconf"
	"go.uber.org/zap"
)

// Decoder decodes JSON request
type Decoder struct {
	log       *zap.Logger
	validator *validation.Validator
}

// NewDecoder creates new decoder
func NewDecoder(log *zap.Logger, cfg *appconf.Cfg) *Decoder {
	return &Decoder{
		log:       log,
		validator: validation.NewValidator(log, cfg),
	}
}

// Validatable interface that request structs must implement
type Validatable interface {
	ValidateWith(v *validation.Validator) (bool, []FieldError)
}

// Sanitizable interface that request structs must implement
type Sanitizable interface {
	Sanitize()
}

// Decode unmarshalls JSON request body
func (d *Decoder) Decode(r *http.Request, val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(val)
	if err != nil {
		d.log.Error("Decode request", zap.Error(err))
		return NewErrorFromMessage("invalid request", http.StatusBadRequest)
	}

	// validate
	if v, ok := val.(Validatable); ok {
		// if it doesn't implement Validator, assume no validation is needed
		if valid, validationErrors := v.ValidateWith(d.validator); !valid {
			return NewErrorWithFields(errors.New("validation error"), http.StatusBadRequest, validationErrors)
		}
	}

	// sanitize
	if s, ok := val.(Sanitizable); ok {
		s.Sanitize()
	}

	return nil
}
