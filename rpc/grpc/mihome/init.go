/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	"fmt"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register server
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/service/mihome",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server", "gpio", "sensors/ener314rt"},
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("mode", "", "RX mode")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			// Convert mode to a MiHomeMode value
			mode_, _ := app.AppFlags.GetString("mode")
			if mode, err := miHomeModeFromString(mode_); err != nil {
				return nil, err
			} else {
				return gopi.Open(Service{
					Server: app.ModuleInstance("rpc/server").(gopi.RPCServer),
					MiHome: app.ModuleInstance("sensors/ener314rt").(sensors.MiHome),
					Mode:   mode,
				}, app.Logger)
			}
		},
	})

	// Register client
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/client/mihome",
		Type:     gopi.MODULE_TYPE_CLIENT,
		Requires: []string{"rpc/clientpool"},
		Run: func(app *gopi.AppInstance, _ gopi.Driver) error {
			if clientpool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool); clientpool == nil {
				return gopi.ErrAppError
			} else {
				clientpool.RegisterClient("sensors.MiHome", NewMiHomeClient)
				return nil
			}
		},
	})
}

// stringFromMode returns an upper-case string from a MiHomeMode
// or returns an empty string otherwise
func stringFromMiHomeMode(mode sensors.MiHomeMode) string {
	return strings.TrimPrefix(fmt.Sprint(mode), "MIHOME_MODE_")
}

// modeFromString returns MiHomeMode given a string, or returns
// an error otherwise. Case-insensitive
func miHomeModeFromString(value string) (sensors.MiHomeMode, error) {
	value_upper := strings.TrimSpace(strings.ToUpper(value))
	all_modes := make([]string, 0)
	for mode := sensors.MIHOME_MODE_NONE; mode <= sensors.MIHOME_MODE_MAX; mode++ {
		if mode_string := stringFromMiHomeMode(mode); mode_string == value_upper {
			// Return mode
			return mode, nil
		} else {
			all_modes = append(all_modes, strings.ToLower(mode_string))
		}
	}
	// Return error
	return sensors.MIHOME_MODE_NONE, fmt.Errorf("Invalid -mode value, values are %v", strings.Join(all_modes, ", "))
}
