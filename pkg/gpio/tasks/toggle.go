package tasks

import (
	"github.com/homepi/homepi/src/db/models"
	"github.com/stianeikeland/go-rpio"
)

type Toggle struct {
	Accessory *models.Accessory
}

// Run lamp accessory
func (t *Toggle) Run(pin rpio.Pin) error {

	// Set pin as output
	pin.Output()

	// Toggling the pin
	pin.Toggle()

	return nil
}
