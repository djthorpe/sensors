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
	_ "github.com/djthorpe/sensors/hw/tsl2561"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MODULE_NAME = "sensors/tsl2561"
)

const (
	COMMAND_MEASURE = iota
	COMMAND_STATUS
	COMMAND_HELP
)

////////////////////////////////////////////////////////////////////////////////
// CONVERT INTO VALUES

func GetGainFromUint(value uint) (sensors.TSL2561Gain, error) {
	switch value {
	case 1:
		return sensors.TSL2561_GAIN_1, nil
	case 16:
		return sensors.TSL2561_GAIN_16, nil
	default:
		return sensors.TSL2561_GAIN_MAX, fmt.Errorf("Invalid gain value: %v", value)
	}
}

func GetIntegrateTimeFromFloat(value float64) (sensors.TSL2561IntegrateTime, error) {
	switch value {
	case 13.7:
		return sensors.TSL2561_INTEGRATETIME_13P7MS, nil
	case 101:
		return sensors.TSL2561_INTEGRATETIME_101MS, nil
	case 402:
		return sensors.TSL2561_INTEGRATETIME_402MS, nil
	default:
		return sensors.TSL2561_INTEGRATETIME_MAX, fmt.Errorf("Invalid integrate_time value: %v", value)
	}
}

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
// SET PARAMETERS

func set_gain(app *gopi.AppInstance, device sensors.TSL2561) error {
	if gain, exists := app.AppFlags.GetUint("gain"); exists {
		if value, err := GetGainFromUint(gain); err != nil {
			return err
		} else if err := device.SetGain(value); err != nil {
			return err
		}
	}
	return nil
}

func set_integrate_time(app *gopi.AppInstance, device sensors.TSL2561) error {
	if integrate_time, exists := app.AppFlags.GetFloat64("integrate_time"); exists {
		if value, err := GetIntegrateTimeFromFloat(integrate_time); err != nil {
			return err
		} else if err := device.SetIntegrateTime(value); err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	// Determine the command to run
	command := COMMAND_HELP
	if args := app.AppFlags.Args(); len(args) == 1 && args[0] == "status" {
		command = COMMAND_STATUS
	} else if len(args) == 0 || len(args) == 1 && args[0] == "measure" {
		command = COMMAND_MEASURE
	}

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.TSL2561); device == nil {
		return errors.New("TSL2561 module not found")
	} else {
		// set parameters
		if err := set_gain(app, device); err != nil {
			return err
		} else if err := set_integrate_time(app, device); err != nil {
			return err
		}

		// run command
		switch command {
		case COMMAND_MEASURE:
			if err := measure(device); err != nil {
				return err
			}
		case COMMAND_STATUS:
			if err := status(device); err != nil {
				return err
			}
		default:
			return gopi.ErrHelp
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
