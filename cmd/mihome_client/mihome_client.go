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
	"io"
	"os"
	"reflect"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"

	// Protocol Buffer definitions
	pb "github.com/djthorpe/sensors/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	start chan pb.MiHomeClient
)

////////////////////////////////////////////////////////////////////////////////
// RECEIEVE MESSAGES LOOP

func ReceiveLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	// Receive the service
	service := <-start

	// Create the context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Call cancel in the background
	go func() {
		<-done
		cancel()
	}()

	// Receive a stream of messages
	stream, err := service.Receive(ctx, &pb.ReceiveRequest{})
	if err != nil {
		return err
	}
	for {
		if message, err := stream.Recv(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			fmt.Printf("Event=%v\n", message)
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN

func HasService(services []string, service string) bool {
	if services == nil {
		return false
	}
	for _, v := range services {
		if v == service {
			return true
		}
	}
	return false
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	client := app.ModuleInstance("rpc/client").(gopi.RPCClientConn)
	start = make(chan pb.MiHomeClient)

	if services, err := client.Connect(); err != nil {
		return err
	} else if has_service := HasService(services, "MiHome"); has_service == false {
		return fmt.Errorf("Invalid MiHome gateway address (missing service)")
	} else if obj, err := client.NewService(reflect.ValueOf(pb.NewMiHomeClient)); err != nil {
		return err
	} else if service, ok := obj.(pb.MiHomeClient); service == nil || ok == false {
		return errors.New("Invalid MiHome service")
	} else {
		// Send the service to the receive loop
		start <- service
	}

	// Wait for signal
	app.WaitForSignal()

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, ReceiveLoop))
}
