package tasks

import (
	"errors"

	"github.com/homepi/homepi/src/db/models"
	"github.com/stianeikeland/go-rpio"
)

type Task interface {
	Run(pinID rpio.Pin) error
}

// Run the accessory's task
func RunAccessory(a *models.Accessory) (Task, error) {

	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()

	var task Task
	switch a.Task {
	case models.TaskDoor:
		task = &Door{Accessory: a}
	case models.TaskLamp:
		task = &Toggle{Accessory: a}
	case models.TaskToggle:
		task = &Toggle{Accessory: a}
	default:
		return nil, errors.New("task is invalid")
	}

	pin := rpio.Pin(a.PinID)
	rpio.PinMode(pin, rpio.Clock)

	return nil, task.Run(pin)
}
