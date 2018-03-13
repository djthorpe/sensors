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
	// Frameworks

	"os"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/hw/rfm69"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if rfm69 := app.ModuleInstance("sensors/rfm69").(sensors.RFM69); rfm69 == nil {
		return gopi.ErrAppError
	} else if openthings := app.ModuleInstance("protocol/openthings"); openthings == nil {
		return gopi.ErrAppError
	} else {
		// TODO: Reset the device

		// Set FSK mode
		if err := SetFSKMode(rfm69); err != nil {
			return err
		}

		// Report
		app.Logger.Info("rfm69=%v", rfm69)
		app.Logger.Info("openthings=%v", openthings)

	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/rfm69", "protocol/openthings")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
