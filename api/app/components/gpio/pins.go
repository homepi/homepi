package gpio

import (
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func GetPins(db *gorm.DB) *models.Pins {
	pins := &models.Pins{
		Top: []*models.Pin{
			{Id: 0, Disable: true, Used: false, Name: "5V, Power"},
			{Id: 0, Disable: true, Used: false, Name: "5V, Power"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_14, UART0_TXD"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_15, UART0_TXD"},
			{Id: 18, Disable: false, Used: false, Name: "GPIO_18"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 23, Disable: false, Used: false, Name: "GPIO_23"},
			{Id: 24, Disable: false, Used: false, Name: "GPIO_24"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 25, Disable: false, Used: false, Name: "GPIO_25"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_8, SPI0_CE0_N"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_7, SPI0_CE0_N"},
			{Id: 0, Disable: true, Used: false, Name: "ID_SC, I2C ID EEPROM"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 12, Disable: false, Used: false, Name: "GPIO_12"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 16, Disable: false, Used: false, Name: "GPIO_16"},
			{Id: 20, Disable: false, Used: false, Name: "GPIO_20"},
			{Id: 21, Disable: false, Used: false, Name: "GPIO_21"},
		},
		Bottom: []*models.Pin{
			{Id: 0, Disable: true, Used: false, Name: "3V3, Power"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_2, SDA1 I2C"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_3, SDA1 I2C"},
			{Id: 4, Disable: false, Used: false, Name: "GPIO_4"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 17, Disable: false, Used: false, Name: "GPIO_17"},
			{Id: 27, Disable: false, Used: false, Name: "GPIO_27"},
			{Id: 22, Disable: false, Used: false, Name: "GPIO_22"},
			{Id: 0, Disable: true, Used: false, Name: "3V3 Power"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_10, SPI0_MOSI"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_9, SPI0_MOSI"},
			{Id: 0, Disable: true, Used: false, Name: "GPIO_11, SPI0_SCLK"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
			{Id: 0, Disable: true, Used: false, Name: "ID_SD, I2C ID EEPROM"},
			{Id: 5, Disable: false, Used: false, Name: "GPIO_5"},
			{Id: 6, Disable: false, Used: false, Name: "GPIO_16"},
			{Id: 13, Disable: false, Used: false, Name: "GPIO_12"},
			{Id: 19, Disable: false, Used: false, Name: "GPIO_19"},
			{Id: 26, Disable: false, Used: false, Name: "GPIO_26"},
			{Id: 0, Disable: true, Used: false, Name: "Groud"},
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
				pins.Bottom[index].Used = true
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
