/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register server
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/mihome:service",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server", "sensors/mihome"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Service{
				Server: app.ModuleInstance("rpc/server").(gopi.RPCServer),
				MiHome: app.ModuleInstance("sensors/mihome").(sensors.MiHome),
			}, app.Logger)
		},
	})
	// Register client
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/mihome:client",
		Type:     gopi.MODULE_TYPE_CLIENT,
		Requires: []string{"rpc/clientpool"},
		Run: func(app *gopi.AppInstance, _ gopi.Driver) error {
			if clientpool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool); clientpool == nil {
				return gopi.ErrAppError
			} else {
				clientpool.RegisterClient("mihome.MiHome", NewMiHomeClient)
				return nil
			}
		},
	})

}
