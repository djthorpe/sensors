/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package energenie

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register pimote using GPIO
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/ener314",
		Requires: []string{"gpio"},
		Type:     gopi.MODULE_TYPE_OTHER,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(ENER314{
				GPIO: app.ModuleInstance("gpio").(gopi.GPIO),
			}, app.Logger)
		},
	})
}
