package validators

import (
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func init() {

	RegisterValidator(&Validator{
		Name:    "pin_type",
		CallValidationEvenIfNull: true,
		HandleFunc: func(fl validator.FieldLevel) bool {

			value := fl.Field().String()

			if value != "" {
				switch value {
				case "1", "2":
					return true
				}
			}

			return false
		},
		TranslationRegister: func(ut ut.Translator) error {
			return ut.Add("pin_type", "{0} is not a valid pin type", true)
		},
		TranslationFunc: func(ut ut.Translator, fe validator.FieldError) (t string) {
			t, _ = ut.T("pin_type", fe.Value().(string))
			return
		},
	})
}