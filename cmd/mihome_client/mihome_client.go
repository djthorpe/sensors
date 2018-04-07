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

func HasService(app *gopi.AppInstance, client gopi.RPCClient, service string) (bool, error) {
	if services, err := client.Modules(); err != nil {
		return false, err
	} else {
		for _, v := range services {
			if v == service {
				return true, nil
			}
		}
		return false, nil
	}
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	client := app.ModuleInstance("rpc/client").(gopi.RPCClient)

	if err := client.Connect(); err != nil {
		return err
	} else if has_service, err := HasService(app, client, "MiHome"); err != nil {
		return err
	} else if has_service == false {
		return errors.New("Invalid MiHome Gateway address")
	} else {
		// Create the gRPC client - pass in the constructor method
		service := client.NewService(reflect.ValueOf(pb.NewMiHomeClient)).(pb.MiHomeClient)
		// Receive a stream of messages
		if stream, err := service.Receive(context.Background(), &pb.ReceiveRequest{}); err != nil {
			return err
		} else {
			for {
				if message, err := stream.Recv(); err == io.EOF {
					break
				} else if err != nil {
					return err
				} else {
					fmt.Println(message)
				}
			}
		}
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
