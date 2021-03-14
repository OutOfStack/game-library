package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   = validator.New()
	translator *ut.UniversalTranslator
)

func init() {
	enLocale := en.New()
	translator = ut.New(enLocale, enLocale)
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, lang)
}

// Decode unmarshalls JSON request body
func Decode(c *gin.Context, val interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(val)
	if err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
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
