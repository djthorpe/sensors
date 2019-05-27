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
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi-rpc/sys/dns-sd"
	_ "github.com/djthorpe/gopi-rpc/sys/grpc"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/sys/sensordb"

	// Clients
	_ "github.com/djthorpe/sensors/rpc/grpc/mihome"
)

const (
	DISCOVERY_TIMEOUT = 700 * time.Millisecond
)

////////////////////////////////////////////////////////////////////////////////

func Conn(app *gopi.AppInstance) (gopi.RPCServiceRecord, error) {
	addr, _ := app.AppFlags.GetString("addr")
	timeout, exists := app.AppFlags.GetDuration("rpc.timeout")
	if exists == false {
		timeout = DISCOVERY_TIMEOUT
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if service, _, _, err := app.Service(); err != nil {
		return nil, err
	} else {
		service_ := fmt.Sprintf("_%v._tcp", service)
		pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
		if services, err := pool.Lookup(ctx, service_, addr, 0); err != nil {
			return nil, err
		} else if len(services) == 0 {
			return nil, gopi.ErrNotFound
		} else if len(services) > 1 {
			var names []string
			for _, service := range services {
				names = append(names, strconv.Quote(service.Name()))
			}
			return nil, fmt.Errorf("More than one service returned, use -addr to choose between %v", strings.Join(names, ","))
		} else {
			return services[0], nil
		}
	}
}

func MiHomeStub(app *gopi.AppInstance, sr gopi.RPCServiceRecord) (sensors.MiHomeClient, error) {
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	if sr == nil || pool == nil {
		return nil, gopi.ErrBadParameter
	} else if conn, err := pool.Connect(sr, 0); err != nil {
		return nil, err
	} else if stub := pool.NewClient("mihome.MiHome", conn); stub == nil {
		return nil, gopi.ErrBadParameter
	} else if stub_, ok := stub.(sensors.MiHomeClient); ok == false {
		_ = stub.(sensors.MiHomeClient)
		return nil, fmt.Errorf("Stub is not an sensors.MiHomeClient")
	} else if err := stub_.Ping(); err != nil {
		return nil, err
	} else {
		return stub_, nil
	}
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if record, err := Conn(app); err != nil {
		return err
	} else if client, err := MiHomeStub(app, record); err != nil {
		return err
	} else if err := Run(app, client); err != nil {
		return err
	}

	// Success
	return nil
}

func Usage(flags *gopi.Flags) {
	fh := os.Stdout

	fmt.Fprintf(fh, "%v: MiHome Message Rx and Tx\nhttps://github.com/djthorpe/sensors/\n\n", flags.Name())
	fmt.Fprintf(fh, "Syntax:\n\n")
	fmt.Fprintf(fh, "  %v (<flags>...)\n\n", flags.Name())
	fmt.Fprintf(fh, "Command line flags:\n\n")
	flags.PrintDefaults()
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/mihome:client", "sensordb", "discovery")

	// Set subtype
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_SUBTYPE, "mihome")

	// Set usage function
	config.AppFlags.SetUsageFunc(Usage)

	// Set flags
	config.AppFlags.FlagString("addr", "", "Service name or gateway address")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main))
}
