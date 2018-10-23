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
	"fmt"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	mihome "github.com/djthorpe/sensors/rpc/grpc/mihome"
)

var (
	// Client communication object
	client *mihome.Client
)

////////////////////////////////////////////////////////////////////////////////

func GetClient(app *gopi.AppInstance) (*mihome.Client, error) {
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
	} else if client_ := pool.NewClient("sensors.MiHome", conn); client_ == nil {
		return nil, gopi.ErrAppError
	} else if client = client_.(*mihome.Client); client == nil {
		return nil, gopi.ErrAppError
	} else if err := client.Ping(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func ReceiveTask(app *gopi.AppInstance, start chan<- struct{}, done <-chan struct{}) error {

	// Connect to the service
	client, err := GetClient(app)
	if err != nil {
		return err
	}

	// Message to start the Main method
	start <- gopi.DONE

	messages := make(chan *mihome.Message)
	go func() {
		for {
			message := <-messages
			if message == nil {
				fmt.Println("END OF MESSAGES")
				break
			} else {
				fmt.Println(message)
			}
		}
	}()

	// Receive messages until done
	if err := client.Receive(done, messages); err != nil {
		return err
	} else {
		return nil
	}
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Wait until CTRL+C is pressed
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()
	done <- gopi.DONE

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client/mihome")

	// Set the RPCServiceRecord for server discovery
	config.Service = "mihome"

	// Address for remote service
	config.AppFlags.FlagString("addr", "", "Gateway address")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, ReceiveTask))
}
