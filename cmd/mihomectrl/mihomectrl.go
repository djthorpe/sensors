/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Control Energenie MiHome devices
package main

import (
	"context"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/hw/energenie"
	_ "github.com/djthorpe/sensors/hw/rfm69"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Command struct {
	description string
	callback    func(app *gopi.AppInstance, mihome sensors.MiHome) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS AND VARIABLES

var (
	COMMANDS = map[string]*Command{
		"reset": &Command{"Reset the radio module", CommandReset},
		"rx":    &Command{"Receive Data Mode", CommandReceive},
	}
)

////////////////////////////////////////////////////////////////////////////////
// RESET COMMAND

func CommandReset(app *gopi.AppInstance, mihome sensors.MiHome) error {
	return mihome.ResetRadio()
}

func CommandReceive(app *gopi.AppInstance, mihome sensors.MiHome) error {
	timeout, _ := app.AppFlags.GetDuration("timeout")
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return mihome.Receive(ctx, sensors.MIHOME_MODE_MONITOR)
}

////////////////////////////////////////////////////////////////////////////////
// HELP FUNCTION

func Usage(flags *gopi.Flags) {
	fmt.Fprintf(os.Stderr, "Usage of %v:\n\n", flags.Name())
	fmt.Fprintf(os.Stderr, "     %v <flags>... <commands>...\n\n", flags.Name())
	fmt.Fprintf(os.Stderr, "Commands:\n\n")

	for key, command := range COMMANDS {
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", key, command.description)
	}

	fmt.Fprintf(os.Stderr, "\nFlags:\n\n")
	flags.PrintDefaults()
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome); mihome == nil {
		return gopi.ErrAppError
	} else if args := app.AppFlags.Args(); len(args) == 0 {
		return gopi.ErrHelp
	} else {
		// Collate the commands to execute
		commands := make([]*Command, 0, len(args))
		for _, arg := range app.AppFlags.Args() {
			if command, exists := COMMANDS[arg]; exists == false {
				return fmt.Errorf("Invalid command: %v", arg)
			} else {
				commands = append(commands, command)
			}
		}
		// Execute the commands
		for i, command := range commands {
			app.Logger.Info("Running command: %v (%v)", args[i], command.description)
			if err := command.callback(app, mihome); err != nil {
				return err
			}
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/mihome")
	config.AppFlags.SetUsageFunc(Usage)

	// Timeout flag for receive timeout
	config.AppFlags.FlagDuration("timeout", 0, "Timeout for receive mode")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
