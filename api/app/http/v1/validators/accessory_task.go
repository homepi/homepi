package validators

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/homepi/homepi/api/db/models"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	RegisterValidator(&Validator{
		Name:                     "accessory_task",
		CallValidationEvenIfNull: true,
		HandleFunc: func(fl validator.FieldLevel) bool {

			value := fl.Field().Int()

			switch models.Task(value) {
			case models.TaskDoor, models.TaskLamp, models.TaskToggle:
				return true
			}

			return false
		},
		TranslationRegister: func(ut ut.Translator) error {
			return ut.Add("accessory_task", "Task is not valid", true)
		},
		TranslationFunc: func(ut ut.Translator, fe validator.FieldError) (t string) {
			t, _ = ut.T("accessory_task", "task")
			return
		},
	})
}
