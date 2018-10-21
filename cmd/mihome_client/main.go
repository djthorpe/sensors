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

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	// Obtain Client Pool and Gateway address
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	addr, _ := app.AppFlags.GetString("addr")

	// Lookup remote service and run
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	if records, err := pool.Lookup(ctx, "", addr, 1); err != nil {
		done <- gopi.DONE
		return err
	} else if len(records) == 0 {
		done <- gopi.DONE
		return gopi.ErrDeadlineExceeded
	} else if conn, err := pool.Connect(records[0], 0); err != nil {
		done <- gopi.DONE
		return err
	} else if client_ := pool.NewClient("sensors.MiHome", conn); client_ == nil {
		done <- gopi.DONE
		return gopi.ErrAppError
	} else if client := client_.(*mihome.Client); client == nil {
		done <- gopi.DONE
		return gopi.ErrAppError
	} else {
		if celcius, err := client.MeasureTemperature(); err != nil {
			done <- gopi.DONE
			return err
		} else {
			// Report the temperature
			fmt.Println("celcius=", celcius)
		}

		if protos, err := client.Protocols(); err != nil {
			done <- gopi.DONE
			return err
		} else {
			// Report the protocols available
			fmt.Println("protos=", protos)
		}

		// Wait until CTRL+C is pressed
		app.Logger.Info("Waiting for CTRL+C")
		app.WaitForSignal()
		done <- gopi.DONE
	}

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
	os.Exit(gopi.CommandLineTool2(config, Main))
}
