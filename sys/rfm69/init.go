/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register RFM69 communication through SPI
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/rfm69/spi",
		Requires: []string{"spi"},
		Type:     gopi.MODULE_TYPE_OTHER,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(RFM69{
				SPI: app.ModuleInstance("spi").(gopi.SPI),
			}, app.Logger)
		},
	})
}
