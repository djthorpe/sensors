/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (

	// Frameworks
	"fmt"

	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi-rpc/sys/dns-sd"
	_ "github.com/djthorpe/gopi-rpc/sys/grpc"
	_ "github.com/djthorpe/gopi/sys/logger"

	// Clients
	_ "github.com/djthorpe/sensors/rpc/grpc/mihome"
)

////////////////////////////////////////////////////////////////////////////////

func ReceiveTask(app *gopi.AppInstance, start chan<- struct{}, done <-chan struct{}) error {
	/*
		// Connect to the service
		client, err := GetClient(app)
		if err != nil {
			app.Logger.Error("ReceiveTask: %v", err)
			return err
		}

		// Message to start the Main method
		start <- gopi.DONE

		// Create a goroutine to receive the messages and print them out and end the goroutine
		// when the message channel is closed. Null events are sent regularly to ensure the
		// channel is still active, ignore these.
		messages := make(chan sensors.Message)
		go func() {
			for {
				message := <-messages
				if message == nil {
					// Closed channel
					break
				} else {
					app.Logger.Info("%v", message)
				}
			}
		}()

		// Receive messages until done
		if err := client.Receive(done, messages); err != nil {
			return err
		} else {
			return nil
		}
	*/
	return nil
}

func Run(app *gopi.AppInstance, client sensors.MiHomeClient) error {
	fmt.Println(client)
	return gopi.ErrNotImplemented
}
