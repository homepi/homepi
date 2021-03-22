package tasks

import (
	"github.com/homepi/homepi/api/db/models"
	"github.com/stianeikeland/go-rpio"
	"gorm.io/gorm"
)

type Toggle struct {
	db        *gorm.DB
	Accessory *models.Accessory
}

// Run lamp accessory
func (t *Toggle) Run(pin rpio.Pin) error {

	// Set pin as output
	pin.Output()

	// Toggling the pin
	pin.Toggle()

	t.Accessory.UpdateStatus(t.db, pin.Read())

	return nil
}
