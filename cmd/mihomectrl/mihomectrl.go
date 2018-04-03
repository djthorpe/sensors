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
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/mutablehome"
	"github.com/djthorpe/sensors"
	"github.com/olekukonko/tablewriter"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/mutablehome/sys/linux"
	_ "github.com/djthorpe/sensors/hw/energenie"
	_ "github.com/djthorpe/sensors/hw/rfm69"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Command struct {
	description string
	callback    func(app *gopi.AppInstance) error
}

type State struct {
	commands chan *Command
	mihome   sensors.MiHome
	devices  mutablehome.Devices
	cancel   context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS AND VARIABLES

var (
	COMMANDS = map[string]*Command{
		"reset":   &Command{"Reset the radio module", CommandReset},
		"rx":      &Command{"Receive Data Mode", CommandReceive},
		"temp":    &Command{"Measure Temperature", CommandTemp},
		"devices": &Command{"List Devices", CommandDevices},
	}
)

var (
	state *State
)

////////////////////////////////////////////////////////////////////////////////
// RECEIVE MESSAGES

////////////////////////////////////////////////////////////////////////////////
// COMMANDS

func CommandReset(app *gopi.AppInstance) error {
	app.Logger.Info("Resetting device")
	return state.mihome.ResetRadio()
}

func CommandReceive(app *gopi.AppInstance) error {
	app.Logger.Info("Receiving data")
	timeout, _ := app.AppFlags.GetDuration("timeout")

	// Obtain the context
	var ctx context.Context
	if timeout != 0 {
		ctx, state.cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, state.cancel = context.WithCancel(context.Background())
	}

	// Set cancel function for duration of the receive
	defer func() { state.cancel = nil }()

	// Print out the header
	fmt.Printf("%7s %2v %-10s %-22s\n", "Time", "Sz", "Sensor ID", "Pair Status")
	fmt.Printf("%7s %2v %-10s %-22s\n", "-------", "--", "----------", "----------------------")

	// Perform the receive
	return state.mihome.Receive(ctx, sensors.MIHOME_MODE_MONITOR)
}

func CommandTemp(app *gopi.AppInstance) error {
	app.Logger.Info("Measuring Temperature")
	if temp, err := state.mihome.MeasureTemperature(); err != nil {
		return err
	} else {
		fmt.Printf("Temperature=%vC\n", temp)
		return nil
	}
}

func CommandDevices(app *gopi.AppInstance) error {
	if device_db := app.ModuleInstance("mutablehome/devices").(mutablehome.Devices); device_db == nil {
		return fmt.Errorf("Missing devices database")
	} else {
		// Print out a list of all devices
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Device", "Name", "Product", "Type", "Status"})
		devices := device_db.Devices(mutablehome.DEVICE_TYPE_ANY, mutablehome.PAIR_STATUS_ANY)
		for _, device := range devices {
			table.Append([]string{
				device.Hash(),
				device.Name,
				fmt.Sprint(device.ProductId),
				fmt.Sprint(device.Type),
				fmt.Sprint(device.PairStatus),
			})
		}
		table.Render()
	}
	return nil
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
// PROCESS EVENTS

func ProcessEvent(evt sensors.OTEvent) error {
	// Deal with corrupted packets
	if evt.Reason() != nil {
		fmt.Printf("%-20s %v\n", evt.Timestamp().Format(time.Stamp), evt.Reason())
		return nil
	}
	// Get the message
	message := evt.Message()
	// Get device
	if device, err := state.devices.Device(uint64(message.SensorID()), mutablehome.DEVICE_TYPE_ENERGENIE_MONITOR, uint64(message.ProductID())); err != nil {
		return err
	} else {
		fmt.Printf("%7s %2v %-10s %-22s\n",
			evt.Timestamp().Format(time.Kitchen),
			message.Size(),
			device.Hash(),
			device.PairStatus)
	}
	// Print out records
	for _, record := range message.Records() {
		fmt.Printf("%30s %s\n", "", record)
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RECEIEVE MESSAGES LOOP

func ReceiveLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome)
	if mihome == nil {
		return gopi.ErrAppError
	}
	events := mihome.Subscribe()

FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case e := <-events:
			if err := ProcessEvent(e.(sensors.OTEvent)); err != nil {
				return err
			}
		}
	}

	// Unsubscribe
	mihome.Unsubscribe(events)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// EXECUTE COMMAND LOOP

func CommandLoop(app *gopi.AppInstance, done <-chan struct{}) error {

FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case cmd := <-state.commands:
			// Execute a command
			if err := cmd.callback(app); err != nil {
				return err
			}
		}
	}
	app.Logger.Info("Press CTRL+C to quit")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Get devices and mihome module instances
	if state.mihome = app.ModuleInstance("sensors/mihome").(sensors.MiHome); state.mihome == nil {
		return gopi.ErrAppError
	}
	if state.devices = app.ModuleInstance("mutablehome/devices").(mutablehome.Devices); state.devices == nil {
		return gopi.ErrAppError
	}

	// Get arguments on command line
	args := app.AppFlags.Args()
	if len(args) == 0 {
		done <- gopi.DONE
		return gopi.ErrHelp
	} else if err := state.AddCommands(app.AppFlags.Args()); err != nil {
		done <- gopi.DONE
		return err
	}

	// Wait for CTRL+C
	app.WaitForSignal()

	// Send out cancel
	state.Cancel()

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func NewState() *State {
	this := new(State)
	return this
}

func (this *State) AddCommands(args []string) error {
	// Create the channel
	this.commands = make(chan *Command, len(args))

	// Emit the commands to execute
	for _, arg := range args {
		if command, exists := COMMANDS[arg]; exists == false {
			return fmt.Errorf("Invalid command: %v", arg)
		} else {
			this.commands <- command
		}
	}

	// Return success
	return nil
}

func (this *State) Cancel() {
	if this.cancel != nil {
		this.cancel()
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/mihome", "mutablehome/devices")
	config.AppFlags.SetUsageFunc(Usage)

	// Timeout flag for receive timeout
	config.AppFlags.FlagDuration("timeout", 0, "Timeout for receive mode")

	// Create the application state
	state = NewState()

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, CommandLoop, ReceiveLoop))
}
