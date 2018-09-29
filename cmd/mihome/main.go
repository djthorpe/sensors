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
	_ "github.com/djthorpe/sensors/sys/ener314rt"
	_ "github.com/djthorpe/sensors/sys/rfm69"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if args := app.AppFlags.Args(); len(args) == 0 {
		return gopi.ErrHelp
	} else {
		if err := RunCommand(app, args[0], args[1:]); err != nil {
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
	config := gopi.NewAppConfig("gpio", "sensors/ener314rt")

	// Usage
	config.AppFlags.SetUsageFunc(func(flags *gopi.Flags) {
		fmt.Fprintf(os.Stderr, "Usage:\n  %v <flags...> <%v>\n\n", flags.Name(), strings.Join(CommandNames(), "|"))
		PrintCommands(os.Stderr)
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flags.PrintDefaults()
	})

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
