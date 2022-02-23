package tasks

import (
	"time"

	"github.com/homepi/homepi/src/db/models"
	"github.com/stianeikeland/go-rpio"
)

type Door struct {
	Accessory *models.Accessory
}

// Run door accessory.
func (d *Door) Run(pin rpio.Pin) error {
	// Set pin as output
	pin.Output()

	// Turn on the pin!
	pin.Low()

	// Sleep 1 second before turning off the pin
	time.Sleep(time.Second)

	// Turn off the pin
	pin.High()

	return nil
}
