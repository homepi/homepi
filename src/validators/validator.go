package validators

import "gopkg.in/go-playground/validator.v9"

var validationMap []*Validator

type Validator struct {
	Name                     string
	CallValidationEvenIfNull bool
	HandleFunc               validator.Func
	TranslationRegister      validator.RegisterTranslationsFunc
	TranslationFunc          validator.TranslationFunc
}

func RegisterValidator(validator *Validator) {
	validationMap = append(validationMap, validator)
}
