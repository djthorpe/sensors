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
	// Register mihome module
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/mihome",
		Type:     gopi.MODULE_TYPE_OTHER,
		Requires: []string{"sensors/ener314rt"},
		Config: func(config *gopi.AppConfig) {
			// MiHome flags
			config.AppFlags.FlagString("mihome.mode", "monitor", "RX mode")
			config.AppFlags.FlagUint("mihome.repeat", 0, "Default TX Repeat")
			config.AppFlags.FlagFloat64("mihome.tempoffset", 0, "Temperature Calibration Value")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			// Convert mode to a MiHomeMode value
			mode_, _ := app.AppFlags.GetString("mihome.mode")
			repeat, _ := app.AppFlags.GetUint("mihome.repeat")
			tempoffset, _ := app.AppFlags.GetFloat64("mihome.tempoffset")
			if mode, err := miHomeModeFromString(mode_); err != nil {
				return nil, err
			} else {
				return gopi.Open(MiHome{
					Radio:      app.ModuleInstance("sensors/ener314rt").(sensors.ENER314RT),
					Mode:       mode,
					Repeat:     repeat,
					TempOffset: float32(tempoffset),
				}, app.Logger)
			}
		},
		Run: func(app *gopi.AppInstance, driver gopi.Driver) error {
			// Register protocols with driver. Codecs have OTHER as module type
			// and name starting with "sensors/protocol"
			for _, module := range gopi.ModulesByType(gopi.MODULE_TYPE_OTHER) {
				if strings.HasPrefix(module.Name, "sensors/protocol/") == false {
					continue
				}
				// Get protocol instance and register it
				if proto, ok := app.ModuleInstance(module.Name).(sensors.Proto); ok == false {
					return fmt.Errorf("Invalid protocol: %v: %v", module.Name, proto)
				} else if err := driver.(sensors.MiHome).AddProto(proto); err != nil {
					return err
				}
			}
			// Return success
			return nil
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
	return sensors.MIHOME_MODE_NONE, fmt.Errorf("Invalid -mihome.mode value: values are %v", strings.Join(all_modes, ", "))
}
