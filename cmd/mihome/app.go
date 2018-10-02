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
	Fn          func([]string) error
	Description string
}

type MiHomeApp struct {
	log    gopi.Logger
	gpio   gopi.GPIO
	mihome sensors.MiHome
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	COMMANDS = map[string]CommandFunc{
		"reset_gpio":   CommandFunc{nil, "Reset GPIO"},
		"reset_radio":  CommandFunc{nil, "Reset RFM69 Radio"},
		"measure_temp": CommandFunc{nil, "Measure Temperature"},
		"on":           CommandFunc{nil, "On TX (optionally use 1,2,3,4 as additional argument)"},
		"off":          CommandFunc{nil, "Off TX (optionally use 1,2,3,4 as additional argument)"},
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

////////////////////////////////////////////////////////////////////////////////
// APP Commands

func NewApp(app *gopi.AppInstance) *MiHomeApp {
	this := new(MiHomeApp)
	this.log = app.Logger
	this.gpio = app.ModuleInstance("gpio").(gopi.GPIO)
	this.mihome = app.ModuleInstance("sensors/ener314rt").(sensors.MiHome)
	return this
}

func (this *MiHomeApp) Run(cmd string, args []string) error {
	if command, exists := COMMANDS[cmd]; exists == false {
		return gopi.ErrHelp
	} else if err := command.Fn(args); err != nil {
		return err
	}
	return nil
}

func (this *MiHomeApp) ResetGPIO(args []string) error {
	// Ensure pins are in correct state for SPI0
	this.gpio.SetPinMode(gopi.GPIOPin(7), gopi.GPIO_OUTPUT)
	this.gpio.SetPinMode(gopi.GPIOPin(8), gopi.GPIO_OUTPUT)
	this.gpio.SetPinMode(gopi.GPIOPin(9), gopi.GPIO_ALT0)
	this.gpio.SetPinMode(gopi.GPIOPin(10), gopi.GPIO_ALT0)
	this.gpio.SetPinMode(gopi.GPIOPin(11), gopi.GPIO_ALT0)
	// Success
	return nil
}

func (this *MiHomeApp) ResetRadio(args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	return this.mihome.ResetRadio()
}

func (this *MiHomeApp) MeasureTemp(args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	if celcius, err := this.mihome.MeasureTemperature(); err != nil {
		return err
	} else {
		fmt.Printf("Temperature=%v\u00B0C\n", celcius)
	}

	// Return success
	return nil
}

func (this *MiHomeApp) TransmitOn(args []string) error {
	if sockets, err := toSockets(args); err != nil {
		return err
	} else if err := this.mihome.On(sockets...); err != nil {
		return err
	}
	// Return success
	return nil
}

func (this *MiHomeApp) TransmitOff(args []string) error {
	if sockets, err := toSockets(args); err != nil {
		return err
	} else if err := this.mihome.Off(sockets...); err != nil {
		return err
	}
	// Return success
	return nil
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
