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

////////////////////////////////////////////////////////////////////////////////

var (
	start = make(chan *MiHomeApp)
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if mihome := NewApp(app); mihome == nil {
		start <- nil
		done <- gopi.DONE
		return gopi.ErrHelp
	} else {
		start <- mihome
		app.WaitForSignal()
		mihome.Cancel()
		done <- gopi.DONE
	}

	return nil
}

func Tasks(app *gopi.AppInstance, done <-chan struct{}) error {
	defer app.SendSignal()

	if mihome := <-start; mihome == nil {
		<-done
		return nil
	} else if err := mihome.Run(); err != nil {
		<-done
		return err
	}

	return nil
}

func Messages(app *gopi.AppInstance, done <-chan struct{}) error {
	mihome := app.ModuleInstance("sensors/ener314rt").(gopi.Publisher)
	evts := mihome.Subscribe()
FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case evt := <-evts:
			fmt.Println("got event=%v", evt)
		}
	}

	// Unsubscribe
	mihome.Unsubscribe(evts)

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
	os.Exit(gopi.CommandLineTool(config, Main, Tasks, Messages))
}
