/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the Energenie ENER314 over GPIO
package main

import (
	"errors"
	"os"
	"strconv"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/hw/energenie"
	_ "github.com/djthorpe/sensors/hw/rfm69"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	//MODULE_NAME = "sensors/ener314"
	MODULE_NAME = "sensors/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func GetCommand(app *gopi.AppInstance) (string, []uint, error) {

	// Establish command as on or off
	command := ""
	cmd_on_value, cmd_on := app.AppFlags.GetBool("on")
	cmd_off_value, cmd_off := app.AppFlags.GetBool("off")
	if (cmd_on == cmd_off) || (cmd_on && cmd_on_value == false) || (cmd_off && cmd_off_value == false) {
		return "", nil, errors.New("Requires either -on or -off flag")
	} else {
		switch {
		case cmd_on:
			command = "on"
		case cmd_off:
			command = "off"
		}
	}

	// Get socket argument
	if len(app.AppFlags.Args()) == 0 {
		return command, nil, nil
	} else if len(app.AppFlags.Args()) == 1 {
		sockets := make([]uint, 0)
		for _, arg := range strings.Split(app.AppFlags.Args()[0], ",") {
			if socket, err := strconv.ParseUint(arg, 10, 32); err != nil {
				return "", nil, err
			} else {
				sockets = append(sockets, uint(socket))
			}
		}
		return command, sockets, nil
	} else {
		return "", nil, errors.New("Expects zero or one argument of socket numbers")
	}
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.ENER314); device == nil {
		return errors.New("ENER314 module not found")
	} else if command, sockets, err := GetCommand(app); err != nil {
		return err
	} else if command == "on" {
		if err := device.On(sockets...); err != nil {
			return err
		}
	} else if command == "off" {
		if err := device.Off(sockets...); err != nil {
			return err
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig(MODULE_NAME)

	// Add on additional flags
	config.AppFlags.FlagBool("on", false, "Switch on")
	config.AppFlags.FlagBool("off", false, "Switch off")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
