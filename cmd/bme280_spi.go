/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the BME280 sensor over the SPI bus
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/bme280"
)

////////////////////////////////////////////////////////////////////////////////

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	if bme280 := app.ModuleInstance("sensors/bme280/spi"); bme280 == nil {
		return errors.New("BME280 module not found")
	}

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/bme280/spi")
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
