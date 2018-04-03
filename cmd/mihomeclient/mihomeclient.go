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
	"errors"
	"fmt"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func BrowseLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	if mdns := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); mdns == nil {
		return errors.New("Module rpc/discovery missing")
	} else {
		events := mdns.Subscribe()
	FOR_LOOP:
		for {
			select {
			case <-done:
				break FOR_LOOP
			case evt := <-events:
				fmt.Println(evt)
			}
		}
	}
	// Success
	return nil
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Create context
	timeout, _ := app.AppFlags.GetDuration("timeout")
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	if mdns := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); mdns == nil {
		return errors.New("Module rpc/discovery missing")
	} else if err := mdns.Browse(ctx, "_mihome._tcp"); err != nil {
		return err
	}

	if client := NewClient(); client == nil {
		return errors.New("Client could not be created")
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/discovery")
	config.AppFlags.FlagDuration("timeout", 2*time.Second, "RPC discovery timeout")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, BrowseLoop))
}
