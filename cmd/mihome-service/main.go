package main

import (
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi-rpc"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/gopi-rpc/sys/dns-sd"
	_ "github.com/djthorpe/gopi-rpc/sys/grpc"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/protocol/ook"
	_ "github.com/djthorpe/sensors/protocol/openthings"
	_ "github.com/djthorpe/sensors/sys/ener314rt"
	_ "github.com/djthorpe/sensors/sys/mihome"
	_ "github.com/djthorpe/sensors/sys/rfm69"

	// RPC Services
	_ "github.com/djthorpe/sensors/rpc/grpc/mihome"
)

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/mihome:service", "sensors/protocol/ook", "sensors/protocol/openthings", "discovery")

	// Set subtype
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_SUBTYPE, "mihome")

	// Run the server and register all the services
	os.Exit(rpc.Server(config))
}
