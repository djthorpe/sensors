/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"errors"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func ServerLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	if server := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil {
		return errors.New("Module rpc/server missing")
	} else {
		// Create the helloworld module
		if service := NewService(); service == nil {
			return errors.New("Service missing")
		} else {
			// Start server - will end when Stop is called
			server.Start(service)
		}
	}

	// wait for done
	<-done

	// Bomb out
	return nil
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Get the server
	if server := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil {
		return errors.New("Module rpc/server missing")
	} else {
		// Wait for CTRL+C
		app.Logger.Info("Press CTRL+C to finish")
		app.WaitForSignal()

		// Indicate we want to stop the server - shutdown
		// after we have serviced requests
		server.Stop(false)
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/server")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, ServerLoop))
}
