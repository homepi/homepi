package validators

import (
	"github.com/homepi/homepi/pkg/libstr"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidator(validatorStruct interface{}) (errors map[string]interface{}) {
	errors = make(map[string]interface{})
	if err := validate.Struct(validatorStruct); err != nil {
		switch err := err.(type) {
		case validator.ValidationErrors:
			for _, e := range err {
				errors[libstr.ToSnakeCase(e.Field())] = e.Translate(trans)
			}
		}
	}
	return
}
