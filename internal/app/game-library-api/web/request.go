package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Decode unmarshalls JSON request body
func Decode(c *gin.Context, val interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(val)
	if err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err = validate.Struct(val); err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, vErr := range vErrs {
			field := FieldError{
				Field: vErr.Field(),
				Error: vErr.Translate(lang),
			}
			fields = append(fields, field)
		}
		return &Error{
			Err:    errors.New("validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}
