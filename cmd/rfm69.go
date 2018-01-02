/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the RFM69 device
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
	_ "github.com/djthorpe/sensors/rfm69"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MODULE_NAME = "sensors/rfm69"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.RFM69); device == nil {
		return errors.New("Module not found: " + MODULE_NAME)
	} else {
		fmt.Println(device)
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
