package web

import (
	"cloud.google.com/go/civil"
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
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
	en_translations.RegisterDefaultTranslations(validate, lang)

	validate.RegisterValidation("date", validateDate)

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
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = validate.RegisterTranslation(tag, lang, registerFn, transFn)
}
