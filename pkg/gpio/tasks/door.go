package tasks

import (
	"time"

	"github.com/homepi/homepi/src/db/models"
	"github.com/stianeikeland/go-rpio"
)

type Door struct {
	Accessory *models.Accessory
}

// Run door accessory
func (d *Door) Run(pin rpio.Pin) error {

	// Set pin as output
	pin.Output()

	// Turn on the pin!
	pin.Low()

	// Turn off the pin after 1 second!
	time.AfterFunc(time.Second, pin.High)

	return nil
}
