package validators

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

var (
	trans    ut.Translator
	validate *validator.Validate
)

func Configure() (err error) {

	validate = validator.New()

	english := en.New()
	uni := ut.New(english, english)

	trans, _ = uni.GetTranslator("en")
	err = enTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return
	}

	err = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "The {0} field is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) (t string) {
		t, _ = ut.T("required", fe.Field())
		return
	})
	if err != nil {
		return
	}

	err = validate.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
		return ut.Add("required_without", "One of the field must be present. {0} or {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) (t string) {
		t, _ = ut.T("required_without", fe.Field(), fe.Param())
		return
	})
	if err != nil {
		return
	}

	err = validate.RegisterTranslation("required_without_all", trans, func(ut ut.Translator) error {
		return ut.Add("required_without_all", "{0} is a required field without {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) (t string) {
		param := strings.ReplaceAll(fe.Param(), " ", " and ")
		t, _ = ut.T("required_without_all", fe.Field(), param)
		return
	})
	if err != nil {
		return
	}

	for _, validation := range validationMap {
		err = validate.RegisterTranslation(
			validation.Name,
			trans,
			validation.TranslationRegister,
			validation.TranslationFunc)
		if err != nil {
			return
		}
		err = validate.RegisterValidation(
			validation.Name,
			validation.HandleFunc,
			validation.CallValidationEvenIfNull)
		if err != nil {
			return
		}
	}

	return
}
