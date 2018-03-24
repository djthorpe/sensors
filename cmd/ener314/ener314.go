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
	"fmt"
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

var (
	MODULE_NAMES = map[string]string{
		"pimote": "sensors/ener314",
		"mihome": "sensors/mihome",
	}
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

	if module_key, exists := app.AppFlags.GetString("interface"); exists == false {
		return errors.New("Missing -interface flag")
	} else if module_value, exists := MODULE_NAMES[module_key]; exists == false {
		return errors.New("Invalid -interface flag")
	} else if device := app.ModuleInstance(module_value).(sensors.ENER314); device == nil {
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
	// Enumerate the modules
	var keys, values []string
	for k, v := range MODULE_NAMES {
		keys = append(keys, k)
		values = append(values, v)
	}

	// Create the configuration
	config := gopi.NewAppConfig(values...)

	// Add on additional flags
	config.AppFlags.FlagString("interface", "mihome", fmt.Sprintf("Interface (%v)", strings.Join(keys, ",")))
	config.AppFlags.FlagBool("on", false, "Switch on")
	config.AppFlags.FlagBool("off", false, "Switch off")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
