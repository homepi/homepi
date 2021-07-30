package gpio

import (
	"github.com/homepi/homepi/src/db/models"
	"gorm.io/gorm"
)

func GetPins(db *gorm.DB) *models.Pins {
	pins := &models.Pins{
		Top: []*models.Pin{
			{ID: 0, Disable: true, Used: false, Name: "5V, Power"},
			{ID: 0, Disable: true, Used: false, Name: "5V, Power"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_14, UART0_TXD"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_15, UART0_TXD"},
			{ID: 18, Disable: false, Used: false, Name: "GPIO_18"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 23, Disable: false, Used: false, Name: "GPIO_23"},
			{ID: 24, Disable: false, Used: false, Name: "GPIO_24"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 25, Disable: false, Used: false, Name: "GPIO_25"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_8, SPI0_CE0_N"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_7, SPI0_CE0_N"},
			{ID: 0, Disable: true, Used: false, Name: "ID_SC, I2C ID EEPROM"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 12, Disable: false, Used: false, Name: "GPIO_12"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 16, Disable: false, Used: false, Name: "GPIO_16"},
			{ID: 20, Disable: false, Used: false, Name: "GPIO_20"},
			{ID: 21, Disable: false, Used: false, Name: "GPIO_21"},
		},
		Bottom: []*models.Pin{
			{ID: 0, Disable: true, Used: false, Name: "3V3, Power"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_2, SDA1 I2C"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_3, SDA1 I2C"},
			{ID: 4, Disable: false, Used: false, Name: "GPIO_4"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 17, Disable: false, Used: false, Name: "GPIO_17"},
			{ID: 27, Disable: false, Used: false, Name: "GPIO_27"},
			{ID: 22, Disable: false, Used: false, Name: "GPIO_22"},
			{ID: 0, Disable: true, Used: false, Name: "3V3 Power"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_10, SPI0_MOSI"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_9, SPI0_MOSI"},
			{ID: 0, Disable: true, Used: false, Name: "GPIO_11, SPI0_SCLK"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
			{ID: 0, Disable: true, Used: false, Name: "ID_SD, I2C ID EEPROM"},
			{ID: 5, Disable: false, Used: false, Name: "GPIO_5"},
			{ID: 6, Disable: false, Used: false, Name: "GPIO_16"},
			{ID: 13, Disable: false, Used: false, Name: "GPIO_12"},
			{ID: 19, Disable: false, Used: false, Name: "GPIO_19"},
			{ID: 26, Disable: false, Used: false, Name: "GPIO_26"},
			{ID: 0, Disable: true, Used: false, Name: "Groud"},
		},
	}
	CheckPins(db, pins)
	return pins
}

func CheckPins(db *gorm.DB, pins *models.Pins) {

	for index, pin := range pins.Top {
		if !pin.Disable {
			used, err := pin.Check(db)
			if err != nil {
				continue
			}
			if used {
				pins.Top[index].Used = true
			}
		}
	}

	for index, pin := range pins.Bottom {
		if !pin.Disable {
			used, err := pin.Check(db)
			if err != nil {
				continue
			}
			if used {
				pins.Bottom[index].Used = true
			}
		}
	}

}
