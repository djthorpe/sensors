/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved

   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"os"
	"sync"

	"github.com/djthorpe/sensors"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"

	// RPC Services
	_ "github.com/djthorpe/sensors/cmd/mihome_gateway/service"
)

////////////////////////////////////////////////////////////////////////////////

type Command uint

const (
	COMMAND_RESET Command = iota
)

////////////////////////////////////////////////////////////////////////////////

var (
	cancel   context.CancelFunc
	lock     sync.Mutex
	commands chan Command
)

////////////////////////////////////////////////////////////////////////////////

func GetContext() context.Context {
	lock.Lock()
	defer lock.Unlock()

	// Cancel previously running receive
	if cancel != nil {
		cancel()
	}

	// Create a new context
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())

	// Return context
	return ctx
}

func Cancel() {
	lock.Lock()
	defer lock.Unlock()
	if cancel != nil {
		cancel()
		cancel = nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func ReceiveLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome)
	events := mihome.Subscribe()
FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case evt := <-events:
			if ot_evt, ok := evt.(sensors.OTEvent); ot_evt != nil && ok {
				app.Logger.Info("Received event: %v", evt)
			}
		}
	}

	// Stop listening for messages
	mihome.Unsubscribe(events)

	// Finished
	return nil
}

func CommandLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome)
FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		default:
			app.Logger.Info("In Receive mode")
			if err := mihome.Receive(GetContext(), sensors.MIHOME_MODE_MONITOR); err != nil {
				app.Logger.Error("CommandLoop: %v", err)
			}
			app.Logger.Info("End of receive mode")
		}
	}

	// Finished
	return nil
}

func CommandCancel(app *gopi.AppInstance, done <-chan struct{}) error {
	// Wait for done and then cancel
	app.Logger.Debug("CommandCancel waiting for done")
	<-done
	app.Logger.Debug("CommandCancel cancelling")
	Cancel()
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("service/mihome:grpc")

	// Set the RPCServiceRecord for server discovery
	config.Service = "mihome"

	// Channel for incoming commands
	commands = make(chan Command)

	// Run the server and register all the services
	// Note the CommandLoop needs to go last as it blocks on Receive() until
	// Cancel is called from the CommandCancel task
	os.Exit(gopi.RPCServerTool(config))
}
