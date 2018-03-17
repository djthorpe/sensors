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
	"fmt"
	"time"

	// Frameworks

	"os"

	"github.com/djthorpe/gopi"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/hw/rfm69"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app := NewApp(app); app == nil {
		return gopi.ErrAppError
	} else if err := app.ResetRadio(); err != nil {
		return err
	} else {
		// Set OOK mode
		if err := app.SetOOKMode(); err != nil {
			return err
		}

		app.SetAddress(0x6C6C6)

		for i := 0; i < 4; i++ {
			fmt.Println("OOK_OFF_ALL")
			if err := app.SendOOK(OOK_OFF_1); err != nil {
				return err
			}

			time.Sleep(time.Second)

			fmt.Println("OOK_ON_ALL")
			if err := app.SendOOK(OOK_ON_1); err != nil {
				return err
			}

			time.Sleep(time.Second)
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/rfm69", "linux/gpio", "protocol/openthings")

	// Add on additional flags
	ConfigFlags(config)

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
