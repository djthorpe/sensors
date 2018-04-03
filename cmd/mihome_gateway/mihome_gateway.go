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
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func EventProcess(evt gopi.RPCEvent, server gopi.RPCServer, discovery gopi.RPCServiceDiscovery) error {
	switch evt.Type() {
	case gopi.RPC_EVENT_SERVER_STARTED:
		fmt.Printf("Server started, addr=%v\n", server.Addr())
		if err := discovery.Register(server.Service("x", "mihome")); err != nil {
			return err
		}
	case gopi.RPC_EVENT_SERVER_STOPPED:
		fmt.Printf("Server stopped\n")
		// TODO: Unregister (same as register but with ttl=0)
	default:
		fmt.Printf("Error: Unhandled event: %v\n", evt)
	}
	return nil
}

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	if server := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil {
		return errors.New("Module rpc/server missing")
	} else if mdns := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); mdns == nil {
		return errors.New("Module rpc/discovery missing")
	} else {
		// Listen for events
		events := server.Subscribe()
	FOR_LOOP:
		for {
			select {
			case evt := <-events:
				if rpc_evt, ok := evt.(gopi.RPCEvent); rpc_evt != nil && ok {
					EventProcess(rpc_evt, server, mdns)
				}
			case <-done:
				break FOR_LOOP
			}
		}

		// Stop listening for events
		server.Unsubscribe(events)
	}

	return nil
}

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
	config := gopi.NewAppConfig("rpc/server", "rpc/discovery")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, ServerLoop, EventLoop))
}
