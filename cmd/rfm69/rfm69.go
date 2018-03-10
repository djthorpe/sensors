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
	"github.com/olekukonko/tablewriter"

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

func setParametersMode(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("mode"); exists == false {
		return nil
	} else if mode, err := stringToMode(value); err != nil {
		return err
	} else if err := device.SetMode(mode); err != nil {
		return err
	}

	// Success
	return nil
}

func setParametersDataMode(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("datamode"); exists == false {
		return nil
	} else if mode, err := stringToDataMode(value); err != nil {
		return err
	} else if err := device.SetDataMode(mode); err != nil {
		return err
	}

	// Success
	return nil
}

func setParametersModulation(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("modulation"); exists == false {
		return nil
	} else if modulation, err := stringToModulation(value); err != nil {
		return err
	} else if err := device.SetModulation(modulation); err != nil {
		return err
	}

	// Success
	return nil
}

func setParameters(app *gopi.AppInstance, device sensors.RFM69) error {
	if err := setParametersMode(app, device); err != nil {
		return err
	}
	if err := setParametersDataMode(app, device); err != nil {
		return err
	}
	if err := setParametersModulation(app, device); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func status(device sensors.RFM69) error {

	// Output register information
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Register", "Value"})

	// Output mode, data mode and modulation
	table.Append([]string{"mode", modeToString(device.Mode())})
	table.Append([]string{"datamode", dataModeToString(device.DataMode())})
	table.Append([]string{"modulation", modulationToString(device.Modulation())})

	// Node and Broadcast addresses
	table.Append([]string{"node_addr", fmt.Sprintf("%02X", device.NodeAddress())})
	table.Append([]string{"broadcast_addr", fmt.Sprintf("%02X", device.BroadcastAddress())})

	table.Render()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.RFM69); device == nil {
		return errors.New("Module not found: " + MODULE_NAME)
	} else if err := setParameters(app, device); err != nil {
		return err
	} else if err := status(device); err != nil {
		return err
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the RFM69 instance
	config := gopi.NewAppConfig(MODULE_NAME)

	// Parameters
	config.AppFlags.FlagString("mode", "", "Device Mode (sleep,standby,fs,tx,rx)")
	config.AppFlags.FlagString("datamode", "", "Data Mode (package,nosync,sync)")
	config.AppFlags.FlagString("modulation", "", "Modulation (fsk,fsk_1.0,fsk_0.5,fsk_0.3,ook,ook_br,ook_2br)")
	config.AppFlags.FlagString("node_addr", "", "Node Address")
	config.AppFlags.FlagString("broadcast_addr", "", "Broadcast Address")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
