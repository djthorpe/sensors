package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/sys/sensordb"

	// RPC Services
	_ "github.com/djthorpe/sensors/rpc/grpc/sensordb"
)

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/service/sensordb")

	// Set the RPCServiceRecord for server discovery
	config.Service = "sensordb"

	// Run the server and register all the services
	os.Exit(gopi.RPCServerTool(config))
}
