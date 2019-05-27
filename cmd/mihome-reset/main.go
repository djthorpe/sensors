/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Control Energenie MiHome devices
package main

import (
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Register modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func ResetRadio(gpio gopi.GPIO, pin gopi.GPIOPin) error {
	// If reset is not defined, then return not implemented
	if pin == gopi.GPIO_PIN_NONE {
		return gopi.ErrNotImplemented
	}

	// Ensure pin is output
	gpio.SetPinMode(pin, gopi.GPIO_OUTPUT)

	// Pull reset high for 100ms and then low for 5ms
	gpio.WritePin(pin, gopi.GPIO_HIGH)
	time.Sleep(time.Millisecond * 100)
	gpio.WritePin(pin, gopi.GPIO_LOW)
	time.Sleep(time.Millisecond * 5)

	return nil
}

func ResetSPI(gpio gopi.GPIO) error {
	// Ensure pins are in correct state for SPI0
	gpio.SetPinMode(gopi.GPIOPin(7), gopi.GPIO_OUTPUT)
	gpio.SetPinMode(gopi.GPIOPin(8), gopi.GPIO_OUTPUT)
	gpio.SetPinMode(gopi.GPIOPin(9), gopi.GPIO_ALT0)
	gpio.SetPinMode(gopi.GPIOPin(10), gopi.GPIO_ALT0)
	gpio.SetPinMode(gopi.GPIOPin(11), gopi.GPIO_ALT0)

	return nil
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	if gpio := app.ModuleInstance("gpio").(gopi.GPIO); gpio == nil {
		app.Logger.Error("Missing gpio module")
		return gopi.ErrAppError
	} else if err := ResetSPI(gpio); err != nil {
		return err
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("gpio")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
