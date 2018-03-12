/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme280

import (
	"errors"

	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register bme280 using I2C
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/bme280:i2c",
		Requires: []string{"i2c"},
		Type:     gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("i2c.slave", 0, "I2C Slave address")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			slave, _ := app.AppFlags.GetUint("i2c.slave")
			if slave > 0x7F {
				return nil, errors.New("Invalid -i2c.slave flag")
			}
			return gopi.Open(BME280_I2C{
				Slave: uint8(slave),
				I2C:   app.ModuleInstance("i2c").(gopi.I2C),
			}, app.Logger)
		},
	})

	// Register bme280 using SPI
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/bme280:spi",
		Requires: []string{"spi"},
		Type:     gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("spi.speed", 0, "SPI Communication Speed, Hz")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			speed, _ := app.AppFlags.GetUint("spi.speed")
			return gopi.Open(BME280_SPI{
				Speed: uint32(speed),
				SPI:   app.ModuleInstance("spi").(gopi.SPI),
			}, app.Logger)
		},
	})
}
