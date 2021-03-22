package tasks

import (
	"time"

	"github.com/homepi/homepi/api/db/models"
	"github.com/stianeikeland/go-rpio"
	"gorm.io/gorm"
)

type Door struct {
	db        *gorm.DB
	Accessory *models.Accessory
}

// Run door accessory
func (d *Door) Run(pin rpio.Pin) error {

	// Set pin as output
	pin.Output()

	// Turn on the pin!
	pin.Low()

	select {
	case <-time.After(time.Second):
		// Turn off the pin after 1 second!
		pin.High()
	}

	return nil
}
