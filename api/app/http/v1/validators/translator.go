package validators

import (
	"github.com/homepi/homepi/api/app/components/strings"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidator(validatorStruct interface{}) (errors map[string]interface{}) {

	errors = make(map[string]interface{})
	err := validate.Struct(validatorStruct)

	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			for _, e := range err.(validator.ValidationErrors) {
				errors[strings.ToSnakeCase(e.Field())] = e.Translate(trans)
			}
		}
	}

	return
}
