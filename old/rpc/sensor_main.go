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
)

////////////////////////////////////////////////////////////////////////////////

// EventProcess registers and unregisters discovery
func EventProcess(evt gopi.RPCEvent, server gopi.RPCServer, discovery gopi.RPCServiceDiscovery) error {
	switch evt.Type() {
	case gopi.RPC_EVENT_SERVER_STARTED:
		fmt.Printf("Server started, addr=%v\n", server.Addr())
		if err := discovery.Register(server.Service()); err != nil {
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

// EventLoop processes stop and start messages
func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	discovery, ok := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)
	if discovery == nil || ok == false {
		return errors.New("rpc/discovery missing")
	}

	// Listen for events
	c := server.Subscribe()
FOR_LOOP:
	for {
		select {
		case evt := <-c:
			if rpc_evt, ok := evt.(gopi.RPCEvent); rpc_evt != nil && ok {
				EventProcess(rpc_evt, server, discovery)
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// Stop listening for events
	server.Unsubscribe(c)

	return nil
}

// ServerLoop starts the server
func ServerLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	// Create the sensors module
	if service := new(SensorService); service == nil {
		return errors.New("SensorService missing")
	} else if err := server.Start(service); err != nil {
		return err
	}

	// wait for done
	<-done

	// Bomb out
	return nil
}

// MainLoop simply waits for signal and then stops server
func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	app.WaitForSignal()

	// Indicate we want to stop the server - shutdown
	// after we have serviced requests
	server.Stop(false)

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig("rpc/server")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, ServerLoop, EventLoop))
}
