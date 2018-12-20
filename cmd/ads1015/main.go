/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the ADS1015 ADC converter
package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/i2c"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/sys/ads1x15"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Get sensor
	adc := app.ModuleInstance("sensors/ads1015").(sensors.ADS1015)
	fmt.Println(adc)

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/ads1015")

	config.AppFlags.SetUsageFunc(func(flags *gopi.Flags) {
		fmt.Fprintf(os.Stderr, "Usage: %v <flags>\n\n", flags.Name())
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flags.PrintDefaults()
	})

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, MainLoop))
}
