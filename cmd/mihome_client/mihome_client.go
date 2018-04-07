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
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"

	// Protocol Buffer definitions
	pb "github.com/djthorpe/sensors/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////

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

	if services, err := client.Connect(); err != nil {
		return err
	} else if has_service := HasService(services, "MiHome"); has_service == false {
		return fmt.Errorf("Invalid MiHome Gateway address (services are %v)", strings.Join(services, ","))
	} else if service_obj, err := client.NewService(reflect.ValueOf(pb.NewMiHomeClient)); err != nil {
		return err
	} else if service, ok := service_obj.(pb.MiHomeClient); service == nil || ok == false {
		_ = service_obj.(pb.MiHomeClient)
		return errors.New("Invalid MiHome service")
	} else {
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
					fmt.Printf("Event=%v\n", message)
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
