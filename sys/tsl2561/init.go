/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package tsl2561

import (
	"errors"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register bme280 using I2C
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/tsl2561",
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
			return gopi.Open(TSL2561{
				Slave: uint8(slave),
				I2C:   app.ModuleInstance("i2c").(gopi.I2C),
			}, app.Logger)
		},
	})
}
