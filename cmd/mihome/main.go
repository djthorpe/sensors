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
	"os"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/protocol/ook"
	_ "github.com/djthorpe/sensors/protocol/openthings"
	_ "github.com/djthorpe/sensors/sys/ener314rt"
	_ "github.com/djthorpe/sensors/sys/rfm69"
)

var (
	mihome *MiHomeApp
)

////////////////////////////////////////////////////////////////////////////////

func ProductName(message sensors.OTMessage) string {
	if message.Manufacturer() == sensors.OT_MANUFACTURER_ENERGENIE {
		product := strings.TrimPrefix(fmt.Sprint(sensors.MiHomeProduct(message.Product())), "MIHOME_PRODUCT_")
		return product
	} else {
		manufacturer := strings.TrimPrefix(fmt.Sprint(message.Manufacturer()), "OT_MANUFACTURER_")
		return fmt.Sprintf("%v<%02X>", manufacturer, message.Product())
	}
}

func Sensor(message sensors.OTMessage) string {
	return fmt.Sprintf("%s<%05X>", ProductName(message), message.Sensor())
}

func PrintEvent(evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok {
		records := ""
		for _, record := range message.Records() {
			records += fmt.Sprint(record) + " "
		}
		fmt.Printf("%20s %s\n", Sensor(message), strings.TrimSpace(records))
	} else if message, ok := evt.(sensors.OOKMessage); ok {
		fmt.Println("OOKMessage", message)
	} else {
		fmt.Println("Other", evt)
	}
}

/*
type OTProto interface {
	Proto

	// Create a new message
	New(manufacturer OTManufacturer, product uint8, sensor uint32) (OTMessage, error)
}

type OTMessage interface {
	Message

	Manufacturer() OTManufacturer
	Product() uint8
	Sensor() uint32
	Records() []OTRecord
}
*/

func Receive(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {

	// Subscribe to events from the ENER314RT
	mihome := app.ModuleInstance("sensors/ener314rt").(gopi.Publisher)
	if mihome == nil {
		return gopi.ErrAppError
	}
	evts := mihome.Subscribe()

	// Event loop
	start <- gopi.DONE
FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		case evt := <-evts:
			PrintEvent(evt)
		}
	}

	// Unsubscribe
	mihome.Unsubscribe(evts)

	// Return success
	return nil
}

func Run(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
	// On exit of this method, send an exit signal to Main
	defer app.SendSignal()

	// Get command-line arguments
	if mihome = NewApp(app); mihome == nil {
		return gopi.ErrHelp
	} else {
		// Indicate tasks is running
		start <- gopi.DONE
		if err := mihome.Run(stop); err != nil {
			return err
		}
	}
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	// Wait for signal
	app.WaitForSignal()
	// Send cancel for long-running operations
	if mihome != nil {
		mihome.Cancel()
	}
	// Send done signal to tasks
	done <- gopi.DONE
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("gpio", "sensors/ener314rt", "sensors/protocol/ook", "sensors/protocol/openthings")

	// Usage
	config.AppFlags.SetUsageFunc(func(flags *gopi.Flags) {
		fmt.Fprintf(os.Stderr, "Usage:\n  %v <flags...> <%v>\n\n", flags.Name(), strings.Join(CommandNames(), "|"))
		PrintCommands(os.Stderr)
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flags.PrintDefaults()
	})

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, Receive, Run))
}
