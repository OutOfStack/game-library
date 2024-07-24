package web

import (
	"log"

	"cloud.google.com/go/civil"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   = validator.New()
	translator *ut.UniversalTranslator
	lang       ut.Translator
)

const (
	dateFieldLength = 10
)

func init() {
	enLocale := en.New()
	translator = ut.New(enLocale, enLocale)
	lang, _ = translator.GetTranslator("en")
	if err := entranslations.RegisterDefaultTranslations(validate, lang); err != nil {
		log.Fatalf("register default translation for validator: %v", err)
	}

	if err := validate.RegisterValidation("date", validateDate); err != nil {
		log.Fatalf("register 'date' validation: %v", err)
	}

	errDate := "{0} must be in format 'YYYY-MM-DD'"
	addTranslation("date", errDate)
}

// validateDate is a rule for validating date format (YYYY-MM-DD)
func validateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()

	if len(date) != dateFieldLength {
		return false
	}

	_, err := civil.ParseDate(date)
	return err == nil
}

// addTranslation adds error message for specified tag
func addTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag = fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.Error()
		}
		return t
	}

	_ = validate.RegisterTranslation(tag, lang, registerFn, transFn)
}
