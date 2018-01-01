/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the Energenie ENER314 over GPIO
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/energenie"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MODULE_NAME = "sensors/ener314"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.ENER314); device == nil {
		return errors.New("ENER314 module not found")
	} else {
		if socket_on, exists := app.AppFlags.GetUint("on"); exists {
			if err := device.On(socket_on); err != nil {
				return err
			}
		} else if socket_off, exists := app.AppFlags.GetUint("off"); exists {
			if err := device.Off(socket_off); err != nil {
				return err
			}
		} else {
			return errors.New("Expect either -on or -off flag")
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig(MODULE_NAME)
	config.AppFlags.FlagUint("on", 0, "Switch on")
	config.AppFlags.FlagUint("off", 0, "Switch off")

	// Create the application
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(runLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
