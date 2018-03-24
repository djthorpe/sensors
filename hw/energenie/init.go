/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package energenie

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
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

	// Register mihome using SPI & RFM69
	gopi.RegisterModule(gopi.Module{
		Name:     "sensors/mihome",
		Requires: []string{"gpio", "sensors/rfm69"},
		Type:     gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			// GPIO pin configurations
			config.AppFlags.FlagUint("gpio.reset", 25, "Reset Pin (Logical)")
			config.AppFlags.FlagUint("gpio.led1", 27, "Green LED Pin (Logical)")
			config.AppFlags.FlagUint("gpio.led2", 22, "Red LED Pin (Logical)")

			// MiHome flags
			config.AppFlags.FlagString("mihome.cid", "", "20-bit Command Device ID (hexadecimal)")
			config.AppFlags.FlagUint("mihome.repeat", 0, "Command TX Repeat")

			// Default spi.slave to 1
			if err := config.AppFlags.SetUint("spi.slave", 1); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			if gpio, ok := app.ModuleInstance("gpio").(gopi.GPIO); !ok {
				return nil, fmt.Errorf("Missing or invalid GPIO module")
			} else if radio, ok := app.ModuleInstance("sensors/rfm69").(sensors.RFM69); !ok {
				return nil, fmt.Errorf("Missing or invalid Radio module")
			} else {
				config := MiHome{
					GPIO:     gpio,
					Radio:    radio,
					PinReset: gopi.GPIO_PIN_NONE,
					PinLED1:  gopi.GPIO_PIN_NONE,
					PinLED2:  gopi.GPIO_PIN_NONE,
				}
				if reset, _ := app.AppFlags.GetUint("gpio.reset"); reset > 0 && reset <= 0xFF {
					config.PinReset = gopi.GPIOPin(reset)
				}
				if led1, _ := app.AppFlags.GetUint("gpio.led1"); led1 > 0 && led1 <= 0xFF {
					config.PinLED1 = gopi.GPIOPin(led1)
				}
				if led2, _ := app.AppFlags.GetUint("gpio.led2"); led2 > 0 && led2 <= 0xFF {
					config.PinLED2 = gopi.GPIOPin(led2)
				}
				if cid, exists := app.AppFlags.GetString("mihome.cid"); exists {
					config.CID = cid
				}
				if repeat, exists := app.AppFlags.GetUint("mihome.repeat"); exists {
					config.Repeat = repeat
				}
				return gopi.Open(config, app.Logger)
			}
		},
	})

}
