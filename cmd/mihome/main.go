/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/sys/ener314rt"
	_ "github.com/djthorpe/sensors/sys/rfm69"
)

////////////////////////////////////////////////////////////////////////////////

func ResetGPIO(app *gopi.AppInstance, gpio gopi.GPIO) error {
	app.Logger.Info("ResetGPIO")

	// Ensure pins are in correct state for SPI0
	gpio.SetPinMode(gopi.GPIOPin(7), gopi.GPIO_OUTPUT)
	gpio.SetPinMode(gopi.GPIOPin(8), gopi.GPIO_OUTPUT)
	gpio.SetPinMode(gopi.GPIOPin(9), gopi.GPIO_ALT0)
	gpio.SetPinMode(gopi.GPIOPin(10), gopi.GPIO_ALT0)
	gpio.SetPinMode(gopi.GPIOPin(11), gopi.GPIO_ALT0)

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if gpio := app.ModuleInstance("gpio").(gopi.GPIO); gpio == nil {
		app.Logger.Error("Missing gpio module")
		return gopi.ErrAppError
	} else {
		// Perform the GPIO Reset
		if err := ResetGPIO(app, gpio); err != nil {
			return err
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("gpio", "sensors/ener314rt")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
