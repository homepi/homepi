package tasks

import (
	"errors"

	"github.com/homepi/homepi/api/db/models"
	"github.com/stianeikeland/go-rpio"
	"gorm.io/gorm"
)

type Task interface {
	Run(pinID rpio.Pin) error
}

// Run the accessory's task
func RunAccessory(db *gorm.DB, a *models.Accessory) (task Task, err error) {

	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()

	switch a.Task {
	case models.TaskDoor:
		task = &Door{Accessory: a, db: db}
	case models.TaskLamp:
		task = &Toggle{Accessory: a, db: db}
	case models.TaskToggle:
		task = &Toggle{Accessory: a, db: db}
	default:
		return nil, errors.New("task is invalid")
	}

	pin := rpio.Pin(a.PinId)
	rpio.PinMode(pin, rpio.Clock)

	err = task.Run(pin)
	return
}
