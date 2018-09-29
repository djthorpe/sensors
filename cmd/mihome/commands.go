/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"io"
	"strconv"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type CommandFunc struct {
	Fn          func(*gopi.AppInstance, []string) error
	Description string
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	COMMANDS = map[string]CommandFunc{
		"reset_gpio":   CommandFunc{ResetGPIO, "Reset GPIO"},
		"reset_radio":  CommandFunc{ResetRadio, "Reset RFM69 Radio"},
		"measure_temp": CommandFunc{MeasureTemp, "Measure Temperature"},
		"on":           CommandFunc{TransmitOn, "On TX (optionally use 1,2,3,4 as additional argument)"},
		"off":          CommandFunc{TransmitOff, "Off TX (optionally use 1,2,3,4 as additional argument)"},
	}
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func CommandNames() []string {
	commands := make([]string, 0, len(COMMANDS))
	for k, _ := range COMMANDS {
		commands = append(commands, k)
	}
	return commands
}

func PrintCommands(out io.Writer) {
	fmt.Fprintf(out, "Commands:\n")
	for k, v := range COMMANDS {
		fmt.Fprintf(out, "  %-10s\t%s\n", k, v.Description)
	}
}

func RunCommand(app *gopi.AppInstance, cmd string, args []string) error {
	if command, exists := COMMANDS[cmd]; exists == false {
		return gopi.ErrHelp
	} else if err := command.Fn(app, args); err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INSTANCES

func GPIO(app *gopi.AppInstance) gopi.GPIO {
	return app.ModuleInstance("gpio").(gopi.GPIO)
}

func MiHome(app *gopi.AppInstance) sensors.MiHome {
	return app.ModuleInstance("sensors/ener314rt").(sensors.MiHome)
}

////////////////////////////////////////////////////////////////////////////////
// UTILTY FUNCTIONS

func toSockets(args []string) ([]uint, error) {
	// Where there are no arguments, assume 'all'
	if len(args) == 0 {
		return nil, nil
	}

	// Else read uintegers into an array
	ret := make([]uint, 0, len(args))
	for _, value := range args {
		if uint_value, err := strconv.ParseUint(value, 10, 64); err != nil {
			return nil, err
		} else {
			ret = append(ret, uint(uint_value))
		}
	}

	// Return success
	return ret, nil
}

////////////////////////////////////////////////////////////////////////////////
// COMMANDS

func ResetGPIO(app *gopi.AppInstance, args []string) error {
	app.Logger.Info("ResetGPIO")
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	if gpio := GPIO(app); gpio == nil {
		return gopi.ErrAppError
	} else {
		// Ensure pins are in correct state for SPI0
		gpio.SetPinMode(gopi.GPIOPin(7), gopi.GPIO_OUTPUT)
		gpio.SetPinMode(gopi.GPIOPin(8), gopi.GPIO_OUTPUT)
		gpio.SetPinMode(gopi.GPIOPin(9), gopi.GPIO_ALT0)
		gpio.SetPinMode(gopi.GPIOPin(10), gopi.GPIO_ALT0)
		gpio.SetPinMode(gopi.GPIOPin(11), gopi.GPIO_ALT0)

		// Success
		return nil
	}
}

func ResetRadio(app *gopi.AppInstance, args []string) error {
	app.Logger.Info("ResetRadio")
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	if mihome := MiHome(app); mihome == nil {
		return gopi.ErrAppError
	} else {
		return mihome.ResetRadio()
	}
}

func MeasureTemp(app *gopi.AppInstance, args []string) error {
	app.Logger.Info("MeasureTemp")
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	if mihome := MiHome(app); mihome == nil {
		return gopi.ErrAppError
	} else if celcius, err := mihome.MeasureTemperature(); err != nil {
		return err
	} else {
		fmt.Printf("Temperature=%v\u00B0C\n", celcius)
	}

	// Return success
	return nil
}

func TransmitOn(app *gopi.AppInstance, args []string) error {
	app.Logger.Info("TransmitOn %v", args)
	if sockets, err := toSockets(args); err != nil {
		return err
	} else if mihome := MiHome(app); mihome == nil {
		return gopi.ErrAppError
	} else if err := mihome.On(sockets...); err != nil {
		return err
	}
	// Return success
	return nil
}

func TransmitOff(app *gopi.AppInstance, args []string) error {
	app.Logger.Info("TransmitOff %v", args)
	if sockets, err := toSockets(args); err != nil {
		return err
	} else if mihome := MiHome(app); mihome == nil {
		return gopi.ErrAppError
	} else if err := mihome.Off(sockets...); err != nil {
		return err
	}
	// Return success
	return nil
}
