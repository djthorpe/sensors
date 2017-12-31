/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the BME280 sensor over the I2C bus
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
	"github.com/olekukonko/tablewriter"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/tsl2561"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MODULE_NAME = "sensors/tsl2561"
)

////////////////////////////////////////////////////////////////////////////////
// STATUS

func status(device sensors.TSL2561) error {

	// Output register information
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Register", "Value"})

	chip_id, chip_version := device.ChipIDVersion()
	table.Append([]string{"chip_id", fmt.Sprintf("0x%02X", chip_id)})
	table.Append([]string{"chip_version", fmt.Sprintf("0x%02X", chip_version)})
	table.Append([]string{"integrate_time", fmt.Sprint(device.IntegrateTime())})
	table.Append([]string{"gain", fmt.Sprint(device.Gain())})

	table.Render()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SAMPLE

func measure(device sensors.TSL2561) error {

	if lux, err := device.ReadSample(); err != nil {
		return err
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.SetHeader([]string{"Measurement", "Value"})
		table.Append([]string{"illuminance", fmt.Sprintf("%.2f Lux", lux)})
		table.Render()
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.TSL2561); device == nil {
		return errors.New("TSL2561 module not found")
	} else if err := status(device); err != nil {
		return err
	} else if err := measure(device); err != nil {
		return err
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
	config.AppFlags.FlagUint("gain", 0, "Sample gain (1,16)")
	config.AppFlags.FlagFloat64("integrate_time", 0, "Integration time, milliseconds (13.7, 101 or 402)")

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
