/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the BME680 sensor
package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/sys/bme680"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	// Create the configuration
	config := gopi.NewAppConfig(MODULE_NAME)

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
