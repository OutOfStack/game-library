package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Decode unmarshalls JSON request body
func Decode(r *http.Request, val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(val)
	if err != nil {
		return NewError(err, http.StatusBadRequest)
	}

	if err = validate.Struct(val); err != nil {
		var vErrs validator.ValidationErrors
		if ok := errors.As(err, &vErrs); !ok {
			return err
		}

		var fields []FieldError
		for _, vErr := range vErrs {
			fields = append(fields, FieldError{
				Field: vErr.Field(),
				Error: vErr.Translate(lang),
			})
		}
		return NewErrorWithFields(errors.New("validation error"), http.StatusBadRequest, fields)
	}

	return nil
}
