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

	// Frameworks
	"github.com/djthorpe/gopi"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/hw/energenie"
	_ "github.com/djthorpe/sensors/hw/rfm69"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if mihome := app.ModuleInstance("sensors/mihome"); mihome == nil {
		return gopi.ErrAppError
	} else {
		app.Logger.Info("mihome=%v", mihome)
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/mihome")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
