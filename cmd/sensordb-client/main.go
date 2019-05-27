/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	sensordb "github.com/djthorpe/sensors/rpc/grpc/sensordb"
)

var (
	// Client communication object
	client *sensordb.Client
)

////////////////////////////////////////////////////////////////////////////////

func GetClient(app *gopi.AppInstance) (*sensordb.Client, error) {
	// Obtain Client Pool and Gateway address
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	addr, _ := app.AppFlags.GetString("addr")

	// Check for existing client
	if client != nil {
		return client, nil
	}

	// Lookup remote service and run
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	if records, err := pool.Lookup(ctx, "", addr, 1); err != nil {
		return nil, err
	} else if len(records) == 0 {
		return nil, gopi.ErrDeadlineExceeded
	} else if conn, err := pool.Connect(records[0], 0); err != nil {
		return nil, err
	} else if client_ := pool.NewClient("sensors.SensorDB", conn); client_ == nil {
		return nil, gopi.ErrAppError
	} else if client = client_.(*sensordb.Client); client == nil {
		return nil, gopi.ErrAppError
	} else if err := client.Ping(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func EventLoop(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {

	// Connect to the service
	client, err := GetClient(app)
	if err != nil {
		app.Logger.Error("EventLoop: Client: %v", err)
		app.SendSignal()
		return err
	}

	// Message to start the Main method
	start <- gopi.DONE

	// ping every one second
	ping := time.NewTicker(time.Second)

	// Wait for stop
FOR_LOOP:
	for {
		select {
		case <-stop:
			// Break out of the event loop
			ping.Stop()
			break FOR_LOOP
		case <-ping.C:
			// Ensure connection is still up periodically
			if err := client.Ping(); err != nil {
				app.Logger.Error("EventLoop: Ping: %v", err)
			}
		}
	}

	// End Loop
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Main method simply waits until CTRL+C is pressed
	// and then signals background tasks to end
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()
	done <- gopi.DONE

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client/sensordb")

	// Set the RPCServiceRecord for server discovery
	config.Service = "sensordb"

	// Address for remote service
	config.AppFlags.FlagString("addr", "", "Gateway address")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, EventLoop))
}
